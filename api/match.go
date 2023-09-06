package api

import (
	"net/http"

	"github.com/adarsh-a-tw/tt-backend/api/dto"
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
