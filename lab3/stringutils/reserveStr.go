package stringutils

/**Функция, возвращающая перевернутутю строку*/
func ReserveStr(s string) (result string) {
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		result += string(runes[i])
	}
	return result
}
