package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/orientwalt/htdf/cmd/hsbatchsend/database"
	"github.com/orientwalt/htdf/cmd/hsbatchsend/log"
	"github.com/orientwalt/htdf/cmd/hsbatchsend/models"
	batchSendRpc "github.com/orientwalt/htdf/cmd/hsbatchsend/rpc"
	"github.com/orientwalt/htdf/utils/unit_convert"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

var (
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

type SendSysParam struct {
	ChainId   string
	Gas       string
	FeeAmount string
}

func RefreshGrantAccount(nodeRpc *batchSendRpc.NodeRpc) (err error) {

	accounts, err := nodeRpc.GetAccountList()
	if err != nil {
		log.Instance().Error("GetAccountList error|err=%s", err)
		return err
	}
	//fmt.Printf("accounts=%v\n", accounts)

	for _, account := range accounts {

		balance, seq, err := nodeRpc.GetAccountBalanceSeq(account.Address, unit_convert.BigDenom)
		if err != nil {
			log.Instance().Error("GetAccountBalanceSeq error|err=%s", err)
			return err
		}

		log.Instance().Debug("accountInfo|address=%s|balance=%s|seq=%d", account.Address, balance, seq)

		grantAccountInDB, err := models.GetGrantAccount(account.Address)
		if err != nil {
			log.Instance().Error("GetGrantAccount error|err=%s", err)
			return err
		}
		log.Instance().Debug("grantAccountInDB=%v", grantAccountInDB)

		var currGrantAccount models.GrantAccount
		currGrantAccount.ID = grantAccountInDB.ID
		currGrantAccount.Address = account.Address
		currGrantAccount.Balance, _ = decimal.NewFromString(balance)
		currGrantAccount.Seq = int64(seq)

		err = models.UpdateOrInsertGrantAccount(&currGrantAccount)
		if err != nil {
			log.Instance().Error("UpdateOrInsertGrantAccount error|err=%s", err)
			return err
		}

	}

	return nil
}

func InsertSendDetai(nodeRpc *batchSendRpc.NodeRpc) (errCode int) {

	//
	err := RefreshGrantAccount(nodeRpc)
	if err != nil {
		log.Instance().Error("RefreshGrantAccount error|err=%s", err)
		return models.Err_Process
	}

	//
	grantAccounts, err := models.GetGrantAccounts(100)
	if err != nil {
		log.Instance().Error("GetGrantAccounts error|err=%s", err)
		return models.Err_Process
	}
	log.Instance().Debug("fromAddress=%v", grantAccounts)

	if len(grantAccounts) == 0 {
		log.Instance().Error("grantAccounts not found")
		return models.Err_NotFound
	}

	db := database.GetDB().Begin()
	defer func() {
		if errCode != models.No_Error {
			db.Rollback()
		}
	}()

	rowsToGet := 1000 * len(grantAccounts)

	for {
		//
		toSends, err := models.GetToSends(rowsToGet)
		if err != nil {
			log.Instance().Error("GetToSends error|err=%s", err)
			return models.Err_Process
		}

		if len(toSends) == 0 {
			break
		}

		log.Instance().Debug("curr toSends|len=%d", len(toSends))

		for _, tosend := range toSends {
			var sendDetail models.SendDetail
			sendDetail.ToSendId = tosend.ID
			sendDetail.ToAddress = tosend.Address
			sendDetail.Amount = tosend.Amount

			choiceIndex := sendDetail.ToSendId % uint(len(grantAccounts))

			log.Instance().Debug("sendDetail.ToSendId=%d|len(grantAccounts)=%d|choiceIndex=%d",
				sendDetail.ToSendId, len(grantAccounts), choiceIndex)

			currGrantAccount := grantAccounts[choiceIndex]

			sendDetail.FromAddress = currGrantAccount.Address
			sendDetail.Seq = currGrantAccount.Seq

			err = models.InsertSendDetail(&sendDetail)
			if err != nil {
				log.Instance().Error("InsertSendDetail error|err=%s", err)
				return models.Err_Process
			}

			//grantAccount seq++
			currGrantAccount.Seq++

			models.UpdateToSendStatus(tosend.ID, 1)
		}

	}

	if err := db.Commit().Error; err != nil {
		log.Instance().Error("commit error|err=%s", err)
		return models.Err_Process
	}

	return 0
}

func SendTx(nodeRpc *batchSendRpc.NodeRpc, sendSysParam *SendSysParam) (errCode int) {

	grantAccounts, err := models.GetGrantAccounts(100)
	if err != nil {
		log.Instance().Error("GetGrantAccounts error|err=%s", err)
		return models.Err_Process
	}

	passwd := "12345678"

	for _, grantAccount := range grantAccounts {

		sendDetail, err := models.GetSendDetail(grantAccount.Address, uint64(grantAccount.Seq), 0)
		if err != nil {

			if err != gorm.ErrRecordNotFound {
				log.Instance().Info("send detail not found, sleep|err=%s", err)

				time.Sleep(10 * time.Second)
				continue
			}

			log.Instance().Error("GetGrantAccounts error|err=%s", err)
			return models.Err_Process
		}

		strMemo := "batch_" + strconv.Itoa(int(grantAccount.Seq))
		strSendReq := fmt.Sprintf(SENDREQMODEL, grantAccount.Address, strMemo, passwd, sendSysParam.ChainId,
			grantAccount.Seq, sendSysParam.Gas, sendSysParam.FeeAmount, sendDetail.Amount, sendDetail.ToAddress)

		log.Instance().Info("strSendReq=%s", strSendReq)

		sendResp, err := nodeRpc.SendTx(strSendReq)
		if err != nil {
			log.Instance().Error("SendTx error, sleep|err=%s|seq=%d", err, grantAccount.Seq)

			time.Sleep(10 * time.Second)
			continue
		}

		log.Instance().Info("SendTx ok|txHash=%s|height=%s|seq=%d", sendResp.TxHash, sendResp.Height, grantAccount.Seq)

		errCode := UpdateSendDetail(sendDetail.ID, 1)
		if errCode != 0 {
			log.Instance().Error("UpdateSendDetail error|errCode=%d", errCode)
		} else {
			log.Instance().Debug("UpdateSendDetail ok")
		}

		grantAccount.Seq++

	}

	return 0
}

func UpdateSendDetail(id uint, status int) (errCode int) {

	db := database.GetDB().Begin()
	defer func() {
		if errCode != models.No_Error {
			db.Rollback()
		}
	}()

	err := models.UpdateSendDetailStatus(id, status)
	if err != nil {
		log.Instance().Error("UpdateSendDetailStatus error|err=%s", err)
		return models.Err_Process
	}

	if err := db.Commit().Error; err != nil {
		log.Instance().Error("commit error|err=%s", err)
		return models.Err_Process
	}

	return 0
}
