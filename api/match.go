package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/adarsh-a-tw/tt-backend/api/dto"
	"github.com/adarsh-a-tw/tt-backend/service"
	"github.com/gin-gonic/gin"
)

func (a *Api) GetMatchInfoList(ctx *gin.Context) {
	var queryParams struct {
		Filter string `form:"status" binding:"oneof=ONGOING UPCOMING PAST"`
	}

	if err := ctx.ShouldBindQuery(&queryParams); err != nil && queryParams.Filter != "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid status, choices are ONGOING, UPCOMING & PAST"})
		return
	}

	matchInfoList, err := a.svc.GetMatchInfoList(queryParams.Filter)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	matchInfo := make([]dto.MatchInfoResponse, 0, len(matchInfoList))
	for _, mi := range matchInfoList {
		opponents := make([]dto.OpponentResponse, 2)
		for i, opp := range mi.Opponents {
			opponents[i] = dto.OpponentResponse{Name: opp.Name, IsWinner: opp.IsWinner}
		}
		matchInfo = append(matchInfo, dto.MatchInfoResponse{
			Id:        mi.Id,
			Format:    string(mi.Format),
			Status:    string(mi.Status),
			Stage:     string(mi.Stage),
			Opponents: opponents,
		})
	}

	response := gin.H{"matches": matchInfo}

	ctx.JSON(http.StatusOK, response)
}

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
