package lngsdb

import (
	. "lngs/common"
	"log"
	"math/rand"
)

var (
	dbmanagers = make([]*DbManager, 0) // 所有的dbmanager列表
)

type DbManager struct {
	conn         *DBConn
	commandQueue commandQueue
}

func NewDbManager(conn *DBConn) *DbManager {
	dbm := &DbManager{conn, make(commandQueue)}
	dbmanagers = append(dbmanagers, dbm)
	log.Printf("DbManager created for DB connection %v, number of dbmanagers is %d\n", conn, len(dbmanagers))
	return dbm
}

func (self *DbManager) Loop() {
	db := self.conn.db
	for {
		cmd, ok := <-self.commandQueue

		if !ok {
			break
		}

		switch cmd.Command {
		case "insertdb":
			{
				insertdbArgs := cmd.data.([]interface{})
				collectionName := insertdbArgs[0].(string)
				doc := insertdbArgs[1]
				err := db.Collection(collectionName).Insert(doc)
				// send back result
				if err != nil {
					log.Println(err)
				}
				EntityManager
			}
		}
	}
}

func (self *DbManager) PostCommand(cmd *Command) {
	self.commandQueue <- cmd
}

func PostDbCommand(cmd *Command) {
	r := rand.Intn(len(dbmanagers))
	dbmanager := dbmanagers[r] // choose random dbmanager from all
	dbmanager.PostCommand(cmd)
}
