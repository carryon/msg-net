package jrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/bocheninc/base/rpc"
	"github.com/bocheninc/msg-net/logger"
	"github.com/bocheninc/msg-net/peer"
	pb "github.com/bocheninc/msg-net/protos"
	"github.com/bocheninc/msg-net/util"
)

type RouterInterface interface {
	PeerIDIterFunc(function func(*pb.Peer))
	String() string
}

var (
	virtualP0 *peer.Peer
	// allRequest map[int64]chan *msgnetMessage
	allRequest sync.Map
)

const (
	chainRPCMsg = 105
)

func sendMessage(id int64, dst string, payload []byte) ([]byte, error) {
	fmt.Println("000000000000000:", string(payload))
	if virtualP0.Send(dst, util.Serialize(msgnetMessage{Cmd: chainRPCMsg, Payload: payload}), nil) {
		ch := make(chan *msgnetMessage)
		allRequest.Store(id, ch)
		select {
		case m := <-ch:
			return m.Payload, nil
		case <-time.NewTicker(time.Second * 30).C:
			return nil, errors.New("send msg time out")
		}
	}
	return nil, errors.New("msgnet can't send msg to l0")
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
	var result Response
	json.Unmarshal(msg.Payload, &result)
	if v, ok := allRequest.Load(result.ID); ok {
		allRequest.Delete(result.ID)
		ch, f := v.(chan *msgnetMessage)
		if f {
			ch <- msg
			close(ch)
		}
	}
	return nil
}

// RunRPCServer start rpc proxy
func RunRPCServer(port, address string, router RouterInterface) {
	virtualP0 = peer.NewPeer(fmt.Sprintf("__virtual:%d", time.Now().Nanosecond()), []string{address}, chainMessageHandle)
	virtualP0.Start()

	server := rpc.NewServer()

	server.Register(&MsgNet{Router: router})

	var (
		listener net.Listener
		err      error
	)
	if listener, err = net.Listen("tcp", ":"+port); err != nil {
		logger.Errorf("TestServer error %+v", err)
	}
	rpc.NewHTTPServer(server, []string{"*"}).Serve(listener)
}
