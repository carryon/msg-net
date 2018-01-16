package jrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
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
	GetIPByPeerID(id string) string
	String() string
}

var (
	virtualP0 *peer.Peer
	// allRequest map[int64]chan *msgnetMessage
	allRequest  sync.Map
	allStatus   sync.Map
	routerInter RouterInterface
)

const (
	chainRPCMsg = 105
	nodeStatus  = 107
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
	switch msg.Cmd {
	case chainRPCMsg:
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
	case nodeStatus:
		m := make(map[string]interface{})
		json.Unmarshal(msg.Payload, m)
		m["time"] = time.Now()
		ip := m["localSever"].(map[string]interface{})["ip"].(string)
		servers := m["localSever"].(map[string]interface{})["servers"].([]interface{})
		routerInter.PeerIDIterFunc(func(peer *pb.Peer) {
			if !strings.Contains(peer.Id, "__virtual") || !strings.Contains(peer.Id, "detector") {
				if strings.Contains(routerInter.GetIPByPeerID(peer.Id), ip) {
					ids := strings.Split(peer.Id, ":")
					servers = append(servers, &ServerProcess{ChainID: ids[0], NodeID: ids[1], Status: "运行中"})
				}
			}
		})

		allStatus.Store(ip, m)
	}
	return nil
}

// RunRPCServer start rpc proxy
func RunRPCServer(port, address string, router RouterInterface) {
	virtualP0 = peer.NewPeer(fmt.Sprintf("__virtual:%d", time.Now().Nanosecond()), []string{address}, chainMessageHandle)
	virtualP0.Start()
	routerInter = router
	server := rpc.NewServer()
	server.Register(&MsgNet{Router: router, StatusTimeOut: 60 * time.Second})

	var (
		listener net.Listener
		err      error
	)
	if listener, err = net.Listen("tcp", ":"+port); err != nil {
		logger.Errorf("TestServer error %+v", err)
	}
	rpc.NewHTTPServer(server, []string{"*"}).Serve(listener)
}
