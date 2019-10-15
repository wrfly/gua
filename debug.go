package gua

import (
	"encoding/json"
	"fmt"
)

var _EnableDebug = false

func debug(format string, a ...interface{}) {
	if _EnableDebug {
		fmt.Printf("[debug] "+format+"\n", a...)
	}
}

func debugJSON(format string, a interface{}) {
	if _EnableDebug {
		bs, _ := json.Marshal(a)
		fmt.Printf("[debug] "+format+"\n", bs)
	}
}
