package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConf(t *testing.T) {
	path := "./app.dev.ini"
	config, err := InitConfig(path)
	assert.Nil(t, err)
	assert.Equal(t, 2380, config.HttpPort)
	assert.Equal(t, "http://raftapp.dev-pds.svc.cluster.local:2379", config.PeerAddr)
	assert.Equal(t, "http://raftapp.dev-pds.svc.cluster.local:2380", config.LeaderAddr)
	assert.Equal(t, 2, config.WaitToClose)
	assert.Equal(t, "dev", config.Env)
	assert.Equal(t, "raft", config.DbName)
	assert.Equal(t, "mysql", config.DbHost)
	assert.Equal(t, "3306", config.DbPort)
	assert.Equal(t, "raft", config.DbUser)
	assert.Equal(t, "raft", config.DbPassword)
	assert.Equal(t, false, config.DbLogEnable)
	assert.Equal(t, 300, config.DbMaxConnect)
	assert.Equal(t, 10, config.DbIdleConnect)
}
