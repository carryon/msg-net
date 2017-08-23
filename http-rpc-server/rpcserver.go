package http_rpc_server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/bocheninc/msg-net/logger"
	"github.com/bocheninc/msg-net/peer"
	"github.com/bocheninc/msg-net/utils"
	"github.com/twinj/uuid"
)

var (
	virtualP0   *peer.Peer
	channelmap  map[string]chan []byte
	wg          sync.WaitGroup
	refroute    pRouterHandler
	defaultchan chan []byte
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

	sendbytes := utils.Serialize(data)

	virtualP0.Send(chainstrId, sendbytes, nil)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	recvchan := make(chan []byte)

	keystr := u.String()

	logger.Debugf("The send request uid is %s ", keystr)

	channelmap[keystr] = recvchan

	defer wg.Done()

	for {
		select {
		case recieve_msg := <-recvchan:
			var httpoutput SerialDataType
			transformResponseData(recieve_msg, &httpoutput)
			w.Header().Set("content-type", "application/json")
			w.Write(httpoutput)
			close(recvchan)
			delete(channelmap, keystr)
			return
		case <-time.After(time.Second * 5):
			var httpoutput SerialDataType
			w.Header().Set("content-type", "application/json")
			httpoutput, _ = json.Marshal(struct {
				Err string `json:"error"`
			}{
				"Timeout for request",
			})
			w.Write(httpoutput)
			close(recvchan)
			delete(channelmap, keystr)
			return
		}
	}
}

func chainMessageHandle(srcID, dstID string, payload []byte, signature []byte) error {
	msg := &sendMessage{}
	utils.Deserialize(payload, msg)
	wg.Add(1)

	logger.Debugf("before recontruct data of response from chain msg type: %v, payload: %v", msg.Cmd, msg.Payload)

	identityid, serialise_data := recontructData(msg.Cmd, msg.Payload)
	logger.Debugf("The identity id from response is %s ", identityid.String())
	go func() {
		if channel, ok := channelmap[identityid.String()]; ok {
			channel <- serialise_data
		} else {
			defaultchan <- serialise_data
		}
	}()

	logger.Debugf("rpc response rom block chain  dst: %s, src: %s\n", dstID, srcID)

	wg.Wait()

	return nil
}

func RunRpcServer(port, address string, proute pRouterHandler) {
	http.HandleFunc("/", parseRpcPost)
	refroute = proute
	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			logger.Errorf("Run msgnet rpc server fail %s", err.Error())
		} else {
			logger.Debug("Run msgnet rpc server")
		}
	}()

	virtualP0 = peer.NewPeer("01:595959", []string{address}, chainMessageHandle)
	virtualP0.Start()
	channelmap = make(map[string]chan []byte)
	defaultchan = make(chan []byte)
	for {
		select {
		case <-defaultchan:
			wg.Done()
			logger.Debugf("Unknown message from the msg-net for rpc msg")
		}
	}

}
