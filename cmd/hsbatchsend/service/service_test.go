package service

import (
	"github.com/magiconair/properties/assert"
	"github.com/orientwalt/htdf/cmd/hsbatchsend/rpc"
	"testing"
)

func TestUpdateGrantAccount(t *testing.T) {

	nodeRpc := rpc.GetInstanceNodeRpc()
	nodeRpc.SetUrl("http://192.168.10.120:1317")
	nodeRpc.SetDebug(false)

	err := RefreshGrantAccount(nodeRpc)
	assert.Equal(t, err, nil)

}
