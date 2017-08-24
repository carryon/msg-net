package http_rpc_server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/bocheninc/msg-net/logger"
	"github.com/bocheninc/msg-net/peer"
	"github.com/bocheninc/msg-net/util"
	"github.com/twinj/uuid"
)

var (
	virtualP0  *peer.Peer
	channelmap map[string]chan []byte
	refroute   pRouterHandler
	mux        sync.Mutex
)

func parseRpcPost(w http.ResponseWriter, request *http.Request) {
	// Read body
	b, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msg message
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if len(msg.ChainId) == 0 {
		http.Error(w, "should specifiy the chain id", 500)
		return
	}

	chainstrId := msg.ChainId

	isPeerExist := refroute.IsPeerExist(chainstrId)
	if !isPeerExist {
		http.Error(w, "chain id is not valid", 500)
		return
	}

	u := uuid.NewV4()

	output, _ := json.Marshal(struct {
		Data       []byte
		IdentityId uuid.Uuid `json:"uuid"`
	}{b, u})

	data := sendMessage{}
	data.Cmd = ChainRpcMsg
	data.Payload = output

	sendbytes := util.Serialize(data)

	virtualP0.Send(chainstrId, sendbytes, nil)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	recvchan := make(chan []byte)

	keystr := u.String()

	logger.Debugf("The send request uid is %s ", keystr)

	mux.Lock()
	channelmap[keystr] = recvchan
	mux.Unlock()

	var httpOutput SerialDataType

	select {
	case buf := <-recvchan:
		transformResponseData(buf, &httpOutput)
	case <-time.After(time.Second * 5):
		httpOutput = []byte(`{"result","timeout"}`)
		mux.Lock()
		delete(channelmap, keystr)
		mux.Unlock()
	}

	close(recvchan)
	w.Header().Set("content-type", "application/json")
	w.Write(httpOutput)
}

func chainMessageHandle(srcID, dstID string, payload []byte, signature []byte) error {
	msg := &sendMessage{}
	util.Deserialize(payload, msg)

	logger.Debugf("before recontruct data of response from chain msg type: %v, payload: %v", msg.Cmd, msg.Payload)

	identityId, serialise_data := recontructData(msg.Cmd, msg.Payload)
	logger.Debugf("The identity id from response is %s ", identityId.String())

	mux.Lock()
	if channel, ok := channelmap[identityId.String()]; ok {
		delete(channelmap, identityId.String())
		channel <- serialise_data
	} else {
		logger.Debugf("Unknown message from the msg-net for rpc msg")
	}
	mux.Unlock()

	logger.Debugf("rpc response rom block chain  dst: %s, src: %s\n", dstID, srcID)

	return nil
}

func RunRpcServer(port, address string, proute pRouterHandler) {
	refroute = proute
	virtualP0 = peer.NewPeer("01:595959", []string{address}, chainMessageHandle)
	virtualP0.Start()
	channelmap = make(map[string]chan []byte)
	http.HandleFunc("/", parseRpcPost)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.Errorf("Run msgnet rpc server fail %s", err.Error())
	} else {
		logger.Debug("Run msgnet rpc server")
	}
}
