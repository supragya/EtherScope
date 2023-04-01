package util

import (
	"fmt"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

func ENOK(err error) {
	ENOKS(2, err)
}

func ENOKS(skip int, err error) {
	if err != nil {
		_, file, no, ok := runtime.Caller(skip)
		if ok {
			fileSplit := strings.Split(file, "/")
			log.WithFields(log.Fields{
				"file": fileSplit[len(fileSplit)-1],
				"line": no,
			}).Fatalln(err)
		}
		log.Fatalln(err)
	}
}

func ENOKF(err error, info interface{}) {
	if err != nil {
		ENOK(fmt.Errorf("%s: %v", err.Error(), info))
	}
}
