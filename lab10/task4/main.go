package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

/**Завершение программы при ошибке*/
func check(err error) {
	if err != nil {
		fmt.Println("Ошибка:", err)
		os.Exit(1)
	}
}

/**Создание сертификата и нового ключа; подписывает сертификат родительским*/
func createCert(name string, template, parent *x509.Certificate, parentKey *rsa.PrivateKey) (*x509.Certificate, *rsa.PrivateKey) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	check(err)
	if parent == nil {
		parent = template
		parentKey = key
	}
	der, err := x509.CreateCertificate(rand.Reader, template, parent, &key.PublicKey, parentKey)
	check(err)
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	check(os.WriteFile("certs/"+name+".crt", certPEM, 0644))
	check(os.WriteFile("certs/"+name+".key", keyPEM, 0600))
	cert, err := x509.ParseCertificate(der)
	check(err)
	return cert, key
}

func main() {
	check(os.MkdirAll("certs", 0755))
	start := time.Now().Add(-time.Hour)
	end := start.AddDate(1, 0, 0)
	//центр сертификации
	ca, caKey := createCert("ca", &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Test CA"},
		NotBefore:             start,
		NotAfter:              end,
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}, nil, nil)
	//сертификат сервера
	createCert("server", &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: "server"},
		DNSNames:     []string{"localhost"},
		NotBefore:    start,
		NotAfter:     end,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}, ca, caKey)
	//сертификат клиента
	createCert("client", &x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject:      pkix.Name{CommonName: "client"},
		NotBefore:    start,
		NotAfter:     end,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}, ca, caKey)

	fmt.Println("Сертификаты и ключи созданы в папке certs:")
	fmt.Println("ca.crt, ca.key - центр сертификации")
	fmt.Println("server.crt, server.key - сервер")
	fmt.Println("client.crt, client.key - клиент")
}
