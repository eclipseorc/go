package lib

import (
	"encoding/json"
	"fmt"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

var dbHand []*xorm.Engine

type DBList struct {
	DBUser string `json:"db_user"`
	DBHome string `json:"db_home"`
	DBPort uint32 `json:"db_port"`
	DBName string `json:"db_name"`
	DBPass string `json:"db_pass"`
}

func UseHand(index int) *xorm.Engine {
	if len(dbHand) >= index {
		return dbHand[index]
	}
	return nil
}

func init() {
	var dbList []DBList
	data, err := ReadConfig("./config/db.json")
	if err != nil {
		panic("读取配置文件错误：" + err.Error())
	}
	json.Unmarshal(data, &dbList)
	fmt.Println("mysql:链接参数：", dbList)

	for _, db := range dbList {
		strconn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db.DBUser,
			db.DBPass,
			db.DBHome,
			db.DBPort,
			db.DBName)
		db_Hand, err := xorm.NewEngine("mysql", strconn)
		if err != nil {
			fmt.Println("Web database link failed", err)
			panic(err.Error())
		}
		if err := db_Hand.Ping(); err != nil {
			fmt.Println("Test link failed", err)
			panic(err.Error())
		}
		db_Hand.SetTableMapper(core.SameMapper{})
		db_Hand.SetColumnMapper(core.SameMapper{})
		db_Hand.ShowSQL(true)
		db_Hand.SetMaxIdleConns(5)
		db_Hand.SetMaxOpenConns(5)
		dbHand = append(dbHand, db_Hand)
	}

}
