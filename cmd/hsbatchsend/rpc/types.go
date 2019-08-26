package rpc

import "github.com/orientwalt/htdf/utils/unit_convert"

type SendResp struct {
	Height    string `json:"height"`
	TxHash    string `json:"txhash"`
	GasWanted string `json:"gas_wanted"`
	GasUsed   string `json:"gas_used"`
}

type AccountValue struct {
	Address  string                 `json:"address"`
	Coins    []unit_convert.BigCoin `json:"coins"`
	Sequence string                 `json:"sequence"`
}

type AccountInfo struct {
	Typed string       `json:"type"`
	Value AccountValue `json:value`
}
