package utils

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
)

type Config struct {
	Server        string `json:"server"`
	Password      string `json:"password"`
	Method        string `json:"method"`
	ServerPort    int    `json:"server_port"`
	Timeout       int    `json:"timeout"`
	LocalPort     int    `json:"local_port"`
	LocalAddress  string `json:"local_address"`
	Protocol      string `json:"protocol"`
	Obfs          string `json:"obfs"`
	ConfigType    string `json:"conf_type"`
	ObfsParam     string `json:"obfs-param"`
	ProtocolParam string `json:"protocol-param"`
	OptionID      int    `json:"aid"`
	OptUUID       string `json:"uuid"`
	ID            string `json:"ps"`
}

func b64decode(a string) (o string, err error) {
	var ierr error
	if strings.Contains(a, "_") {
		a = strings.ReplaceAll(a, "_", "/")
	}

	if strings.Contains(a, "-") {
		a = strings.ReplaceAll(a, "-", "+")
	}

	for i := 0; i < 4; i++ {
		dat, err := base64.StdEncoding.DecodeString(strings.TrimSpace(a))
		if err != nil {
			a += "="
			ierr = err
			continue
		}
		return string(dat), nil
	}
	return "", ierr
}

func ParseVmessUri(u string) (cfg Config, err error) {
	var dats string
	if strings.HasPrefix(u, "vmess://") {
		u = u[8:]
	}
	if dats, err = b64decode(u); err != nil {
		return
	}
	s := make(map[string]interface{})
	err = json.Unmarshal([]byte(dats), &s)

	if err != nil {
		log.Println("json parse err:", err, u)

		return
	}
	// fmt.Println(s)
	if host, ok := s["host"]; ok {
		cfg.Server = host.(string)
	}

	if netAddr, ok := s["add"]; ok {
		cfg.ObfsParam = netAddr.(string)
	}

	if ports, ok := s["port"]; ok {
		cfg.ServerPort = int(ports.(float64))

	}
	if aids, ok := s["aid"]; ok {
		cfg.OptionID = int(aids.(float64))

	}
	if proto, ok := s["net"]; ok {
		cfg.Protocol = proto.(string)
		if cfg.Protocol == "ws" {
			cfg.ProtocolParam = s["path"].(string)
		}
	}

	if uid, ok := s["id"]; ok {
		cfg.OptUUID = uid.(string)
	}

	if sectype, ok := s["type"]; ok {
		cfg.Obfs = sectype.(string)
	}

	if name, ok := s["ps"]; ok {
		cfg.ID = name.(string)
	}
	cfg.ConfigType = "vmess"
	return
}
