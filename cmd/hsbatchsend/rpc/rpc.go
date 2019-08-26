package rpc

import (
	"encoding/json"
	"errors"
	"github.com/orientwalt/htdf/accounts"
	"github.com/orientwalt/htdf/cmd/hsbatchsend/log"
	"github.com/orientwalt/htdf/cmd/hsbatchsend/models"
	"strconv"
	"strings"
	"sync"
)

import (
	"fmt"
	"gopkg.in/resty.v1"
)

type NodeRpc struct {
	hostUrl string
	debug   bool

	restyClientNode_rpc *resty.Client
}

var (
	globalNodeRpc   *NodeRpc
	onceInitNodeRpc sync.Once
)

func initNodeRpc() {
	globalNodeRpc = NewNodeRpc()
}

func GetInstanceNodeRpc() *NodeRpc {
	onceInitNodeRpc.Do(func() {
		initNodeRpc()
	})
	return globalNodeRpc
}

func SetNodeRpc(nodeRpc *NodeRpc) {
	globalNodeRpc = nodeRpc
}

func NewNodeRpc() *NodeRpc {
	nodeRpc := NodeRpc{debug: true}
	nodeRpc.init()
	return &nodeRpc
}

func (nodeRpc *NodeRpc) init() {
	nodeRpc.restyClientNode_rpc = resty.New()
	nodeRpc.restyClientNode_rpc.SetHostURL(nodeRpc.hostUrl)

	// Headers for all request
	nodeRpc.restyClientNode_rpc.SetHeader("Accept", "application/json")
	//nodeRpc.restyClientNode_rpc.SetBasicAuth("x", nodeRpc.password)

	nodeRpc.restyClientNode_rpc.Debug = nodeRpc.debug

}

func init() {

}

func (nodeRpc *NodeRpc) SetUrl(url string) {
	nodeRpc.hostUrl = url
	nodeRpc.restyClientNode_rpc.SetHostURL(url)
}

func (nodeRpc *NodeRpc) SetDebug(debugOn bool) {
	nodeRpc.debug = debugOn
	nodeRpc.restyClientNode_rpc.Debug = debugOn
}

//@brief RPC call enencap
//	respon json text has not uniform format
func (nodeRpc *NodeRpc) RpcCallSimple(url string, body []byte) (string, error) {
	res, err := nodeRpc.restyClientNode_rpc.R().Get(url)
	if err != nil {
		log.Instance().Error("http GET error|err=%s|url=%s", err, url)
		return "", err
	}

	//fmt.Println(res.String())

	return res.String(), nil
}

//@brief RPC call enencap
//	respon json text has no uniform format
func (nodeRpc *NodeRpc) RpcCall(url string, body []byte) (resp string, err error) {
	res, err := nodeRpc.restyClientNode_rpc.R().SetBody(body).Post(url)
	if err != nil {
		log.Instance().Error("http POST error|err=%s|url=%s", err, url)
		return "", err
	}

	//log.Instance().Debug(res.String())

	return res.String(), nil
}

func (nodeRpc *NodeRpc) GetNodeInfo() error {
	url := fmt.Sprintf("/node_info")
	body := []byte(fmt.Sprintf(` `))

	respString, err := nodeRpc.RpcCallSimple(url, body)
	if err != nil {
		log.Instance().Error("RpcCallSimple error|err=%s|url=%s", err, url)
		return err
	}

	log.Instance().Debug("respString=%s\n", respString)

	return nil
}

func (nodeRpc *NodeRpc) SendTx(strBody string) (*SendResp, error) {
	url := fmt.Sprintf("/hs/send")
	body := []byte(fmt.Sprintf(`%s`, strBody))

	resp, err := nodeRpc.RpcCall(url, body)
	if err != nil {
		log.Instance().Error("RpcCall error|err=%s|url=%s", err, url)
		return nil, err
	}

	log.Instance().Debug("resp=%v\n", resp)

	if strings.Index(resp, `"message"`) >= 0 {
		log.Instance().Error("err resp from peer|errStr=%s", resp)
		return nil, errors.New("err resp from peer")
	}

	if len(resp) == 0 {
		log.Instance().Error("null resp from peer", err, url)
		return nil, errors.New("null resp from peer")
	}

	data := []byte(resp)
	var sendResp SendResp

	err = json.Unmarshal(data, &sendResp)
	if err != nil {
		log.Instance().Error("Unmarshal error|err=%s|url=%s", err, url)
		return nil, err
	}

	if sendResp.TxHash == "" {
		log.Instance().Error("null txHash")
		return nil, errors.New("null txhash")
	}

	return &sendResp, nil
}

func (nodeRpc *NodeRpc) GetAccountInfo(address string) (*AccountInfo, int) {

	url := fmt.Sprintf("/auth/accounts/%s", address)
	body := []byte(fmt.Sprintf(` `))

	respString, err := nodeRpc.RpcCallSimple(url, body)
	if err != nil {
		log.Instance().Error("RpcCallSimple error|err=%s|url=%s", err, url)
		return nil, models.Err_Process
	}

	//log.Instance().Debug("respString=%s\n", respString)

	data := []byte(respString)
	var accountInfo AccountInfo

	if len(respString) == 0 {
		log.Instance().Debug("account not found")
		return nil, models.Err_NotFound
	}

	err = json.Unmarshal(data, &accountInfo)
	if err != nil {
		log.Instance().Error("Unmarshal error|err=%s|url=%s", err, url)
		return nil, models.Err_Process
	}

	return &accountInfo, 0
}

func (nodeRpc *NodeRpc) GetAccountBalanceSeq(address string, denom string) (balance string, seq uint64, err error) {
	balance = "0.0"
	accountInfo, errCode := nodeRpc.GetAccountInfo(address)
	if errCode != 0 {
		if errCode == models.Err_NotFound {
			return balance, 0, nil
		}

		log.Instance().Error("GetAccountInfo error|err=%s|address=%s", err, address)
		return "", 0, errors.New("GetAccountInfo")
	}

	//fmt.Printf("accountInfo=%v\n", accountInfo)

	seq, err = strconv.ParseUint(accountInfo.Value.Sequence, 10, 64)
	if err != nil {
		log.Instance().Error("ParseUint error|err=%s|address=%s", err, address)
		return "", 0, err
	}

	for _, coin := range accountInfo.Value.Coins {
		if coin.Denom == denom {
			balance = coin.Amount
		}
	}

	return balance, seq, nil
}

func (nodeRpc *NodeRpc) GetAccountList() (accounts []accounts.Account, err error) {

	url := fmt.Sprintf("/accounts/list?jsonformat=1")
	body := []byte(fmt.Sprintf(` `))

	respString, err := nodeRpc.RpcCallSimple(url, body)
	if err != nil {
		log.Instance().Error("RpcCallSimple error|err=%s|url=%s", err, url)
		return accounts, err
	}

	//log.Instance().Debug("respString=%s\n", respString)

	if len(respString) == 0 {
		log.Instance().Debug("account not found")
		return accounts, nil
	}

	data := []byte(respString)

	err = json.Unmarshal(data, &accounts)
	if err != nil {
		log.Instance().Error("Unmarshal error|err=%s|url=%s", err, url)
		return accounts, err
	}

	return accounts, nil
}
