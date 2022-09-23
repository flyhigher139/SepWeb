package context

import (
	"encoding/json"
	"errors"
	"github.com/igevin/sepweb/pkg/template"
	"net/http"
	"net/url"
)

type Context struct {
	Req              *http.Request
	Resp             http.ResponseWriter
	RespStatusCode   int
	RespData         []byte
	MatchedRoute     string
	PathParams       map[string]string
	cacheQueryValues url.Values
	TplEngine        template.TemplateEngine
}

func (c *Context) BindJson(val any) error {
	data := c.Req.Body
	if data == nil {
		return errors.New("body is empty")
	}
	decoder := json.NewDecoder(data)
	decoder.DisallowUnknownFields()
	return decoder.Decode(val)
}

func (c *Context) FormValue(key string) StringValue {
	if err := c.Req.ParseForm(); err != nil {
		return StringValue{err: err}
	}
	return StringValue{val: c.Req.FormValue(key)}
}

func (c *Context) QueryValue(key string) StringValue {
	if c.cacheQueryValues == nil {
		c.cacheQueryValues = c.Req.URL.Query()
	}
	val, ok := c.cacheQueryValues[key]
	if !ok {
		return StringValue{err: errors.New("web: 找不到这个key")}
	}
	return StringValue{val: val[0]}
}

func (c *Context) PathValue(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}
	return StringValue{val: val}
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Resp, cookie)
}

func (c *Context) RespJSONOK(val any) error {
	return c.RespJSON(http.StatusOK, val)
}

func (c *Context) RespJSON(code int, val any) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.Resp.WriteHeader(code)
	_, err = c.Resp.Write(bs)
	return err
}

func (c *Context) Render(tpl string, data any) error {
	var err error
	c.RespData, err = c.TplEngine.Render(c.Req.Context(), tpl, data)
	c.RespStatusCode = http.StatusOK
	if err != nil {
		c.RespStatusCode = http.StatusInternalServerError
	}
	return err
}
