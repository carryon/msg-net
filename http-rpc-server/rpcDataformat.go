package http_rpc_server

import (
	"encoding/json"

	"github.com/bocheninc/msg-net/logger"
	"github.com/bocheninc/msg-net/util"
	"github.com/twinj/uuid"
)

const (
	ChainRpcMsg = iota + 105
)

const (
	RpcNewAccount = iota + 10
	RpcListAccount
	RpcSignAccount
	RpcTransCreate
	RpcTransBroadcast
	RpcTransBroadQuery
	NetGetLocalpeer
	NetGetPeers
	LedgerBalance
	LedgerHeight
	LastBlockHash
	NumberForBlockHash
	NumberForBlock
	HashForBlock
	HashForTx
	BlockHashForTxs
	BlockNumberForTxs
	MergeTxHashForTxs
)

func transformResponseData(responseData []byte, httpoutput *SerialDataType) {
	encapdata := &sendMessage{}
	util.Deserialize(responseData, encapdata)
	switch encapdata.Cmd {
	case RpcNewAccount:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr  string
			Outputaddr Address
		}{}
		util.Deserialize(tmpdata, &outputdata)
		logger.Debugf("The address is %s", outputdata.Outputaddr.String())
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result string  `json:"result"`
				Err    *string `json:"error"`
			}{
				outputdata.Outputaddr.String(), nil,
			})
		}
	case RpcListAccount:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr       string
			Outaccountslist []string
		}{}
		util.Deserialize(tmpdata, &outputdata)
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result []string `json:"result"`
				Err    *string  `json:"error"`
			}{
				outputdata.Outaccountslist, nil,
			})
		}
	case RpcSignAccount:
		tmpdata := encapdata.Payload
		outsigndata := struct {
			Outputerr string
			Signstr   string
		}{}
		util.Deserialize(tmpdata, &outsigndata)
		if len(outsigndata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outsigndata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result string  `json:"result"`
				Err    *string `json:"error"`
			}{
				outsigndata.Signstr, nil,
			})
		}
	case RpcTransCreate:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr string
			Outtrans  string
		}{}
		util.Deserialize(tmpdata, &outputdata)
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result string  `json:"result"`
				Err    *string `json:"error"`
			}{
				outputdata.Outtrans, nil,
			})
		}
	case RpcTransBroadcast:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr    string
			Outbroadcast BroadcastReply
		}{}
		util.Deserialize(tmpdata, &outputdata)
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result BroadcastReply `json:"result"`
				Err    *string        `json:"error"`
			}{
				outputdata.Outbroadcast, nil,
			})
		}
	case RpcTransBroadQuery:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr     string
			Outbroadquery string
		}{}
		util.Deserialize(tmpdata, &outputdata)
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result string  `json:"result"`
				Err    *string `json:"error"`
			}{
				outputdata.Outbroadquery, nil,
			})
		}
	case NetGetLocalpeer:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr    string
			Outlocalpeer string
		}{}
		util.Deserialize(tmpdata, &outputdata)
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result string  `json:"result"`
				Err    *string `json:"error"`
			}{
				outputdata.Outlocalpeer, nil,
			})
		}
	case NetGetPeers:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr string
			Outpeers  []string
		}{}
		util.Deserialize(tmpdata, &outputdata)
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result []string `json:"result"`
				Err    *string  `json:"error"`
			}{
				outputdata.Outpeers, nil,
			})
		}
	case LedgerBalance:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr  string
			Outbalance Balance
		}{}
		util.Deserialize(tmpdata, &outputdata)
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result Balance `json:"result"`
				Err    *string `json:"error"`
			}{
				outputdata.Outbalance, nil,
			})
		}
	case LedgerHeight:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr string
			Outheight uint32
		}{}
		util.Deserialize(tmpdata, &outputdata)
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result uint32  `json:"result"`
				Err    *string `json:"error"`
			}{
				outputdata.Outheight, nil,
			})
		}
	case LastBlockHash:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr string
			Outhash   Hash
		}{}
		util.Deserialize(tmpdata, &outputdata)
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result string  `json:"result"`
				Err    *string `json:"error"`
			}{
				outputdata.Outhash.String(), nil,
			})
		}
	case NumberForBlockHash:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr    string
			OutBlockhash Hash
		}{}
		util.Deserialize(tmpdata, &outputdata)
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result string  `json:"result"`
				Err    *string `json:"error"`
			}{
				outputdata.OutBlockhash.String(), nil,
			})
		}
	case NumberForBlock:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr string
			Outblock  Block
		}{}
		util.Deserialize(tmpdata, &outputdata)
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result Block   `json:"result"`
				Err    *string `json:"error"`
			}{
				outputdata.Outblock, nil,
			})
		}
	case HashForBlock:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr string
			Outblock  Block
		}{}
		util.Deserialize(tmpdata, &outputdata)
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result Block   `json:"result"`
				Err    *string `json:"error"`
			}{
				outputdata.Outblock, nil,
			})
		}
	case HashForTx:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr string
			Outtx     Transaction
		}{}
		util.Deserialize(tmpdata, &outputdata)
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result Transaction `json:"result"`
				Err    *string     `json:"error"`
			}{
				outputdata.Outtx, nil,
			})
		}
	case BlockHashForTxs:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr string
			OutTxs    Transactions
		}{}
		util.Deserialize(tmpdata, &outputdata)
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result Transactions `json:"result"`
				Err    *string      `json:"error"`
			}{
				outputdata.OutTxs, nil,
			})
		}
	case BlockNumberForTxs:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr string
			OutTxs    Transactions
		}{}
		util.Deserialize(tmpdata, &outputdata)
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result Transactions `json:"result"`
				Err    *string      `json:"error"`
			}{
				outputdata.OutTxs, nil,
			})
		}
	case MergeTxHashForTxs:
		tmpdata := encapdata.Payload
		outputdata := struct {
			Outputerr string
			OutTxs    Transactions
		}{}
		util.Deserialize(tmpdata, &outputdata)
		if len(outputdata.Outputerr) > 0 {
			*httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				outputdata.Outputerr,
			})
		} else {
			*httpoutput, _ = json.Marshal(struct {
				Result Transactions `json:"result"`
				Err    *string      `json:"error"`
			}{
				outputdata.OutTxs, nil,
			})
		}
	default:
		*httpoutput, _ = json.Marshal(struct {
			Err string `json:"error"`
		}{
			"Unknown message for response http request",
		})
		logger.Debugf("Unknown message when transform response http request")
	}

	return
}

