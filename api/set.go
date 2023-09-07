package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/adarsh-a-tw/tt-backend/service"
	"github.com/gin-gonic/gin"
)

func (a *Api) CreateSet(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Params.ByName("match_id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := a.svc.CreateSet(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrGameOverOrSetCountExceeded) || errors.Is(err, service.ErrPreviousSetNotCompleted) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.Status(http.StatusCreated)
}

func (a *Api) UpdateScore(ctx *gin.Context) {
	matchId, err1 := strconv.Atoi(ctx.Params.ByName("match_id"))
	setId, err2 := strconv.Atoi(ctx.Params.ByName("set_id"))
	if err1 != nil || err2 != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request params"})
		return
	}

	var requestBody struct {
		ScoredByA *bool `json:"scored_by_a" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := a.svc.UpdateScore(matchId, setId, *requestBody.ScoredByA); err != nil {
		if errors.Is(err, service.ErrGameOverOrSetCountExceeded) || errors.Is(err, service.ErrPreviousSetNotCompleted) || errors.Is(err, service.ErrSetAlreadyCompleted) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.Status(http.StatusAccepted)
}
