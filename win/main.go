//go:build windows

package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "embed"
	"github.com/jchv/go-webview2"
)

func writeLog(msg string) {
	path := filepath.Join(os.TempDir(), "souhun.log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer file.Close()
	fmt.Fprintln(file, time.Now().Format(time.RFC3339), msg)
}

//go:embed ui.html
var uiHTML string

//go:embed assets/AppIcon.png
var iconPNG []byte

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--check" {
		os.Exit(runCheck())
	}
	openWindow()
}

func runCheck() int {
	report := detect()
	for _, row := range report.Rows {
		mark := "fail"
		if row.OK {
			mark = "ok"
		}
		fmt.Printf("%s=%s\n", row.ID, mark)
	}
	fmt.Printf("ready=%t\n", report.Ready)
	fmt.Printf("bots=%d\n", len(report.Bots))
	fmt.Printf("login_named=%t\n", report.Email != "" || report.Name != "")
	if report.Ready {
		return 0
	}
	return 1
}

func openWindow() {
	model := &appModel{}
	window := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     false,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  "搜魂",
			Width:  980,
			Height: 680,
			Center: true,
		},
	})
	if window == nil {
		writeLog("打不开窗口。这台电脑需要 Microsoft Edge WebView2。")
		os.Exit(1)
	}
	writeLog("窗口已打开")
	defer window.Destroy()
	window.Bind("appState", model.snapshot)
	window.Bind("appRefresh", model.refresh)
	window.Bind("appSelect", model.selectBot)
	window.Bind("appPull", model.pull)
	window.Bind("appWrite", model.write)
	window.Bind("appForget", model.forget)
	page := strings.Replace(uiHTML, "__ICON__", base64.StdEncoding.EncodeToString(iconPNG), 1)
	window.SetHtml(page)
	window.Run()
}
