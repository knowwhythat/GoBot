//go:build !windows

package virtual_desk

import "errors"
import "context"

func StartInVirtualDesk(ctx context.Context, id string) error {
	return errors.New("非windows系统不支持虚拟桌面")
}
func TerminateVirtualDeskFlow(ctx context.Context, id string) error {
	return errors.New("非windows系统不支持虚拟桌面")
}
func StopVisualHost() {
}
