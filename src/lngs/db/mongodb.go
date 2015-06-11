package lngsdb

import (
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"

	// "lngs/common"
)

type DBConn struct {
	session *mgo.Session
	db      *mgo.Database
}

func (self *DBConn) Close() {
	self.session.Close()
}

func ConnectDB() *DBConn {
	// config := lngscommon.GetConfig()

	session, err := mgo.Dial("127.0.0.1:27017") //连接数据库
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	db := session.DB("lngs") //数据库名称
	// collection := db.C("person") //如果该集合已经存在的话，则直接返回
	return &DBConn{session, db}
}
