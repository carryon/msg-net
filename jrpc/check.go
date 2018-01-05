package jrpc

import (
	"fmt"
	"net"
)

func checkIP(params []interface{}) (string, error) {
	if len(params) != 1 {
		return "", fmt.Errorf("params len must 1, not %v", len(params))
	}
	ipString, ok := params[0].(string)
	if !ok {
		return "", fmt.Errorf("param :%v must string", params[0])
	}
	ip := net.ParseIP(ipString)
	if ip != nil {
		return ip.String(), nil
	}
	return "", fmt.Errorf("ip :%s form error", ipString)
}
