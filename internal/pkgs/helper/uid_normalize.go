package helper

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyUID         = errors.New("uid kartu tidak boleh kosong")
	ErrInvalidUIDLength = errors.New("panjang uid kartu tidak valid (harus 4-10 byte / 8-20 karakter hex)")
	ErrInvalidUIDFormat = errors.New("uid kartu harus berupa karakter heksadesimal (0-9, A-F)")

	uidReplacer = strings.NewReplacer(
		":", "",
		"-", "",
		" ", "",
		"\t", "",
		"\n", "",
		"\r", "",
	)
)

func NormalizeUID(raw string) string {
	cleaned := uidReplacer.Replace(raw)
	return strings.ToUpper(cleaned)
}

func ParseUID(raw string) (string, error) {
	normalized := NormalizeUID(raw)
	if normalized == "" {
		return "", ErrEmptyUID
	}

	n := len(normalized)
	if n < 8 || n > 20 || n%2 != 0 {
		return "", fmt.Errorf("%w: got %d hex chars", ErrInvalidUIDLength, n)
	}

	for i := range n {
		c := normalized[i]
		isDigit := c >= '0' && c <= '9'
		isUpperHex := c >= 'A' && c <= 'F'
		if !isDigit && !isUpperHex {
			return "", ErrInvalidUIDFormat
		}
	}

	return normalized, nil
}

func IsValidUID(raw string) bool {
	_, err := ParseUID(raw)
	return err == nil
}
