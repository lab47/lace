package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	. "github.com/candid82/joker/core"
)

var client = &http.Client{}

func extractMethod(request Map) string {
	if ok, m := request.Get(MakeKeyword("method")); ok {
		switch m := m.(type) {
		case String:
			return m.S
		case Keyword:
			return m.ToString(false)[1:]
		case Symbol:
			return m.ToString(false)
		default:
			panic(StubNewError(fmt.Sprintf("method must be a string, keyword or symbol, got %s", m.GetType().ToString(false))))
		}
	}
	return "get"
}

func getOrPanic(m Map, k Object, errMsg string) Object {
	if ok, v := m.Get(k); ok {
		return v
	}
	panic(StubNewError(errMsg))
}

func mapToReq(env *Env, request Map) (*http.Request, error) {
	method := strings.ToUpper(extractMethod(request))
	urlv, err := AssertString(env, getOrPanic(request, MakeKeyword("url"), ":url key must be present in request map"), "url must be a string")
	if err != nil {
		return nil, err
	}

	url := urlv.S
	var reqBody io.Reader
	if ok, b := request.Get(MakeKeyword("body")); ok {
		sv, err := AssertString(env, b, "body must be a string")
		if err != nil {
			return nil, err
		}
		reqBody = strings.NewReader(sv.S)
	}
	req, err := http.NewRequest(method, url, reqBody)
	PanicOnErr(err)
	if ok, headers := request.Get(MakeKeyword("headers")); ok {
		h, err := AssertMap(env, headers, "headers must be a map")
		if err != nil {
			return nil, err
		}
		for iter := h.Iter(); iter.HasNext(); {
			p := iter.Next()
			key, err := AssertString(env, p.Key, "header name must be a string")
			if err != nil {
				return nil, err
			}

			val, err := AssertString(env, p.Value, "header value must be a string")
			if err != nil {
				return nil, err
			}
			req.Header.Add(key.S, val.S)
		}
	}
	if ok, host := request.Get(MakeKeyword("host")); ok {
		s, err := AssertString(env, host, "host must be a string")
		if err != nil {
			return nil, err
		}
		req.Host = s.S
	}
	return req, nil
}

func reqToMap(host String, port String, req *http.Request) (Map, error) {
	defer req.Body.Close()
	res := EmptyArrayMap()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	res.Add(MakeKeyword("request-method"), MakeKeyword(strings.ToLower(req.Method)))
	res.Add(MakeKeyword("body"), MakeString(string(body)))
	res.Add(MakeKeyword("uri"), MakeString(req.URL.Path))
	res.Add(MakeKeyword("query-string"), MakeString(req.URL.RawQuery))
	res.Add(MakeKeyword("server-name"), host)
	res.Add(MakeKeyword("server-port"), port)
	res.Add(MakeKeyword("remote-addr"), MakeString(req.RemoteAddr[:strings.LastIndexByte(req.RemoteAddr, byte(':'))]))
	res.Add(MakeKeyword("protocol"), MakeString(req.Proto))
	res.Add(MakeKeyword("scheme"), MakeKeyword("http"))
	headers := EmptyArrayMap()
	for k, v := range req.Header {
		headers.Add(MakeString(strings.ToLower(k)), MakeString(strings.Join(v, ",")))
	}
	res.Add(MakeKeyword("headers"), headers)
	return res, nil
}

func respToMap(resp *http.Response) (Map, error) {
	defer resp.Body.Close()
	res := EmptyArrayMap()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res.Add(MakeKeyword("body"), MakeString(string(body)))
	res.Add(MakeKeyword("status"), MakeInt(resp.StatusCode))
	respHeaders := EmptyArrayMap()
	for k, v := range resp.Header {
		respHeaders.Add(MakeString(k), MakeString(strings.Join(v, ",")))
	}
	res.Add(MakeKeyword("headers"), respHeaders)
	// TODO: 32-bit issue
	res.Add(MakeKeyword("content-length"), MakeInt(int(resp.ContentLength)))
	return res, nil
}

func mapToResp(env *Env, response Map, w http.ResponseWriter) error {
	status := 0
	if ok, s := response.Get(MakeKeyword("status")); ok {
		sv, err := AssertInt(env, s, "HTTP response status must be an integer")
		if err != nil {
			return err
		}
		status = sv.I
	}
	body := ""
	if ok, b := response.Get(MakeKeyword("body")); ok {
		bodyv, err := AssertString(env, b, "HTTP response body must be a string")
		if err != nil {
			return err
		}

		body = bodyv.S
	}
	if ok, headers := response.Get(MakeKeyword("headers")); ok {
		header := w.Header()
		h, err := AssertMap(env, headers, "HTTP response headers must be a map")
		if err != nil {
			return err
		}
		for iter := h.Iter(); iter.HasNext(); {
			p := iter.Next()
			hnamev, err := AssertString(env, p.Key, "HTTP response header name must be a string")
			if err != nil {
				return err
			}
			hname := hnamev.S
			switch pvalue := p.Value.(type) {
			case String:
				header.Add(hname, pvalue.S)
			case Seqable:
				s := pvalue.Seq()
				for !s.IsEmpty() {
					x, err := AssertString(env, s.First(), "HTTP response header value must be a string")
					if err != nil {
						return err
					}
					header.Add(hname, x.S)
					s = s.Rest()
				}
			default:
				panic(StubNewError("HTTP response header value must be a string or a seq of strings"))
			}
		}
	}
	if status != 0 {
		w.WriteHeader(status)
	}
	io.WriteString(w, body)
	return nil
}

func sendRequest(env *Env, request Map) (Map, error) {
	req, err := mapToReq(env, request)
	if err != nil {
		return nil, err
	}
	//RT.GIL.Unlock()
	resp, err := client.Do(req)
	//RT.GIL.Lock()
	PanicOnErr(err)
	return respToMap(resp)
}

func startServer(env *Env, addr string, handler Callable) (Object, error) {
	i := strings.LastIndexByte(addr, byte(':'))
	host, port := MakeString(addr), MakeString("")
	if i != -1 {
		host = MakeString(addr[:i])
		port = MakeString(addr[i+1:])
	}
	//RT.GIL.Unlock()
	//defer RT.GIL.Lock()
	err := http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//RT.GIL.Lock()
		defer func() {
			//RT.GIL.Unlock()
			if r := recover(); r != nil {
				w.WriteHeader(500)
				io.WriteString(w, "Internal server error")
				fmt.Fprintln(os.Stderr, r)
			}
		}()
		m, err := reqToMap(host, port, req)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return
		}
		response, err := handler.Call(env, []Object{m})
		rm, err := AssertMap(env, response, "HTTP response must be a map")
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return
		}
		mapToResp(env, rm, w)
	}))
	PanicOnErr(err)
	return NIL, nil
}

func startFileServer(addr string, root string) (Object, error) {
	err := http.ListenAndServe(addr, http.FileServer(http.Dir(root)))
	return NIL, err
}
