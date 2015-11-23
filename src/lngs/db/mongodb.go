package lngsdb

import (
	"log"
	"sync"
	"time"

	"gopkg.in/mgo.v2"

	// "lngs/common"
)

func init() {
	go serveMongoDB()
}

type DBConn struct {
	session *mgo.Session
	db      *mgo.Database
}

func (self *DBConn) Close() {
	self.session.Close()
}

var (
	conn     *DBConn = nil
	connLock sync.RWMutex
)

func connectDB() {
	// config := lngscommon.GetConfig()
	connLock.Lock()
	defer connLock.Unlock()

	log.Println("Connecting MongoDB ...")
	session, err := mgo.Dial("127.0.0.1:27017") //连接数据库
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	db := session.DB("lngs") //数据库名称
	// collection := db.C("person") //如果该集合已经存在的话，则直接返回

	conn = &DBConn{session, db}
	log.Println("MongoDB connected", db)
}

func handleDBError(err error) {
	log.Println("DB ERROR: ", err)
}

func serveMongoDB() {
	for {
		if conn == nil {
			connectDB()
		}

		time.Sleep(1e9)
	}
}

func waitForDBConn() {
restart:
	for conn == nil {
		time.Sleep(1e9)
	}

	connLock.RLock()
	if conn == nil {
		connLock.RUnlock()
		goto restart
	}
	// now the conn is locked, and conn != nil

}

func FindDoc(collectionName string, selector interface{}) (Doc, error) {
	waitForDBConn()
	defer connLock.RUnlock()

	doc := make(Doc)
	err := conn.db.C(collectionName).Find(selector).One(doc)
	if err != nil {
		return nil, err
	} else {
		return doc, nil
	}
}

func UpdateDoc(collectionName string, selector interface{}, update Doc) error {
	waitForDBConn()
	defer connLock.RUnlock()

	err := conn.db.C(collectionName).Update(selector, update)
	return err
}

func InsertDoc(collectionName string, doc Doc) error {
	waitForDBConn()
	defer connLock.RUnlock()

	err := conn.db.C(collectionName).Insert(doc)
	return err
}
