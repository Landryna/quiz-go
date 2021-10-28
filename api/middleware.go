package api

import (
	"context"
	"quizGO/config"

	"github.com/valyala/fasthttp"
)

// newCookie creates cookie and sets its expiration time
func newCookie(key, value string) *fasthttp.Cookie {
	c := fasthttp.Cookie{}
	c.SetKey(key)
	c.SetValue(value)
	c.SetMaxAge(config.CookieExpireSeconds)
	return &c
}

// requireTimeout requires processing of function to take time in lower manner than requestTimeOut
func requireTimeout(f fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		contextTimeOut, cancel := context.WithTimeout(context.Background(), config.RequestTimeout)
		defer cancel()

		doneCh := make(chan bool)
		go func() {
			f(ctx)
			close(doneCh)
		}()

		select {
		case <-contextTimeOut.Done():
			ctx.Error("request_timeout", fasthttp.StatusRequestTimeout)
		case <-doneCh:
		}
	}
}

// requireLogin requires that users are need to be logged in (have cookie) to use endpoint
func requireLogin(f fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		token := string(ctx.Request.Header.Cookie(config.SessionID))
		if len(token) == 0 {
			ctx.Error("unauthorized", fasthttp.StatusUnauthorized)
			return
		}
		// check if provided cookie exist in actual 'database'
		for _, u := range users {
			if u.token == token {
				ctx.SetUserValue("_user", u)
				f(ctx)
				return
			}
		}
		ctx.Error("unauthorized", fasthttp.StatusUnauthorized)
	}
}
