//go:build windows

package sys_tray

import (
	"context"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func CreateTray(ctx context.Context, icon []byte) {
	menu := map[string][2]string{
		"Show": {"显示", "显示窗口"},
		"Hide": {"隐藏", "隐藏窗口"},
		"Quit": {"退出", "退出程序"},
	}

	go func() {
		systray.Run(func() {
			systray.SetIcon(icon)
			systray.SetTitle("GoBot")
			systray.SetTooltip("GoBot")
			mShow := systray.AddMenuItem(menu["Show"][0], menu["Show"][1])
			mHide := systray.AddMenuItem(menu["Hide"][0], menu["Hide"][1])
			systray.AddSeparator()
			mQuit := systray.AddMenuItem(menu["Quit"][0], menu["Quit"][1])

			for {
				select {
				case <-mShow.ClickedCh:
					runtime.WindowShow(ctx)
				case <-mHide.ClickedCh:
					runtime.WindowHide(ctx)
				case <-mQuit.ClickedCh:
					systray.Quit()
					runtime.Quit(ctx)
				}
			}
		}, nil)
	}()
}
