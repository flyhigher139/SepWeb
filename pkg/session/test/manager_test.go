package test

import (
	"github.com/google/uuid"
	sepweb "github.com/igevin/sepweb/pkg"
	"github.com/igevin/sepweb/pkg/context"
	"github.com/igevin/sepweb/pkg/handler"
	"github.com/igevin/sepweb/pkg/session"
	"github.com/igevin/sepweb/pkg/session/cookie"
	"github.com/igevin/sepweb/pkg/session/memory"
	"net/http"
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	s := sepweb.NewHttpServer()

	m := session.Manager{
		SessCtxKey: "_sess",
		Store:      memory.NewStore(30 * time.Minute),
		Propagator: cookie.NewPropagator("sessid",
			cookie.WithCookieOption(func(c *http.Cookie) {
				c.HttpOnly = true
			})),
	}

	s.Get("/login", func(ctx *context.Context) {
		// 前面就是你登录的时候一大堆的登录校验
		id := uuid.New()
		sess, err := m.InitSession(ctx, id.String())
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
		// 然后根据自己的需要设置
		err = sess.Set(ctx.Req.Context(), "mykey", "some value")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
	})
	s.Get("/resource", func(ctx *context.Context) {
		sess, err := m.GetSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
		val, err := sess.Get(ctx.Req.Context(), "mykey")
		ctx.RespData = []byte(val)
	})

	s.Get("/logout", func(ctx *context.Context) {
		_ = m.RemoveSession(ctx)
	})

	s.Use(func(next handler.Handle) handler.Handle {
		return func(ctx *context.Context) {
			// 执行校验
			if ctx.Req.URL.Path != "/login" {
				sess, err := m.GetSession(ctx)
				// 不管发生了什么错误，对于用户我们都是返回未授权
				if err != nil {
					ctx.RespStatusCode = http.StatusUnauthorized
					return
				}
				ctx.UserValues["sess"] = sess
				_ = m.Refresh(ctx.Req.Context(), sess.ID())
			}
			next(ctx)
		}
	})

	s.Start(":8081")
}
