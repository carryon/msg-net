package jrpc

import (
	"encoding/json"
)

type Request struct {
	ID     int64       `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

func createRes(id int64, method string, params interface{}) *Request {
	return &Request{
		ID:     id,
		Method: method,
		Params: params,
	}
}
func (res *Request) byte() ([]byte, error) {
	return json.Marshal(res)
}

type Response struct {
	ID     int64       `json:"id"`
	Result interface{} `json:"result"`
	Error  string      `json:"error"`
}

// func (r *Response) byte() []byte {
// 	result, _ := json.Marshal(r)
// 	return result
// }

// func createResp(id int64, result string, err string) *Response {
// 	return &Response{
// 		ID:     id,
// 		Result: result,
// 		Error:  err,
// 	}
// }

type NodesTheLastBlockInfo struct {
	TheLastBlocks []interface{} `json:"theLastBlocks"`
}

type msgnetMessage struct {
	Cmd     int
	Payload []byte
}

type Servers struct {
}

type QueryLogArgs struct {
	ChainID string
	NodeID  string
	Range   []int
}

type QueryConfigArgs struct {
	ChainID string
	NodeID  string
}

type GetBlocksInfoByRangeArgs struct {
	ChainID string
	NodeID  string
	Range   int
}

type GetBlockByHashArgs struct {
	ChainID string
	NodeID  string
	Hash    string
}

type GetBlockByHeightArgs struct {
	ChainID string
	NodeID  string
	Height  int
}

type GetTxByHashArgs struct {
	ChainID string
	NodeID  string
	Hash    string
}

type GetAccountInfoByIDArgs struct {
	ChainID   string
	NodeID    string
	AccountID string
}

type GetAccountInfoByAddrArgs struct {
	ChainID     string
	NodeID      string
	AccountAddr string
}

type GetHistoryTransactionArgs struct {
	ChainID string
	NodeID  string
	Args    *Params
}

type Params struct {
	Addr  string
	Range []int
}
