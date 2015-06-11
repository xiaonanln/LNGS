package lngsdb

import (
	. "lngs/common"
	"log"
)

var (
	dbmanagers = make([]*DbManager, 0) // 所有的dbmanager列表
)

type DbManager struct {
	conn         *DBConn
	requestqueue CommandQueue
}

func NewDbManager(conn *DBConn) *DbManager {
	dbm := &DbManager{conn, make(CommandQueue)}
	dbmanagers = append(dbmanagers, dbm)
	log.Printf("DbManager created for DB connection %v, number of dbmanagers is %d\n", conn, len(dbmanagers))
	return dbm
}

func (self *DbManager) Loop() {
	for {
		request, ok := <-self.requestqueue

		if !ok {
			break
		}

		request = request
	}
}

func (self *DbManager) PostCommand(request *Command) {
	self.requestqueue <- request
}

func PostDbCommand(cmd *Command) {

}
