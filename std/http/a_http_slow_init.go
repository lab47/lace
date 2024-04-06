// This file is generated by generate-std.clj script. Do not edit manually!

package http

import (
	"fmt"
	"os"

	. "github.com/lab47/lace/core"
)

func InternsOrThunks() {
	if VerbosityLevel > 0 {
		fmt.Fprintln(os.Stderr, "Lazily running slow version of http.InternsOrThunks().")
	}
	httpNamespace.ResetMeta(MakeMeta(nil, `Provides HTTP client and server implementations.`, "1.0"))

	httpNamespace.InternVar("send", send_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("request"))),
			`Sends an HTTP request and returns an HTTP response.
  request is a map with the following keys:
  - url (string)
  - method (string, keyword or symbol, defaults to :get)
  - body (string)
  - host (string, overrides Host header if provided)
  - headers (map).
  All keys except for url are optional.
  response is a map with the following keys:
  - status (int)
  - body (string)
  - headers (map)
  - content-length (int)`, "1.0"))

	httpNamespace.InternVar("start-file-server", start_file_server_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("addr"), MakeSymbol("root"))),
			`Starts HTTP server on the TCP network address addr that
  serves HTTP requests with the contents of the file system rooted at root.`, "1.0"))

	httpNamespace.InternVar("start-server", start_server_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("addr"), MakeSymbol("handler"))),
			`Starts HTTP server on the TCP network address addr.`, "1.0"))

}
