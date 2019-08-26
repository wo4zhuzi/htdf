package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/orientwalt/htdf/utils/unit_convert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestRpc(t *testing.T) {

	nodeRpc := GetInstanceNodeRpc()
	nodeRpc.SetUrl("http://192.168.10.120:1317")
	nodeRpc.SetDebug(false)

	err := nodeRpc.GetNodeInfo()
	assert.Equal(t, err, nil)

	address := "htdf1gh8yeqhxx7n29fx3fuaksqpdsau7gxc0rteptt"
	accountInfo, errCode := nodeRpc.GetAccountInfo(address)
	assert.Equal(t, errCode, 0)
	fmt.Printf("accountInfo=%v\n", accountInfo)

	balance, seq, err := nodeRpc.GetAccountBalanceSeq(address, unit_convert.BigDenom)
	assert.Equal(t, err, nil)
	fmt.Printf("address=%s|balance=%s|seq=%d\n", address, balance, seq)

	address = "htdf1c3pdnq4sk4u5k7x4v8g004w3y3rzs66hx9k5k4"
	balance, seq, err = nodeRpc.GetAccountBalanceSeq(address, unit_convert.BigDenom)
	assert.Equal(t, err, nil)
	fmt.Printf("address=%s|balance=%s|seq=%d\n", address, balance, seq)

	accounts, err := nodeRpc.GetAccountList()
	assert.Equal(t, err, nil)
	fmt.Printf("accounts=%v\n", accounts)

}

func TestRpc1(t *testing.T) {
	nodeRpc := GetInstanceNodeRpc()
	nodeRpc.SetUrl("http://192.168.10.120:1317")
	nodeRpc.SetDebug(false)

	address := "htdf1gh8yeqhxx7n29fx3fuaksqpdsau7gxc0rteptt"

	balance, seq, err := nodeRpc.GetAccountBalanceSeq(address, unit_convert.BigDenom)
	assert.Equal(t, err, nil)
	fmt.Printf("address=%s|balance=%s|seq=%d\n", address, balance, seq)
}

type Server struct {
	ServerName string
	ServerIp   string
}

type ServerSlice struct {
	Server    []Server
	ServersID string
}

func TestPost(t *testing.T) {
	//post 第三个参数是io.reader interface
	//strings.NewReader  byte.NewReader bytes.NewBuffer  实现了read 方法
	s := ServerSlice{ServersID: "tearm", Server: []Server{{"beijing", "127.0.0.1"}, {"shanghai", "127.0.0.1"}}}
	b, _ := json.Marshal(s)
	fmt.Println(string(b))
	resp, _ := http.Post("http://baidu.com", "application/x-www-form-urlencoded", strings.NewReader("heel="+string(b)))
	//
	defer resp.Body.Close()
	//io.Reader

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
