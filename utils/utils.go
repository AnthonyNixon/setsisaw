package utils

import (
	"errors"
	"strings"
)

func GetTokenFromHeader(authHeader string) (string, error) {
	fields := strings.Fields(authHeader)
	if len(fields) != 2 {
		return "", errors.New("bad authentication header format")
	}

	return fields[1], nil
}
