package router

import (
	"canada-hires/container"
	"canada-hires/controllers"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
)

type AuthRouter interface {
	Init(r chi.Router)
	InjectAuthRoutes(r chi.Router)
}

type authRouter struct {
	cn             *container.Container
	authController controllers.AuthController
}

func NewAuthRouter(cn *container.Container, authController controllers.AuthController) AuthRouter {
	return &authRouter{
		cn:             cn,
		authController: authController,
	}
}

func (ar *authRouter) InjectAuthRoutes(r chi.Router) {
	err := ar.cn.Invoke(func(authController controllers.AuthController) {
		ar := NewAuthRouter(ar.cn, authController)
		ar.Init(r)
	})

	if err != nil {
		log.Fatal("Failed to initialize routes", "error", err)
	}
}

func (ar *authRouter) Init(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {

		r.Post("/send-login-link", ar.authController.SendLoginLink)
		r.Get("/verify-login/{token}", ar.authController.VerifyLogin)
		r.Post("/logout", ar.authController.Logout)
	})
}
