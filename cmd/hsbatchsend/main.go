package main

import (
	"fmt"
	"github.com/orientwalt/htdf/cmd/hsbatchsend/config"
	"github.com/orientwalt/htdf/cmd/hsbatchsend/log"
	batchSendRpc "github.com/orientwalt/htdf/cmd/hsbatchsend/rpc"
	"github.com/orientwalt/htdf/cmd/hsbatchsend/service"
	"github.com/orientwalt/htdf/cmd/hsbatchsend/utils"
	"os"
	"syscall"
	"time"
)

var (
	loopStop = 0

	SENDREQMODEL = `
    { "base_req": 
      { "from": "%s", 
        "memo": "%s",
        "password": "%s", 
        "chain_id": "%s", 
        "account_number": "0", 
        "sequence": "%d", 
        "gas": "%s", 
        "fees": [ 
              { "denom": "htdf",
                 "amount": "%s" } 
         ], 
         "simulate": false
      },          
      "amount": [ 
              { "denom": "htdf", 
                "amount": "%s" } ],
      "to": "%s"
    }
`
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("usage: %s #configFilePath\n", os.Args[0])
		os.Exit(1)
	}

	// read config
	err := config.Init("config", os.Args[1])
	if err != nil {
		fmt.Printf("GetInstanceConfig error|err=%s\n", err)
		os.Exit(1)
	}

	// log init
	strLogStr := fmt.Sprintf(`{"filename":"%s"}`, config.GetString("log.logFile", ""))
	log.Instance().SetLogger("file", strLogStr)
	log.Instance().SetLevel(config.GetInt("log.logLevel", 6)) //#logLevel  : 1, Alert; 2, Crit; 3, Error; 4, Warn; 5, Notice; 6, Info; 7, Debug
	log.Instance().EnableFuncCallDepth(true)

	log.Instance().Info("app start")

	// signal handler
	go utils.WaitSignal(signalHandler)

	rpcUrl := config.GetString("fullnode.rpcUrl", "")
	restDebug := config.GetInstanceConfig().GetBool("fullnode.restDebug")

	var sendSysParam service.SendSysParam
	sendSysParam.ChainId = config.GetString("fullnode.chainId", "")
	sendSysParam.Gas = config.GetString("fullnode.gas", "0")
	sendSysParam.FeeAmount = config.GetString("fullnode.fee_amount", "0")

	log.Instance().Info("chainId=%s|rpcUrl=%s|restDebug=%v|gas=%s|fee_amount=%s",
		sendSysParam.ChainId, rpcUrl, restDebug, sendSysParam.Gas, sendSysParam.FeeAmount)

	//init rpc object
	nodeRpc := batchSendRpc.GetInstanceNodeRpc()
	nodeRpc.SetUrl(rpcUrl)
	nodeRpc.SetDebug(restDebug)

	//wait for block package, and get the right seq
	time.Sleep(8 * time.Second)

	//insertSendDetail
	errCode := service.InsertSendDetai(nodeRpc)
	if errCode != 0 {
		log.Instance().Error("InsertSendDetai error|err=%s", err)
		os.Exit(1)
	}

	//sendDetail, err := models.GetSendDetail("htdf1gh8yeqhxx7n29fx3fuaksqpdsau7gxc0rteptt", 0, 0)
	//if err != nil {
	//	log.Instance().Error("GetSendDetail error|err=%s", err)
	//	os.Exit(1)
	//}

	//log.Instance().Debug("sendDetail=%v", sendDetail)

	errCode = service.SendTx(nodeRpc, &sendSysParam)
	if errCode != 0 {
		log.Instance().Error("SendTx error|errCode=%d", errCode)
		os.Exit(1)
	}

	//==============================================================================
	//if false {
	//
	//	for iCurrSeq := iStartSeq; iCurrSeq <= iEndSeq; iCurrSeq++ {
	//
	//		log.Instance().Debug("iCurrSeq=%d\n", iCurrSeq)
	//
	//		strMemo := "batch_" + strconv.Itoa(int(iCurrSeq))
	//
	//		strSendReq := fmt.Sprintf(SENDREQMODEL, fromAddr, strMemo, passwd, chainId,
	//			iCurrSeq, gas, fee_amount, amount, toAddr)
	//
	//		log.Instance().Info("strSendReq=%s\n", strSendReq)
	//
	//		sendResp, err := nodeRpc.SendTx(strSendReq)
	//		if err != nil {
	//			log.Instance().Error("SendTx error|err=%s|seq=%d\n", err, iCurrSeq)
	//			os.Exit(1)
	//		}
	//
	//		log.Instance().Info("SendTx ok|txHash=%s|height=%s|seq=%d\n", sendResp.TxHash, sendResp.Height, iCurrSeq)
	//	}
	//
	//	for {
	//		if loopStop != 0 {
	//			break
	//		}
	//
	//		time.Sleep(1 * time.Second)
	//	}
	//
	//}

	log.Instance().Info("app exit")

}

func stopLoop() {
	loopStop = 1
}

func signalHandler(s os.Signal) {

	switch s {
	case syscall.SIGINT:
		//log.Instance().Info("ignore SIGINT..")
	case syscall.SIGHUP:
		log.Instance().Info("ignore SIGHUP..")
	case syscall.SIGKILL:
		log.Instance().Info("ignore SIGKILL..")
	case syscall.SIGTERM:
		log.Instance().Info("process SIGTERM..")

		stopLoop()
	default:
		log.Instance().Info("unknown signal..")
	}

}
