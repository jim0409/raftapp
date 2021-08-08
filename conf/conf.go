package conf

import (
	"os"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
)

var Conf *Config

type Config struct {
	BaseConf
	DbConf
}

// BaseConf inlclude deatils server components
type BaseConf struct {
	ID          int    `ini:"ID"`
	HttpPort    int    `ini:"HttpPort"`    // http port
	PeerAddr    string `ini:"PeerAddr"`    // peer address
	LeaderAddr  string `ini:"LeaderAddr"`  // leader address
	WaitToClose int    `ini:"WaitToClose"` // wait time to close
	Env         string `ini:"Env"`         // 運行環境
}

// DbConf is for mysql settings
type DbConf struct {
	DbName        string `ini:"DbName"`
	DbHost        string `ini:"DbHost"`
	DbPort        string `ini:"DbPort"`
	DbUser        string `ini:"DbUser"`
	DbPassword    string `ini:"DbPassword"`
	DbLogEnable   bool   `ini:"DbLogEnable"`
	DbMaxConnect  int    `ini:"DbMaxConnect"`
	DbIdleConnect int    `ini:"DbIdleConnect"`
}

func InitConfig(path string) (*Config, error) {
	Conf = new(Config)
	if err := ini.MapTo(Conf, path); err != nil {
		return nil, err
	}

	if os.Getenv("k8s-env") != "" {
		id, addr := k8sConfig()
		Conf.ID = id
		// Conf.PeerAddr = fmt.Sprintf("http://%v:2379", addr)
		_ = addr
	}

	return Conf, nil
}

// unde k8s env, given some special parameters
func k8sConfig() (int, string) {
	h, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	id := strings.Split(h, "-")[1]

	iid, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}

	// 因為 k8s stateful set 從 0 開始
	return iid + 1, h
}
