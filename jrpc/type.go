package jrpc

import "encoding/json"

// type Request struct {
// 	ID      int64         `json:"id"`
// 	ChainID string        `json:"chainId"`
// 	NodeID  string        `json:"nodeId"`
// 	Method  string        `json:"method"`
// 	Params  []interface{} `json:"params"`
// }

type Response struct {
	ID     int64  `json:"id"`
	Result string `json:"result"`
	Error  string `json:"error"`
}

func (r *Response) byte() []byte {
	result, _ := json.Marshal(r)
	return result
}

func createResp(id int64, result string, err string) *Response {
	return &Response{
		ID:     id,
		Result: result,
		Error:  err,
	}
}

type msgnetMessage struct {
	Cmd     int
	Payload []byte
}

type Servers struct {
}

type QueryConfigArgs struct {
	ChainID string `json:"chainId"`
	NodeID  string `json:"nodeId"`
}
