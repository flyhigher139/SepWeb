package pkg

import "net/http"

type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter
}
