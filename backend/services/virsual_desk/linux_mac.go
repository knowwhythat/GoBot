//go:build !windows

package virsual_desk

import "errors"

func StartInVirtualDesk(id string) error {
	return errors.New("非windows系统不支持虚拟桌面")
}
func TerminateVirtualDeskFlow(id string) error {
	return errors.New("非windows系统不支持虚拟桌面")
}
