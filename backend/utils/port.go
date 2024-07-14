package utils

import (
	"fmt"
	"net"
	"strconv"
)

func PortCheck(port int) bool {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port)))

	if err != nil {
		return false
	}
	defer func(l net.Listener) {
		_ = l.Close()
	}(l)
	return true
}
