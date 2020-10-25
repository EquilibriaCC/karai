package util

import (
	"fmt"
	"strconv"
	"time"
)

func UnixTimeStampNano() string {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	return timestamp
}

func Handle(msg string, err error) {
	if err != nil {
		fmt.Printf(Brightred+"\n%s: %s"+white, msg, err)
	}
}
