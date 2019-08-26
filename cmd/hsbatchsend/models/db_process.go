package models

import (
	"github.com/jinzhu/gorm"
	"github.com/orientwalt/htdf/cmd/hsbatchsend/database"
	"github.com/orientwalt/htdf/cmd/hsbatchsend/log"
	"github.com/shopspring/decimal"
)

//
type GrantAccount struct {
	gorm.Model
	Address string
	Seq     int64
	Balance decimal.Decimal `sql:"type:decimal(36,18);"`
}

type ToSend struct {
	gorm.Model
	Address string
	Amount  decimal.Decimal `sql:"type:decimal(36,18);"`
	Status  int
}

type SendDetail struct {
	gorm.Model
	ToSendId    uint
	ToAddress   string
	Amount      decimal.Decimal `sql:"type:decimal(36,18);"`
	FromAddress string
	Seq         int64
	Status      int
	TxHash      string
}

func UpdateOrInsertGrantAccount(grantAccount *GrantAccount) error {
	db := database.GetDB()
	err := db.Save(grantAccount).Error
	if err != nil {
		log.Instance().Error("UpdateOrInsertGrantAccount error|err=%s", err)
		return err
	}

	return nil
}

func GetGrantAccount(address string) (*GrantAccount, error) {
	db := database.GetDB()
	record := new(GrantAccount)
	err := db.Where("address = ?", address).First(&record).Error
	return record, err
}

func GetGrantAccounts(limit int) ([]*GrantAccount, error) {
	db := database.GetDB()
	var records []*GrantAccount
	err := db.Limit(limit).Find(&records).Error
	return records, err
}

func GetToSends(rowsToGet int) ([]ToSend, error) {
	db := database.GetDB()
	var records []ToSend
	err := db.Where("status is null or status = ?", 0).Limit(rowsToGet).Find(&records).Error
	return records, err
}

func UpdateToSendStatus(id uint, status int) (err error) {
	db := database.GetDB()
	if err = db.Model(&ToSend{}).Where("id =?", id).Update("status", status).Error; err != nil {
		return err
	}

	return err
}

func InsertSendDetail(sendDetail *SendDetail) error {
	db := database.GetDB()
	err := db.Create(sendDetail).Error
	if err != nil {
		log.Instance().Error("Create sendDetail error|err=%s", err)
		return err
	}

	return nil
}

func GetSendDetail(fromAddress string, seq uint64, status int) (*SendDetail, error) {
	db := database.GetDB()
	sendDetail := new(SendDetail)
	err := db.Where("from_address = ? and seq = ?  and status = ? ", fromAddress, seq, status).First(&sendDetail).Error
	return sendDetail, err
}

func UpdateSendDetailStatus(id uint, status int) (err error) {
	db := database.GetDB()
	if err = db.Model(&SendDetail{}).Where("id =?", id).Update("status", status).Error; err != nil {
		log.Instance().Error("update SendDetail error|err=%s", err)
		return err
	}

	return err
}

func Commit() {
	db := database.GetDB()
	if err := db.Commit().Error; err != nil {
		log.Instance().Error("commit error|err=%s", err)
		return
	}
}
