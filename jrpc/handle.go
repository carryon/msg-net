package jrpc

type MsgNet struct {
}

func (m *MsgNet) getNodeLogByID(key string, relay *string) error {

	return nil
}

func (m *MsgNet) getNodeConfigByID(key string, relay *string) error {

	sendMessage(1, "00",,nil)

	return nil
}
func (m *MsgNet) getAllServers(key string, relay *Servers) error          { return nil }
func (m *MsgNet) getServersByIP(key string, relay *Servers) error         { return nil }
func (m *MsgNet) getServerStatusByIP(key string, relay *Servers) error    { return nil }
func (m *MsgNet) getAllNodeNewBlockInfo(key string, relay *Servers) error { return nil }
func (m *MsgNet) getBlockByRange(key string, relay *Servers) error        { return nil }
func (m *MsgNet) getBlockByHash(key string, relay *Servers) error         { return nil }
func (m *MsgNet) getBlockByHeight(key string, relay *Servers) error       { return nil }
func (m *MsgNet) getTxByHash(key string, relay *Servers) error            { return nil }
func (m *MsgNet) getAccountInfoByID(key string, relay *Servers) error     { return nil }
func (m *MsgNet) getAccountInfoByAddr(key string, relay *Servers) error   { return nil }
func (m *MsgNet) getHistoryTransaction(key string, relay *Servers) error  { return nil }
