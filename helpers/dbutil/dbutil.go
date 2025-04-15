package dbutil

import (
	"strings"
)

// ReplacePlaceholders replaces '?' with PostgreSQL-style numbered placeholders ($1, $2, ...).
func ReplacePlaceholders(query string) string {
	var builder strings.Builder
	count := 1

	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			builder.WriteString("$" + itoa(count))
			count++
		} else {
			builder.WriteByte(query[i])
		}
	}

	return builder.String()
}

// Fast integer to string conversion for small integers
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := [10]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	var buf [6]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = digits[n%10]
		n /= 10
	}
	return string(buf[i:])
}
