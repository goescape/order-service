package routes

import (
	"fmt"
	handlers "order-svc/handlers/user"
	"order-svc/middlewares"
	"strings"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	Router *gin.Engine
	User   *handlers.UserHandler
}

func (r *Routes) Setup(baseURL string) {
	r.Router = gin.New()
	r.Router.Use(middlewares.EnabledCORS(), middlewares.Logger(r.Router))

	if baseURL != "" && baseURL != "/" {
		baseURL = "/" + strings.Trim(baseURL, "/")
	} else {
		baseURL = "/"
	}

	r.setupAPIRoutes(baseURL)
}

func (r *Routes) setupAPIRoutes(baseURL string) {
	apiGroup := r.Router.Group(baseURL)
	r.configureUserRoutes(apiGroup)
}

func (r *Routes) configureUserRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/user")
	userGroup.POST("/register", r.User.HandleUserRegister)
	userGroup.POST("/login", r.User.HandleUserLogin)
}

func (r *Routes) configureOrderRoutes(router *gin.RouterGroup) {
	orderGroup := router.Group("/order")

}

func (r *Routes) Run(port string) {
	if r.Router == nil {
		panic("[ROUTER ERROR] Gin Engine has not been initialized. Make sure to call Setup() before Run().")
	}

	err := r.Router.Run(":" + port)
	if err != nil {
		panic(fmt.Sprintf("[SERVER ERROR] Failed to start the server on port %s: %v", port, err))
	}
}
