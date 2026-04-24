package logger

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "", log.LstdFlags)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
	log.SetOutput(os.Stdout)
}
