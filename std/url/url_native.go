package url

import (
	"net/url"

	. "github.com/candid82/joker/core"
)

func pathUnescape(s string) string {
	res, err := url.PathUnescape(s)
	if err != nil {
		panic(StubNewError("Error unescaping string: " + err.Error()))
	}
	return res
}

func queryUnescape(s string) string {
	res, err := url.QueryUnescape(s)
	if err != nil {
		panic(StubNewError("Error unescaping string: " + err.Error()))
	}
	return res
}
