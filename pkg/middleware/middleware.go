package middleware

import "github.com/igevin/sepweb/pkg/handler"

type Middleware func(next handler.Handle) handler.Handle
