package url

import (
	"net/url"

	. "github.com/candid82/joker/core"
)

func pathUnescape(s string) (string, error) {
	res, err := url.PathUnescape(s)
	if err != nil {
		return "", StubNewError("Error unescaping string: " + err.Error())
	}
	return res, nil
}

func queryUnescape(s string) (string, error) {
	res, err := url.QueryUnescape(s)
	if err != nil {
		return "", StubNewError("Error unescaping string: " + err.Error())
	}
	return res, nil
}
