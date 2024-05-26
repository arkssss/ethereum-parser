package remotely

import (
	"bytes"
	"encoding/json"
	"errors"
	"ethereum-parser/internal/models"
	"ethereum-parser/internal/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	ReqMarshallError = errors.New("req marshall error")
	HttpRequestError = errors.New("http request error")
)

const (
	httpEthEndpoint = "https://cloudflare-eth.com/"
	httpContentType = "application/json"
	httpStatusOk    = 200
)

const (
	CurrentBlockMethod      = "eth_blockNumber"
	CurrentGetBlockByNumber = "eth_getBlockByNumber"
)

// GetCurrentBlock curl -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}'
func GetCurrentBlock() (int64, error) {
	id := requestId()
	req := reqObj{
		JsonRpc: "2.0",
		Method:  CurrentBlockMethod,
		Id:      id,
	}

	respBody, err := request(&req)
	if err != nil {
		log.Printf("getCurrentBlock resp http reqeust error: [%s], requestId: [%d]", err.Error(), id)
		return 0, HttpRequestError
	}

	resp := getBlockResp{}
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		log.Printf("getCurrentBlock resp json unmarshall error: [%s], requestId: [%d]", err.Error(), id)
		return 0, ReqMarshallError
	}

	return utils.HexToInt(resp.Result), nil
}

// GetTransactionByNumber curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x1b4", true],"id":1}'
func GetTransactionByNumber(num int64) ([]models.Transaction, error) {
	id := requestId()
	params := make([]interface{}, 2)
	params[0] = utils.IntToHex(num)
	params[1] = true
	req := reqObj{
		JsonRpc: "2.0",
		Method:  CurrentGetBlockByNumber,
		Id:      id,
		Param:   params,
	}
	respBody, err := request(&req)
	if err != nil {
		log.Printf("http reqeust error: [%s], requestId: [%d]", err.Error(), id)
		return nil, HttpRequestError
	}

	resp := getBlockByNumberResp{}
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		log.Printf("GetTransactionByNumber resp json unmarshall error: [%s], requestId: [%d]", err.Error(), id)
		return nil, ReqMarshallError
	}
	if resp.Result == nil {
		log.Printf(":GetTransactionByNumber resp not valid requestId: [%d]", id)
		return nil, ReqMarshallError
	}

	return resp.Result.Transactions, nil
}

// do http request
func request(reqBody *reqObj) ([]byte, error) {
	if reqBody == nil {
		return nil, errors.New("passed params empty")
	}

	//
	reqB, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.New("marshall request obj error")
	}

	resp, err := http.Post(httpEthEndpoint, httpContentType, strings.NewReader(string(reqB)))
	defer func(response *http.Response) {
		if response != nil {
			_ = response.Body.Close()
		}
	}(resp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("do http error: [%s]", err.Error()))
	}
	if resp == nil {
		return nil, errors.New("http resp nil")
	}
	//
	if resp.StatusCode != httpStatusOk {
		return nil, errors.New("http resp status code not ok")
	}
	respBuffer := new(bytes.Buffer)
	_, err = respBuffer.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.New("http read body error")
	}
	return respBuffer.Bytes(), nil
}

func requestId() int64 {
	return 1
}
