package database

import (
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/orientwalt/htdf/cmd/hsbatchsend/config"
	"github.com/orientwalt/htdf/cmd/hsbatchsend/log"
	"github.com/orientwalt/htdf/cmd/hsbatchsend/utils"
)

var (
	globalDB   *gorm.DB
	onceInitDB sync.Once
)

func InitDB() error {
	//"loop:loop1234@tcp(localhost:3306)/metis1?parseTime=true"

	if db, err := openDatabase(); err == nil && db != nil {
		SetDB(db)
		return nil
	} else {
		return err
	}
}

func openDatabase() (*gorm.DB, error) {
	var db *gorm.DB

	err := utils.Do(func(attempt int) (bool, error) {
		var err error

		v := config.GetInstanceConfig()
		provider := v.GetString("database.provider")
		connectString := v.GetString("database.connectStr")
		maxConns := v.GetInt("database.maxOpenConn")
		maxIdleConns := v.GetInt("database.maxIdleConn")

		db, err = gorm.Open(provider, connectString)
		if err != nil {
			//time.Sleep(1 * time.minute) // wait a minute
		} else {
			// Disable plular
			db.SingularTable(true)
			db.LogMode(v.GetBool("database.log"))
			//db.SetLogger(log.Instance().log.CurrentLogger())
			if maxConns > 0 && maxIdleConns > 0 {
				db.DB().SetMaxOpenConns(maxConns)
				db.DB().SetMaxIdleConns(maxIdleConns)
			}
			log.Instance().Info("Open DB %s success", provider)
		}
		return attempt < 1, err
	})
	if err != nil {
		panic(err.Error())
	}

	return db, err
}

func LogMode(enable bool) {
	if globalDB != nil {
		globalDB.LogMode(enable)
	}
}

func GetDB() *gorm.DB {
	if globalDB == nil {
		onceInitDB.Do(func() {
			InitDB()
		})
	}
	return globalDB
}

func SetDB(db *gorm.DB) {
	globalDB = db
}
