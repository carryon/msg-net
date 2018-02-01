package jrpc

import (
	"encoding/hex"
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
	RouterIterFunc(function func(string, *pb.Router))
}

type Jrpc struct {
	virtualP0                     *peer.Peer
	allRequest, allStatus         sync.Map
	routerInter                   RouterInterface
	statusTimeOut                 time.Duration
	ldpIDs, msgnetIDs, monitorIDs []string
}

const (
	chainRPCMsg   = 105
	nodeStatusMsg = 107
)

func (j *Jrpc) sendMessage(id int64, dst string, payload []byte) ([]byte, error) {
	fmt.Println("000000000000000:", string(payload))
	if j.virtualP0.Send(dst, util.Serialize(msgnetMessage{Cmd: chainRPCMsg, Payload: payload}), nil) {
		ch := make(chan *msgnetMessage)
		j.allRequest.Store(id, ch)
		select {
		case m := <-ch:
			return m.Payload, nil
		case <-time.NewTicker(time.Second * 30).C:
			return nil, errors.New("send msg time out")
		}
	}
	return nil, errors.New("msgnet can't send msg to l0")
}

func (j *Jrpc) chainMessageHandle(srcID, dstID string, payload, signature []byte) error {
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
		if v, ok := j.allRequest.Load(result.ID); ok {
			j.allRequest.Delete(result.ID)
			ch, f := v.(chan *msgnetMessage)
			if f {
				ch <- msg
				close(ch)
			}
		}
	case nodeStatusMsg:

		//initialization process
		confirmedLdps := make(map[string]bool)
		for _, v := range j.ldpIDs {
			confirmedLdps[v] = false
		}

		confirmedMsgnets := make(map[string]bool)
		for _, v := range j.msgnetIDs {
			confirmedMsgnets[v] = false
		}

		confirmedMonitor := make(map[string]bool)
		for _, v := range j.monitorIDs {
			confirmedMonitor[v] = false
		}

		//processing msg
		m := make(map[string]interface{})
		json.Unmarshal(msg.Payload, &m)
		m["timestamp"] = time.Now()
		localSever := m["localServer"].(map[string]interface{})
		ip := localSever["ip"].(string)
		servers := []*ServerProcess{}
		j.routerInter.PeerIDIterFunc(func(peer *pb.Peer) {
			if !strings.Contains(peer.Id, "__virtual") {
				if strings.Contains(j.routerInter.GetIPByPeerID(peer.Id), ip) {
					ids := strings.Split(peer.Id, ":")
					var confirmed bool
					if strings.Contains(peer.Id, "monitor") {
						_, ok := confirmedMonitor[ids[1]]
						if ok {
							confirmed = true
							delete(confirmedMonitor, ids[1])
						}
						servers = append(servers, &ServerProcess{ChainID: "", NodeID: ids[1], Status: "运行中", Type: monitor, Confirmed: confirmed})
					} else {
						nodeID, _ := hex.DecodeString(ids[1])
						_, ok := confirmedLdps[string(nodeID)]
						if ok {
							confirmed = true
							delete(confirmedLdps, string(nodeID))
						}
						servers = append(servers, &ServerProcess{ChainID: ids[0], NodeID: string(nodeID), Status: "运行中", Type: ldpPeer, Confirmed: confirmed})
					}
				}
			}
		})
		j.routerInter.RouterIterFunc(func(key string, router *pb.Router) {
			fmt.Println("-------ssssssss->>>>>>>", key, router.Id, router.Address, ip)
			var confirmed bool
			if strings.Contains(router.Address, ip) {
				_, ok := confirmedMsgnets[router.Id]
				if ok {
					confirmed = true
					delete(confirmedMsgnets, router.Id)

				}
				servers = append(servers, &ServerProcess{ChainID: "", NodeID: router.Id, Status: "运行中", Type: msgnet, Confirmed: confirmed})
			}
		})

		//add confirmed peer
		for k := range confirmedLdps {
			servers = append(servers, &ServerProcess{ChainID: "", NodeID: k, Status: "故障", Type: ldpPeer, Confirmed: true})
		}
		for k := range confirmedMsgnets {
			servers = append(servers, &ServerProcess{ChainID: "", NodeID: k, Status: "故障", Type: msgnet, Confirmed: true})
		}
		for k := range confirmedMonitor {
			servers = append(servers, &ServerProcess{ChainID: "", NodeID: k, Status: "故障", Type: monitor, Confirmed: true})
		}

		localSever["servers"] = servers
		j.allStatus.Store(ip, m)
	}
	return nil
}

// RunRPCServer start rpc proxy
func (j *Jrpc) RunRPCServer(port, address string, router RouterInterface, ldpIDs, msgnetIDs, monitorIDs []string) {
	j.ldpIDs = ldpIDs
	j.msgnetIDs = msgnetIDs
	j.monitorIDs = monitorIDs
	j.virtualP0 = peer.NewPeer(fmt.Sprintf("__virtual:%d", time.Now().Nanosecond()), []string{address}, j.chainMessageHandle)
	j.virtualP0.Start()
	j.routerInter = router
	j.statusTimeOut = 60 * time.Second
	server := rpc.NewServer()
	server.Register(&MsgNet{jrpc: j})
	var (
		listener net.Listener
		err      error
	)
	if listener, err = net.Listen("tcp", ":"+port); err != nil {
		logger.Errorf("TestServer error %+v", err)
	}
	rpc.NewHTTPServer(server, []string{"*"}).Serve(listener)
}
