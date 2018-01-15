package jrpc

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bocheninc/msg-net/logger"
	pb "github.com/bocheninc/msg-net/protos"
)

const (
	ldpQueryLog              = "Ldp.GetLog"
	ldpQueryConfig           = "Ldp.GetConfig"
	ldpGetTheLastBlockInfo   = "Ldp.GetTheLastBlockInfo"
	ldpGetBlockByRange       = "Ldp.GetBlockByRange"
	ldpGetBlockByHash        = "Ldp.GetBlockByHash"
	ldpGetBlockByHeight      = "Ldp.GetBlockByHeight"
	ldpGetTxByHash           = "Ldp.GetTxByHash"
	ldpGetAccountInfoByID    = "Ldp.GetAccountInfoByID"
	ldpGetAccountInfoByAddr  = "Ldp.GetAccountInfoByAddr"
	ldpGetHistoryTransaction = "Ldp.GetHistoryTransaction"
)

type MsgNet struct {
	Router RouterInterface
}

func (m *MsgNet) GetAllServers(key string, relay *Servers) error       { return nil }
func (m *MsgNet) GetServersByIP(key string, relay *Servers) error      { return nil }
func (m *MsgNet) GetServerStatusByIP(key string, relay *Servers) error { return nil }

func (m *MsgNet) GetNodeConfigByID(args QueryConfigArgs, relay *string) error {
	dst := fmt.Sprintf("%s:%s", args.ChainID, hex.EncodeToString([]byte(args.NodeID)))
	result, err := m.getResponse(ldpQueryConfig, dst, []string{})
	if err != nil {
		return err
	}
	*relay = result.(string)
	return nil
}

func (m *MsgNet) GetNodeLogByID(args QueryLogArgs, relay *string) error {
	dst := fmt.Sprintf("%s:%s", args.ChainID, hex.EncodeToString([]byte(args.NodeID)))
	result, err := m.getResponse(ldpQueryLog, dst, [][]int{args.Range})
	if err != nil {
		return err
	}
	*relay = result.(string)
	return nil
}

func (m *MsgNet) GetAllNodeNewBlockInfo(ignore string, relay *NodesTheLastBlockInfo) error {
	nodes := &NodesTheLastBlockInfo{}
	var err error
	m.Router.PeerIDIterFunc(func(peer *pb.Peer) {
		//TODO 所有节点，但是不包括监控节点
		blockInfo, err := m.getResponse(ldpGetTheLastBlockInfo, peer.Id, []string{})
		if err != nil {
			logger.Errorf("can't not get %s response %s", peer.Id, err)
			err = fmt.Errorf("can't not get %s response %s", peer.Id, err)
		} else {
			nodes.TheLastBlocks = append(nodes.TheLastBlocks, blockInfo)
		}
	})
	if err != nil {
		return err
	}

	*relay = *nodes
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
	result, err := m.getResponse(ldpGetBlockByHash, dst, []string{args.Hash})
	if err != nil {
		return err
	}
	*relay = result
	return nil
}

func (m *MsgNet) GetBlockByHeight(args GetBlockByHeightArgs, relay *interface{}) error {
	dst := fmt.Sprintf("%s:%s", args.ChainID, hex.EncodeToString([]byte(args.NodeID)))
	result, err := m.getResponse(ldpGetBlockByHeight, dst, []int{args.Height})
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
	fmt.Println("111111111111:", args)
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
