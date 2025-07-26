package router

import (
	"canada-hires/container"
	"canada-hires/controllers"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
)

type userRouter struct {
	cn             *container.Container
	userController controllers.UserController
	authMW         func(http.Handler) http.Handler
	requireMW      func(http.Handler) http.Handler
}

func NewUserRouter(cn *container.Container, userController controllers.UserController, authMW func(http.Handler) http.Handler, requireMW func(http.Handler) http.Handler) AuthRouter {
	return &userRouter{
		cn:             cn,
		userController: userController,
		authMW:         authMW,
		requireMW:      requireMW,
	}
}

func (ar *userRouter) InjectAuthRoutes(r chi.Router) {
	err := ar.cn.Invoke(func(userController controllers.UserController) {
		ar := NewUserRouter(ar.cn, userController, ar.authMW, ar.requireMW)
		ar.Init(r)
	})

	if err != nil {
		log.Fatal("Failed to initialize routes", "error", err)
	}
}

func (ar *userRouter) Init(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Use(ar.authMW)
		r.With(ar.requireMW).Get("/profile", ar.userController.GetUser)
	})
}
