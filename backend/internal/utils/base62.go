package utils

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func EncodeBase62(num int64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	encoded := ""
	base := int64(len(base62Chars))

	for num > 0 {
		remainder := num % base
		encoded = string(base62Chars[remainder]) + encoded
		num = num / base
	}

	return encoded
}

func DecodeBase62(encoded string) int64 {
	var num int64
	base := int64(len(base62Chars))

	for _, char := range encoded {
		var value int64
		for i, c := range base62Chars {
			if c == char {
				value = int64(i)
				break
			}
		}
		num = num*base + value
	}

	return num
}
