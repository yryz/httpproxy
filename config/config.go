package config

import (
	"encoding/json"
	"io/ioutil"
	"os/user"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Listen     string `json:"listen"`
	SsServer   string `json:"ss_server"`
	SsCipher   string `json:"ss_cipher"`
	SsPassword string `json:"ss_password"`
}

var Conf = &Config{Listen: "127.0.0.1:6666"}

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	configFile := usr.HomeDir + "/.httpproxy/config.json"

	log.Info("load config: ", configFile)
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// log.Debug(hex.Dump(buf))

	err = json.Unmarshal(buf, Conf)
	if err != nil {
		log.Fatal(err)
	}
}
