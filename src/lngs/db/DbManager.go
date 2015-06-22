package lngsdb

import (
	. "lngs/cmdque"
	. "lngs/common"
	"log"
	"math/rand"
)

var (
	dbmanagers          = make([]*DbManager, 0) // 所有的dbmanager列表
	dbDispatcherStarted = false
	dbCommandQueue      = GetCommandQueue("db")
)

func init() {
	go dispatchDbCommands()
}

func debug(msg string, args ...interface{}) {
	Debug("dbmanager", msg, args...)
}

func dispatchDbCommands() {
	for {
		cmd := <-dbCommandQueue

		r := rand.Intn(len(dbmanagers))
		debug("dispatching db command %v to dbmanager %d\n", cmd, r)

		dbmanager := dbmanagers[r] // choose random dbmanager from all
		dbmanager.PostCommand(cmd)
	}
}

type DbManager struct {
	conn         *DBConn
	commandQueue CommandQueue
}

func NewDbManager(conn *DBConn) *DbManager {
	dbm := &DbManager{conn, make(CommandQueue)}
	dbmanagers = append(dbmanagers, dbm)
	debug("DbManager created for DB connection %v, number of dbmanagers is %d", conn, len(dbmanagers))
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
		case "insert":
			{
				dbArgs := cmd.Data.([]interface{})
				collectionName := dbArgs[0].(string)
				doc := dbArgs[1]
				err := db.C(collectionName).Insert(doc)
				// send back result
				if err != nil {
					log.Println(err)
					PostCommandQueue(cmd.EntityId, &Command{"db", "insert_cb", err})
				} else {
					debug("insert %v", doc)
					PostCommandQueue(cmd.EntityId, &Command{"db", "insert_cb", nil})
				}
			}
		case "find":
			{
				dbArgs := cmd.Data.([]interface{})
				collectionName := dbArgs[0].(string)
				query := dbArgs[1]
				cursor := db.C(collectionName).Find(query)
				var doc Doc
				err := cursor.One(&doc)
				if err != nil {
					debug("find %v error: %s", query, err)
					PostCommandQueue(cmd.EntityId, &Command{"db", "find_cb", err})
				} else {
					debug("find %v error: %v", query, err)
					PostCommandQueue(cmd.EntityId, &Command{"db", "find_cb", doc})
				}
			}
		case "update":
			{
				dbArgs := cmd.Data.([]interface{})
				collectionName := dbArgs[0].(string)
				query := dbArgs[1]
				doc := dbArgs[2]
				err := db.C(collectionName).Update(query, doc)
				// send back result
				if err != nil {
					log.Println(err)
					PostCommandQueue(cmd.EntityId, &Command{"db", "update_cb", err})
				} else {
					debug("update %v", doc)
					PostCommandQueue(cmd.EntityId, &Command{"db", "update_cb", nil})
				}
			}
		}
	}
}

func (self *DbManager) PostCommand(cmd *Command) {
	self.commandQueue <- cmd
}

func PostDbCommand(cmd *Command) {
	dbCommandQueue <- cmd
}
