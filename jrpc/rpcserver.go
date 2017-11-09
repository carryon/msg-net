package jrpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/bocheninc/msg-net/logger"
	"github.com/bocheninc/msg-net/peer"
	"github.com/bocheninc/msg-net/util"
)

var (
	virtualP0 *peer.Peer
	// allRequest map[int64]chan *msgnetMessage
	allRequest sync.Map
)

type message struct {
	ID      int64    `json:"id"`
	ChainID string   `json:"chainId"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

type msgnetMessage struct {
	Cmd     int
	Payload []byte
}

const (
	chainRPCMsg = 105
)

func parseRPCPost(w http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	// Read body
	b, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msg message

	if err = json.Unmarshal(b, &msg); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	logger.Debugf("The send request  ", string(b))

	if len(msg.ChainID) == 0 {
		http.Error(w, "should specifiy the chain id", 500)
		return
	}

	now := time.Now().UnixNano()
	output, _ := json.Marshal(struct {
		ID     int64    `json:"id"`
		Method string   `json:"method"`
		Params []string `json:"params"`
	}{now, msg.Method, msg.Params})

	w.Header().Set("content-type", "application/json")
	if virtualP0.Send(msg.ChainID, util.Serialize(msgnetMessage{Cmd: chainRPCMsg, Payload: output}), nil) {
		ch := make(chan *msgnetMessage)
		// allRequest[now] = ch
		allRequest.Store(now, ch)
		select {
		case m := <-ch:
			w.Write(m.Payload)
		case <-time.NewTicker(time.Second * 30).C:
			w.Write([]byte(`{"result","error"}`))
		}
	} else {
		w.Write([]byte(`{"result","error , msgnet send msg to l0"}`))
	}
}

func chainMessageHandle(srcID, dstID string, payload, signature []byte) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Debugln(r)
			return
		}
	}()

	msg := &msgnetMessage{}
	util.Deserialize(payload, msg)
	var reslt message
	json.Unmarshal(msg.Payload, &reslt)

	if v, ok := allRequest.Load(reslt.ID); ok {
		allRequest.Delete(reslt.ID)
		ch, f := v.(chan *msgnetMessage)
		if f {
			logger.Debugln("===================ok")
			ch <- msg
			close(ch)
		}
		logger.Debugln("===================chan err")
	}
	logger.Debugln("=======================map nil")
	// if v, ok := allRequest[reslt.ID]; ok {
	// 	v <- msg
	// 	close(v)
	// }
	return nil
}

// RunRPCServer start rpc proxy
func RunRPCServer(port, address string) {
	virtualP0 = peer.NewPeer(fmt.Sprintf("__virtual:%d", time.Now().Nanosecond()), []string{address}, chainMessageHandle)
	virtualP0.Start()
	http.HandleFunc("/", parseRPCPost)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.Errorf("Run msgnet rpc server fail %s", err.Error())
	} else {
		logger.Debug("Run msgnet rpc server")
	}
}
