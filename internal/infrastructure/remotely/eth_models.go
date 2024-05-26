package remotely

import "ethereum-parser/internal/models"

type reqObj struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Id      int64         `json:"id"`
	Param   []interface{} `json:"params"`
}

type getBlockResp struct {
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Result  string `json:"result"`
}

type getBlockByNumberResp struct {
	JsonRpc string             `json:"jsonrpc"`
	Method  string             `json:"method"`
	Result  *getTransactionObj `json:"result"`
}

type getTransactionObj struct {
	Transactions []models.Transaction `json:"transactions"`
}