func recontructData(msgtype uint8, payload []byte) (uuid.Uuid, []byte) {
	var (
		uid              uuid.Uuid
		retSerializeData []byte
	)
	switch msgtype {
	case RpcNewAccount:
		retdata := struct {
			Outaddress Address
			Outputerr  string
			IdentityId uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		logger.Debugf("The return address is %s", retdata.Outaddress.String())
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr  string
				Outaddress Address
			}{
				retdata.Outputerr,
				retdata.Outaddress,
			}),
		}
		retSerializeData = util.Serialize(encapdata)
	case RpcListAccount:
		retdata := struct {
			Outaccountslist []string
			Outputerr       string
			IdentityId      uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr       string
				Outaccountslist []string
			}{retdata.Outputerr, retdata.Outaccountslist}),
		}
		retSerializeData = util.Serialize(encapdata)
	case RpcSignAccount:
		retdata := struct {
			Outsignstr string
			Outputerr  string
			IdentityId uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr  string
				Outsignstr string
			}{retdata.Outputerr, retdata.Outsignstr}),
		}
		retSerializeData = util.Serialize(encapdata)
	case RpcTransCreate:
		retdata := struct {
			Outtransactionout string
			Outputerr         string
			IdentityId        uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr         string
				Outtransactionout string
			}{retdata.Outputerr, retdata.Outtransactionout}),
		}
		retSerializeData = util.Serialize(encapdata)
	case RpcTransBroadcast:
		retdata := struct {
			Outbroadcastreps BroadcastReply
			Outputerr        string
			IdentityId       uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr        string
				Outbroadcastreps BroadcastReply
			}{retdata.Outputerr, retdata.Outbroadcastreps}),
		}
		retSerializeData = util.Serialize(encapdata)
	case RpcTransBroadQuery:
		retdata := struct {
			Outquerystr string
			Outputerr   string
			IdentityId  uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr   string
				Outquerystr string
			}{retdata.Outputerr, retdata.Outquerystr}),
		}
		retSerializeData = util.Serialize(encapdata)
	case NetGetLocalpeer:
		retdata := struct {
			Outpeerstr string
			Outputerr  string
			IdentityId uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr  string
				Outpeerstr string
			}{retdata.Outputerr, retdata.Outpeerstr}),
		}
		retSerializeData = util.Serialize(encapdata)
	case NetGetPeers:
		retdata := struct {
			Outpeersstr []string
			Outputerr   string
			IdentityId  uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr   string
				Outpeersstr []string
			}{retdata.Outputerr, retdata.Outpeersstr}),
		}
		retSerializeData = util.Serialize(encapdata)
	case LedgerBalance:
		retdata := struct {
			Outbalance Balance
			Outputerr  string
			IdentityId uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr  string
				Outbalance Balance
			}{retdata.Outputerr, retdata.Outbalance}),
		}
		retSerializeData = util.Serialize(encapdata)
	case LedgerHeight:
		retdata := struct {
			Outheight  uint32
			Outputerr  string
			IdentityId uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr string
				Outheight uint32
			}{retdata.Outputerr, retdata.Outheight}),
		}
		retSerializeData = util.Serialize(encapdata)
	case LastBlockHash:
		retdata := struct {
			Lastblockhash Hash
			Outputerr     string
			IdentityId    uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr     string
				Lastblockhash Hash
			}{retdata.Outputerr, retdata.Lastblockhash}),
		}
		retSerializeData = util.Serialize(encapdata)
	case NumberForBlockHash:
		retdata := struct {
			Outblockhash Hash
			Outputerr    string
			IdentityId   uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr    string
				Outblockhash Hash
			}{retdata.Outputerr, retdata.Outblockhash}),
		}
		retSerializeData = util.Serialize(encapdata)
	case NumberForBlock:
		retdata := struct {
			Outputblock Block
			Outputerr   string
			IdentityId  uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		logger.Debugf("The identity id is %s", uid.String())
		logger.Debugf("The block info is blockheader: %v , hash list: %v ", retdata.Outputblock.BlockHeader, retdata.Outputblock.TxHashList)
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr   string
				Outputblock Block
			}{retdata.Outputerr, retdata.Outputblock}),
		}
		retSerializeData = util.Serialize(encapdata)
	case HashForBlock:
		retdata := struct {
			Outputblock Block
			Outputerr   string
			IdentityId  uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr   string
				Outputblock Block
			}{retdata.Outputerr, retdata.Outputblock}),
		}
		retSerializeData = util.Serialize(encapdata)
	case HashForTx:
		retdata := struct {
			OutputTx   Transaction
			Outputerr  string
			IdentityId uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr string
				OutputTx  Transaction
			}{retdata.Outputerr, retdata.OutputTx}),
		}
		retSerializeData = util.Serialize(encapdata)
	case BlockHashForTxs:
		retdata := struct {
			Outputtxs  Transactions
			Outputerr  string
			IdentityId uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr string
				Outputtxs Transactions
			}{retdata.Outputerr, retdata.Outputtxs}),
		}
		retSerializeData = util.Serialize(encapdata)
	case BlockNumberForTxs:
		retdata := struct {
			Outputtxs  Transactions
			Outputerr  string
			IdentityId uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr string
				Outputtxs Transactions
			}{retdata.Outputerr, retdata.Outputtxs}),
		}
		retSerializeData = util.Serialize(encapdata)
	case MergeTxHashForTxs:
		retdata := struct {
			Outputtxs  Transactions
			Outputerr  string
			IdentityId uuid.Uuid
		}{}
		util.Deserialize(payload, &retdata)
		uid = retdata.IdentityId
		encapdata := struct {
			Cmd     uint8
			Payload []byte
		}{
			msgtype, util.Serialize(struct {
				Outputerr string
				Outputtxs Transactions
			}{retdata.Outputerr, retdata.Outputtxs}),
		}
		retSerializeData = util.Serialize(encapdata)
	default:
		uid = uuid.NewV4()
		retSerializeData = util.Serialize(struct{}{})
	}

	logger.Debugf("after recontruct data of response from chain uid: %v, serialize data: %v", uid, retSerializeData)

	return uid, retSerializeData
}
