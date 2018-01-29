package service

import (
	"fmt"
	"git/inspursoft/board/src/common/utils"
	"os"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func connectToDB() {
	hostIP := os.Getenv("HOST_IP")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", fmt.Sprintf("root:root123@tcp(%s:3306)/board?charset=utf8", hostIP))
	if err != nil {
		logs.Error("Failed to connect to DB.")
	}

}

func TestMain(m *testing.M) {
	connectToDB()
	utils.Initialize()
	utils.SetConfig("SSH_KEY_PATH", "/Users/wangkun/keys")
	os.Exit(m.Run())
}
