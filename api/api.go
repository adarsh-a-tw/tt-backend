package api

import (
	"net/http"

	"github.com/adarsh-a-tw/tt-backend/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Api struct {
	svc service.Service
	r   *gin.Engine
}

func New(svc service.Service) *Api {
	r := gin.Default()
	api := &Api{svc, r}
	api.registerMiddlewares()
	api.registerEndpoints()
	return api
}

func (a *Api) registerMiddlewares() {
	a.r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
	}))
}

func (a *Api) registerEndpoints() {
	a.r.GET("/api/matches", a.GetMatchInfoList)
	a.r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	a.r.Use(adminAuthMiddleware())
	a.r.POST("/api/matches/:match_id/sets", a.CreateSet)
	a.r.POST("/api/matches/:match_id/sets/:set_id/score", a.UpdateScore)
	a.r.PATCH("/api/matches/:match_id/sets/:set_id/score", a.UndoScore)
}

func (a *Api) Serve(addr string) error {
	return a.r.Run(addr)
}
