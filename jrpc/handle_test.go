package jrpc

import (
	"fmt"
	"testing"

	"github.com/bocheninc/msg-net/util"
)

func TestJsonMarshal(t *testing.T) {
	res := createRes(22, ldpQueryConfig, []string{})
	data, err := res.byte()

	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(data))

	data = util.Serialize(msgnetMessage{Cmd: chainRPCMsg, Payload: data})

	msg := &msgnetMessage{}
	util.Deserialize(data, msg)
	fmt.Println(string(msg.Payload))
}
