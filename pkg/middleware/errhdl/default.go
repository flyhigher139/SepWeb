package errhdl

import (
	"bytes"
	"github.com/igevin/sepweb/pkg/middleware"
	"html/template"
	"log"
	"net/http"
)

func createResp(name, page string) []byte {
	tpl, err := template.New(name).Parse(page)
	if err != nil {
		log.Fatal(err)
	}
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, nil)
	if err != nil {
		log.Fatal(err)
	}
	return buffer.Bytes()
}

func createResp404() []byte {
	page := `
<html>
	<h1>404 NOT FOUND</h1>
</html>
`
	return createResp("404", page)
}

func createResp500() []byte {
	page := `
<html>
	<h1>500 Internal Server Error</h1>
</html>
`
	return createResp("500", page)
}

func CreateHttpErrorHandleMiddleware() middleware.Middleware {
	return NewMiddlewareBuilder().
		RegisterError(http.StatusNotFound, createResp404()).
		RegisterError(http.StatusInternalServerError, createResp500()).
		Build()
}
