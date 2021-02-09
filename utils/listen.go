package utils

import (
	"log"
	"net"
	"strings"

	"github.com/Qingluan/merkur"
	"github.com/martinlindhe/notify"
)

var (
	TO_STOP               = false
	RE_START              = 0
	NowConfig             = ""
	Socks5ConnectedRemote = []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43}
	Configs               = Init()
)

type ConfigK struct {
	Routes     map[string]string `json:"routes"`
	ListenHost string            `json:"listen"`
}

func Init() ConfigK {
	return ConfigK{
		Routes:     make(map[string]string),
		ListenHost: "0.0.0.0:1080",
	}
}

func Listen() (err error) {
	ln, err := net.Listen("tcp", Configs.ListenHost)
	// if conn.ShowLog < 2 {
	// 	// utils.ColorL("Local Listen:", listenAddr)

	// }
	notify.Notify("FrameV2", "Start", Configs.ListenHost, "")

	for {
		if TO_STOP {
			break
		}
		// if conn.Role == "tester" && conn.GetAliveNum() > conn.Numconn {
		// 	time.Sleep(10 * time.Millisecond)
		// 	continue
		// }
		p1, err := ln.Accept()

		if err != nil {
			if !strings.Contains(err.Error(), "too many open files") {
				LogErr(err)
			}

			continue
		}
		go handleSocks5TcpAndUDP(p1)

	}
	if RE_START > 0 {
		Listen()
	}
	return
}

func handleSocks5TcpAndUDP(p1 net.Conn) {
	defer p1.Close()
	if err := Socks5HandShake(&p1); err != nil {
		// utils.ColorL("socks handshake:", err)
		return
	}

	_, host, _, err := GetLocalRequest(&p1)
	if err != nil {
		LogErr(err)
		return
	}
	// fmt.Println(string(raw))
	// if isUdp {

	// 	utils.ColorL("socks5 UDP-->", host)
	// } else {

	// 	utils.ColorL("socks5 -->", host)
	// }
	if err != nil {
		LogErr(err)
		return
	}
	handleBody(p1, host)
}

func PipeOne(dst, fr net.Conn) (err error) {
	buf := make([]byte, 4096)
	n, err := fr.Read(buf)
	if err != nil {
		return
	}
	log.Println("BUF:", string(buf[:n]))
	_, err = dst.Write(buf[:n])
	return
}

func handleBody(p1 net.Conn, host string) {
	if NowConfig != "" {
		dialer := merkur.NewProxyDialer(NowConfig)
		p2, err := dialer.Dial("tcp", host)

		if err != nil {
			notify.Notify("Frame2", "error", err.Error(), "")
			return
		}
		if p2 == nil {
			notify.Notify("Frame2", "error", "p2 is not connected", "")
			return
		}
		if p1 != nil && p2 != nil {
			log.Println("connect ->", host)
			_, err = p1.Write(Socks5ConnectedRemote)
			if err != nil {
				LogErr(err)
				return
			}
			// PipeOne(p2, p1)
			Pipe(p1, p2)
		} else {
			notify.Notify("Frame2", "err", "connection is not connected!", "")
		}
	} else {
		notify.Notify("Frame2", "error", "No set Proxy item", "")
		return
	}
}
