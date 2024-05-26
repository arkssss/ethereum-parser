package api

import (
	"encoding/json"
	"ethereum-parser/internal/domain/parser"
	"ethereum-parser/internal/models"
	"net/http"
	"strconv"
)

var (
	MessageSuccess = "success"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	RespMap := make(map[string]int)
	RespMap["currentBlock"] = parser.GetParser().GetCurrentBlock()

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(ToResponseJson(MessageSuccess, RespMap))
}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query().Get("address")
	ok := parser.GetParser().Subscribe(v)
	RespMap := make(map[string]string)
	RespMap["subscribed"] = strconv.FormatBool(ok)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(ToResponseJson(MessageSuccess, RespMap))
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query().Get("address")
	trans := parser.GetParser().GetTransactions(v)
	RespMap := make(map[string][]models.Transaction)
	RespMap["transactions"] = trans

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(ToResponseJson(MessageSuccess, RespMap))
}

func ToResponseJson(message string, data interface{}) []byte {
	r := Response{
		Message: message,
		Data:    data,
	}
	b, _ := json.Marshal(r)
	return b
}
