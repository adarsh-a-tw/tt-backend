package api

import (
	"net/http"

	"github.com/adarsh-a-tw/tt-backend/service"
	"github.com/gin-gonic/gin"
)

type Api struct {
	svc service.Service
	r   *gin.Engine
}

func New(svc service.Service) *Api {
	r := gin.Default()
	api := &Api{svc, r}
	api.registerEndpoints()
	return api
}

func (a *Api) registerEndpoints() {
	a.r.GET("/api/matches", a.GetMatchInfoList)

	a.r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}

func (a *Api) Serve(addr string) error {
	return a.r.Run(addr)
}
