package jrpc

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bocheninc/msg-net/logger"
	pb "github.com/bocheninc/msg-net/protos"
)

const (
	ldpQueryLog              = "Browse.GetLog"
	ldpQueryConfig           = "Browse.GetConfig"
	ldpGetTheLastBlockInfo   = "Browse.GetTheLastBlockInfo"
	ldpGetBlockByRange       = "Browse.GetBlockByRange"
	ldpGetBlockByHash        = "Browse.GetBlockByHash"
	ldpGetBlockByHeight      = "Browse.GetBlockByHeight"
	ldpGetTxByHash           = "Browse.GetTxByHash"
	ldpGetAccountInfoByID    = "Browse.GetAccountInfoByID"
	ldpGetAccountInfoByAddr  = "Browse.GetAccountInfoByAddr"
	ldpGetHistoryTransaction = "Browse.GetHistoryTransaction"
)

type MsgNet struct {
	Router        RouterInterface
	StatusTimeOut time.Duration
}

func (m *MsgNet) GetAllServerInfos(ignore string, relay *Servers) error {
	ss := &Servers{}
	m.checkStatus(func(k, v interface{}) {
		server := v.(map[string]interface{})["localServer"]
		ss.Servers = append(ss.Servers, server)
	})
	*relay = *ss
	return nil
}
func (m *MsgNet) GetServerInfoByIP(ip string, relay *interface{}) error {
	m.checkStatus(func(k, v interface{}) {
		if k.(string) == ip {
			*relay = v.(map[string]interface{})["localServer"]
		}
	})
	return nil
}
func (m *MsgNet) GetServerStatusByIP(ip string, relay *interface{}) error {
	m.checkStatus(func(k, v interface{}) {
		if k.(string) == ip {
			*relay = v.(map[string]interface{})["serverStatus"]
		}
	})
	return nil
}

func (m *MsgNet) GetNodeConfigByID(args QueryConfigArgs, relay *string) error {
	dst := fmt.Sprintf("%s:%s", args.ChainID, hex.EncodeToString([]byte(args.NodeID)))
	result, err := m.getResponse(ldpQueryConfig, dst, []string{})
	if err != nil {
		return err
	}
	*relay = result.(string)
	return nil
}

func (m *MsgNet) GetNodeLogByID(args QueryLogArgs, relay *interface{}) error {
	dst := fmt.Sprintf("%s:%s", args.ChainID, hex.EncodeToString([]byte(args.NodeID)))
	result, err := m.getResponse(ldpQueryLog, dst, [][]int{args.Range})
	if err != nil {
		return err
	}
	*relay = result
	return nil
}

func (m *MsgNet) GetAllNodeNewBlockInfo(ignore string, relay *NodesTheLastBlockInfo) error {
	nodes := &NodesTheLastBlockInfo{}
	var err error
	m.Router.PeerIDIterFunc(func(peer *pb.Peer) {
		if !strings.Contains(peer.Id, "__virtual") && !strings.Contains(peer.Id, "monitor") {
			blockInfo, err := m.getResponse(ldpGetTheLastBlockInfo, peer.Id, []string{})
			if err != nil {
				logger.Errorf("can't not get %s response %s", peer.Id, err)
				err = fmt.Errorf("can't not get %s response %s", peer.Id, err)
			} else {
				nodes.TheLastBlocks = append(nodes.TheLastBlocks, blockInfo)
			}
		}
	})
	if err != nil {
		return err
	}
	*relay = *nodes
	return nil
}

func (m *MsgNet) GetNodeNewBlockInfo(args QueryConfigArgs, relay *interface{}) error {
	dst := fmt.Sprintf("%s:%s", args.ChainID, hex.EncodeToString([]byte(args.NodeID)))
	blockInfo, err := m.getResponse(ldpGetTheLastBlockInfo, dst, []string{})
	if err != nil {
		return err
	}
	*relay = blockInfo
	return nil
}

func (m *MsgNet) GetBlockByRange(args GetBlocksInfoByRangeArgs, relay *interface{}) error {
	dst := fmt.Sprintf("%s:%s", args.ChainID, hex.EncodeToString([]byte(args.NodeID)))
	result, err := m.getResponse(ldpGetBlockByRange, dst, []int{args.Range})
	if err != nil {
		return err
	}
	*relay = result
	return nil
}

func (m *MsgNet) GetBlockByHash(args GetBlockByHashArgs, relay *interface{}) error {
	dst := fmt.Sprintf("%s:%s", args.ChainID, hex.EncodeToString([]byte(args.NodeID)))
	result, err := m.getResponse(ldpGetBlockByHash, dst, []interface{}{args.Args})
	if err != nil {
		return err
	}
	*relay = result
	return nil
}

func (m *MsgNet) GetBlockByHeight(args GetBlockByHeightArgs, relay *interface{}) error {
	dst := fmt.Sprintf("%s:%s", args.ChainID, hex.EncodeToString([]byte(args.NodeID)))
	result, err := m.getResponse(ldpGetBlockByHeight, dst, []interface{}{args.Args})
	if err != nil {
		return err
	}
	*relay = result
	return nil
}

func (m *MsgNet) GetTxByHash(args GetTxByHashArgs, relay *interface{}) error {
	dst := fmt.Sprintf("%s:%s", args.ChainID, hex.EncodeToString([]byte(args.NodeID)))
	result, err := m.getResponse(ldpGetTxByHash, dst, []string{args.Hash})
	if err != nil {
		return err
	}
	*relay = result
	return nil
}

func (m *MsgNet) GetAccountInfoByID(args GetAccountInfoByIDArgs, relay *interface{}) error {
	dst := fmt.Sprintf("%s:%s", args.ChainID, hex.EncodeToString([]byte(args.NodeID)))
	result, err := m.getResponse(ldpGetAccountInfoByID, dst, []string{args.AccountID})
	if err != nil {
		return err
	}
	*relay = result
	return nil
}

func (m *MsgNet) GetAccountInfoByAddr(args GetAccountInfoByAddrArgs, relay *interface{}) error {
	dst := fmt.Sprintf("%s:%s", args.ChainID, hex.EncodeToString([]byte(args.NodeID)))
	result, err := m.getResponse(ldpGetAccountInfoByAddr, dst, []string{args.AccountAddr})
	if err != nil {
		return err
	}
	*relay = result
	return nil
}

func (m *MsgNet) GetHistoryTransaction(args GetHistoryTransactionArgs, relay *interface{}) error {
	dst := fmt.Sprintf("%s:%s", args.ChainID, hex.EncodeToString([]byte(args.NodeID)))
	result, err := m.getResponse(ldpGetHistoryTransaction, dst, []interface{}{args.Args})
	if err != nil {
		return err
	}
	*relay = result
	return nil
}

func (m *MsgNet) getResponse(method, dst string, params interface{}) (interface{}, error) {
	now := time.Now().UnixNano()
	res := createRes(now, method, params)
	payload, err := res.byte()
	if err != nil {
		return nil, err
	}
	data, err := sendMessage(now, dst, payload)
	if err != nil {
		return nil, err
	}
	resp := &Response{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	return resp.Result, nil
}

func (m *MsgNet) checkStatus(f func(k, v interface{})) {
	allStatus.Range(func(key, value interface{}) bool {
		created := value.(map[string]interface{})["timestamp"].(time.Time)
		if created.Add(m.StatusTimeOut).Before(time.Now()) {
			logger.Debugln("remove monito :", key)
			allStatus.Delete(key)
		} else {
			f(key, value)
		}
		return true
	})
}
