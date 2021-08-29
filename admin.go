package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/jkomyno/nanoid"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)


const sessionID = "sessionID"
const cookieExpireSeconds = 360

type user struct {
	login string
	password string
	token string
}

var users []user



type API interface {
	Run () error
	Close () error
}
var _ API = &adminAPI{}

type adminAPI struct {
	server fasthttp.Server
}

func newAdminAPI(s fasthttp.Server) *adminAPI {
	return &adminAPI{
		server: s,
	}
}

func (a adminAPI) Run() error {
	r := router.New()
	r.PanicHandler = func (ctx *fasthttp.RequestCtx, i interface{}) {
		ctx.Error("internal server error", fasthttp.StatusInternalServerError)
	}

	r.POST("/login", login)
	r.POST("/logout", requireLogin(logout))
	r.GET("/view", requireLogin(view))
	r.GET("/panic", requireTimeout(requireLogin(testPanic)))
	a.server.Handler = r.Handler
	return a.server.ListenAndServe(":8080")
}

func (a adminAPI) Close() error {
	return a.server.Shutdown()
}

func newCookie (key, value string) *fasthttp.Cookie {
	c := fasthttp.Cookie{}
	//c.SetDomain("admin-quiz")
	c.SetKey(key)
	c.SetValue(value)
	c.SetMaxAge(cookieExpireSeconds)
	return &c
}

func requireTimeout(f fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func (ctx *fasthttp.RequestCtx) {
		context, _ := context.WithTimeout(context.Background(), 15*time.Second)
		doneCh := make(chan bool)

		go func () {
			f(ctx)
			close(doneCh)
		}()

		select {
			case <- context.Done():
				ctx.Error("request_timeout", fasthttp.StatusRequestTimeout)
			case <- doneCh:
				return
		}
	}
}

func requireLogin(f fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		token := string(ctx.Request.Header.Cookie(sessionID))
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

func login(ctx *fasthttp.RequestCtx) {
	// nie chce miec generowanego ciastka za kazdym razem jak robie login
	// zrobic odswiezanie gdy user odwiedzi strony
	token, err := nanoid.Nanoid()
	if err != nil {
		log.Print(err)
	}

	validPostBody := struct {
		Login string
		Password string
	}{}

	postBody := ctx.PostBody()
	err = json.Unmarshal(postBody, &validPostBody); if err != nil {
		ctx.Error("bad request", fasthttp.StatusBadRequest)
		return
	}

	if validPostBody.Password == "" || validPostBody.Login == "" {
		ctx.Error("bad request", fasthttp.StatusBadRequest)
		return
	}
	// weryfikacja czy dany user juz istnieje
	for _, u := range users {
		if u.login == validPostBody.Login {
			if u.password == validPostBody.Password {
				ctx.Error("already logged in", fasthttp.StatusBadRequest)
				return
			} else {
				ctx.Error("wrong password", fasthttp.StatusUnauthorized)
				return
			}
		}
	}

	users = append(users, user{
		login:    validPostBody.Login,
		password: validPostBody.Password,
		token:    token,
	})
	ctx.Response.Header.SetCookie(newCookie(sessionID, token))
}

func logout(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.DelClientCookie(sessionID)
}

func view(ctx *fasthttp.RequestCtx) {
	user, ok := ctx.UserValue("_user").(user)
	if !ok {
		ctx.Error("bad request", fasthttp.StatusBadRequest)
		return
	}
	body := map[string]interface{}{
		fmt.Sprintf("%s with password %s", user.login, user.password):
		fmt.Sprintf("token: %s", user.token),
	}
	b, err := json.Marshal(body)
	if err != nil {
		log.Print(err)
	}
	ctx.Response.SetBody(b)
}

func testPanic(ctx *fasthttp.RequestCtx) {
	time.Sleep(10*time.Second)
	ctx.Error("SUCCESS",200)
}
