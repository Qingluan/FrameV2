package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/martinlindhe/notify"

	"github.com/Qingluan/FrameV2/icon"
	"github.com/Qingluan/merkur/config"
	"github.com/gen2brain/dlgs"
	"github.com/getlantern/systray"
)

// func testIfStart() bool {
// 	cmd := exec.Command(os.Args[0], "-book.ls")
// 	cmd.Env = os.Environ()
// 	data, err := cmd.Output()
// 	if err != nil {
// 		log.Println(err)
// 		// dlgs.Info("Pid", fmt.Sprintln(err))
// 		return false
// 	}
// 	println(string(data))
// 	// dlgs.Info("Pid", string(data))

// 	if strings.Contains(string(data), "json:unexpected") {
// 		return false
// 	}
// 	return true
// }

func execs(cmds string, std bool) (output string) {
	cmd := exec.Command("bash", "-c", cmds)
	cmd.Env = os.Environ()
	if strings.HasPrefix(cmds, "Kcpee") {
		msg := strings.Split(cmds, " ")
		// dlgs.Info("show", cmds)
		cmd = exec.Command(os.Args[0], msg[1:]...)
	}
	if std {
		var stdout bytes.Buffer
		// cmd.Stdout = &stdout
		err := cmd.Start()
		if err != nil {
			// dlgs.Info("Pid", fmt.Sprintln(err))
		}
		// time.Sleep(1 * time.Second)
		return fmt.Sprintf("%s", string(stdout.Bytes()))
	}
	data, err := cmd.Output()
	if err != nil {
		log.Println(err)
		// dlgs.Info("Pid", fmt.Sprintln(err))
	}
	output = strings.TrimSpace(string(data))
	println("output:", output)
	return

}

func RunGui(global func()) {
	onready := func() {
		OnReady(func() {
		})
	}
	Load()
	for k, v := range Configs.Routes {
		NowConfig = v
		notify.Notify("Frame2", "Info", "switch to "+k, "")
		break
	}
	go Listen()
	systray.Run(onready, OnExit)
	return
}

func OnReady(global func()) {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Kcpee ")
	// systray.SetTooltip("点击切换线路")
	switchg := systray.AddMenuItem("Routes", "switch route")
	setGlobal := systray.AddMenuItem("Global Mode", "set global mode")
	addItem := systray.AddMenuItem("Add Item", "add item ")
	clearItem := systray.AddMenuItem("Clear Item", "clear all item ")

	setListen := systray.AddMenuItem("set listen host", "set listen host port  ")

	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// Sets the icon of a menu item. Only available on Mac and Windows.
	mQuit.SetIcon(icon.Data)
	for {
		select {
		case <-clearItem.ClickedCh:
			Configs.Routes = make(map[string]string)
			Save(Configs)
		case <-setListen.ClickedCh:
			res, ok, err := dlgs.Entry("set listen", "set socks5 listen host:", Configs.ListenHost)
			if LogErr(err) {
				continue
			}
			if ok {
				Configs.ListenHost = strings.TrimSpace(res)
				TO_STOP = true
				RE_START = 1
			}
		case <-addItem.ClickedCh:
			res, ok, err := dlgs.Entry("Add item", "input vmess:// | ss:// | socks5:// | or order urls ", "")
			if err != nil || !ok {
				notify.Notify("FrameV2", "err", err.Error(), "")
			}
			if strings.HasPrefix(res, "http") {
				for k, v := range config.ParseOrding(res) {
					name := fmt.Sprintf("link-%d", len(Configs.Routes)+k)
					if strings.HasPrefix(v, "vmess") {
						if cfg, err := ParseVmessUri(v); err != nil {
							LogErr(err)
							continue
						} else {
							name = cfg.ID
						}
					} else if strings.HasPrefix(v, "ssr") {
						if cfg, err := ParseVmessUri(v); err != nil {
							LogErr(err)
							continue
						} else {
							name = cfg.ID
						}
					}

					Configs.Routes[name] = v
				}
				Save(Configs)
			} else {
				res2, ok, err := dlgs.Entry("mark name", "may be some country or some type:[default:"+fmt.Sprintf("link-%d", len(Configs.Routes))+"]", fmt.Sprintf("link-%d", len(Configs.Routes)))

				if err != nil || !ok {
					notify.Notify("FrameV2", "err", err.Error(), "")
				}
				Configs.Routes[res2] = res
				Save(Configs)
			}

		case <-switchg.ClickedCh:
			// items := []string{"Global Mode", "Stop Kcp", "Auto Mode", "Flow Mode"}
			items := []string{}
			for k := range Configs.Routes {
				items = append(items, k)
			}

			item, _, err := dlgs.List("FrameV2", "Select items:", items)
			NowConfig = Configs.Routes[item]
			if err != nil {
				// panic(err)
				notify.Alert("FrameV2", "Error info:", err.Error(), "")
			} else {
				notify.Notify("FrameV2", "switch to", item, "")
			}

			// if !s {
			// 	os.Exit(0)
			// }
		case <-setGlobal.ClickedCh:
			global()
		case <-mQuit.ClickedCh:
			// execs("Kcpee -book.stop ", false)
			out := execs("ps aux | grep ProxyAnyWhere | grep -v grep |` awk '{ print $2} '| xargs kill -9", false)
			notify.Notify("FrameV2", "Info", out, "")
			os.Exit(0)
			break
		}

	}
}

func OnExit() {
	// clean up here

	notify.Notify("FrameV2", "exit FrameV2", "this app exit!!", "")
}

func LogErr(err error) bool {
	if err != nil {
		notify.Alert("FrameV2", "Error info:", err.Error(), "")
		return true
	}
	return false
}

func Save(m ConfigK) {

	config := filepath.Join(Home, "FrameV2.json")
	buf, err := json.Marshal(&m)
	if LogErr(err) {
		return
	}
	err = ioutil.WriteFile(config, buf, os.ModePerm)
	if LogErr(err) {
		return
	}
	notify.Notify("FrameV2", "save configs", "Save => "+config, "")
}

func Load() {
	if _, err := os.Stat(Home); err != nil {
		// nameb,_ := exec.Command("bash", "-c", "whoami").Output()
		// name := string(nameb)
		// _,err exec.Command("bash", "-c", fmt.Sprintf("sudo mkdir /etc/ProxyAnyWhere && sudo chown -R %s:%s /etc/ProxyAnyWhere ",name, name))
		if err := os.MkdirAll(Home, os.ModePerm); err != nil {
			LogErr(err)
			os.Exit(1)
			return
		}
	}
	config := filepath.Join(Home, "FrameV2.json")
	_, err := os.Stat(config)
	if err != nil {
		Save(Configs)

		LogErr(err)
	}
	buf, err := ioutil.ReadFile(config)
	json.Unmarshal(buf, &Configs)
}
