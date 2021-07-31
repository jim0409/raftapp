package conf

import "github.com/go-ini/ini"

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
	return Conf, nil
}

// func InitDB() (db.OPDB, error) {
// 	mysqlAddr := "mysql"
// 	mysqlPort := "3306"
// 	mysqlOpDB := "raft"
// 	mysqlUsr := "raft"
// 	mysqUsrPwd := "raft"

// 	return db.NewDBConfiguration(mysqlUsr, mysqUsrPwd, "mysql", mysqlOpDB, mysqlPort, mysqlAddr).NewDBConnection()
// }
