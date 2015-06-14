package lngscommon

import (
	"fmt"
	"log"
)

func Debug(category string, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	log.Printf("= %s = %s\n", category, msg)
}
