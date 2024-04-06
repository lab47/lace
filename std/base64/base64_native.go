package base64

import (
	"encoding/base64"

	. "github.com/lab47/lace/core"
)

func decodeString(s string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", StubNewError("Invalid base64 string: " + err.Error())
	}
	return string(decoded), nil
}

func encodeString(s string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(s)), nil
}
