package api

import (
	"quizGO/config"

	"github.com/fasthttp/router"
	"github.com/jkomyno/nanoid"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type user struct {
	login    string
	password string
	token    string
}

var users []user

var _ Service = &adminService{}

type adminService struct {
	logger *logrus.Logger
}

func NewAdminService(log *logrus.Logger) *adminService {
	return &adminService{
		logger: log,
	}
}

func (a *adminService) Route(r *router.Router) *router.Router {
	r.POST("/api/v1/login", requireTimeout(a.login))
	r.POST("/api/v1/register", requireTimeout(a.register))
	r.POST("/api/v1/logout", requireTimeout(a.logout))
	return r
}

func (a *adminService) login(ctx *fasthttp.RequestCtx) {
	loginObj := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}

	if err := decode(ctx.PostBody(), &loginObj); err != nil {
		a.logger.WithError(err).Error("login: invalid payload")
		ctx.Error("invalid payload", fasthttp.StatusBadRequest)
		return
	}

	for _, u := range users {
		if u.login == loginObj.Login {
			if u.password == loginObj.Password {
				ctx.Response.Header.SetCookie(newCookie(config.SessionID, u.token))
				ctx.SetStatusCode(fasthttp.StatusOK)
				return
			} else {
				ctx.Error("invalid password", fasthttp.StatusUnauthorized)
				return
			}
		}
	}
	ctx.Error("unauthorized", fasthttp.StatusUnauthorized)
}

func (a *adminService) register(ctx *fasthttp.RequestCtx) {
	registerObj := struct {
		Login         string `json:"login"`
		Password      string `json:"password"`
		RetryPassword string `json:"retry_password"`
	}{}

	if err := decode(ctx.PostBody(), &registerObj); err != nil {
		a.logger.WithError(err).Error("register: invalid payload")
		ctx.Error("invalid payload", fasthttp.StatusBadRequest)
		return
	}

	if registerObj.Password != registerObj.RetryPassword {
		a.logger.Error("register: password mismatch")
		ctx.Error("password mismatch", fasthttp.StatusUnauthorized)
		return
	}

	for _, u := range users {
		if u.login == registerObj.Login {
			a.logger.Error("register: user already exists")
			ctx.Error("user already exists", fasthttp.StatusUnauthorized)
			return
		}
	}

	token, err := nanoid.Nanoid()
	if err != nil {
		a.logger.WithError(err).Error("register: nanoid")
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}

	u := user{
		login:    registerObj.Login,
		password: registerObj.Password,
		token:    token,
	}

	users = append(users, u)
	ctx.Response.Header.SetCookie(newCookie(config.SessionID, token))
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (a *adminService) logout(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.DelClientCookie(config.SessionID)
}
