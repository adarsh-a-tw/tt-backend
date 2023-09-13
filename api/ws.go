package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adarsh-a-tw/tt-backend/api/dto"
	"github.com/adarsh-a-tw/tt-backend/service"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]*int)

func handleClient(c *websocket.Conn, svc service.Service) {
	defer func() {
		delete(clients, c)
		log.Println("Closing Websocket")
		c.Close()
	}()
	clients[c] = nil

	var m dto.MatchSubscribeRequest
	err := c.ReadJSON(&m)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("[err] client error: %v", err)
		}
		return
	}

	clients[c] = &m.MatchId

	var resp dto.MatchDetailResponse
	md, err := svc.GetMatchDetails(m.MatchId)
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp = NewMatchDetailsResponse(md)
	}
	notifyMatch(c, svc, resp)

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

func NewMatchDetailsResponse(md *service.MatchDetail) dto.MatchDetailResponse {
	var resp dto.MatchDetailResponse
	opponents := make([]dto.OpponentResponse, 0)
	for _, opp := range md.Opponents {
		opponents = append(opponents, dto.OpponentResponse{
			Id:       opp.Id,
			Name:     opp.Name,
			IsWinner: opp.IsWinner,
		})
	}

	sets := make([]dto.SetResponse, 0)
	for _, s := range md.Sets {
		setLogs := make([]dto.SetLogResponse, 0)
		for _, sl := range s.Logs {
			setLogs = append(setLogs, dto.SetLogResponse{
				Id:        sl.Id,
				OppAScore: sl.OppAScore,
				OppBScore: sl.OppBScore,
				ScoredByA: sl.ScoredByA,
			})
		}

		sets = append(sets, dto.SetResponse{
			Id:             s.Id,
			SetNumber:      s.SetNumber,
			OpponentAScore: s.OpponentAScore,
			OpponentBScore: s.OpponentBScore,
			IsCompleted:    s.IsCompleted,
			Logs:           setLogs,
		})
	}

	resp.Data = dto.MatchDetail{
		Id:        md.Id,
		Format:    md.Format,
		Stage:     md.Stage,
		Status:    md.Status,
		Opponents: opponents,
		Sets:      sets,
	}
	return resp
}

func notifyMatch(c *websocket.Conn, svc service.Service, resp dto.MatchDetailResponse) {
	log.Printf("Notifying client for match id: %d\n", resp.Data.Id)
	wr, err := c.NextWriter(websocket.TextMessage)
	defer func() {
		_ = wr.Close()
	}()

	if err != nil {
		log.Println("[err] creating writer", err)
		return
	}

	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("[err] creating match response", err)
		return
	}

	_, err = wr.Write(jsonBytes)
	if err != nil {
		log.Println("[err] writing match response", err)
	}
}

func NotifySubscribers(matchId int, svc service.Service) {
	var resp dto.MatchDetailResponse
	md, err := svc.GetMatchDetails(matchId)
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp = NewMatchDetailsResponse(md)
	}

	for c, mId := range clients {
		if mId != nil && *mId == matchId {
			go notifyMatch(c, svc, resp)
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request, svc service.Service) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go handleClient(conn, svc)
}
