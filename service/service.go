package service

import (
	"encoding/json"
	"net/http"
	"raftapp/raft"
	"strconv"

	"github.com/gin-gonic/gin"
)

type KeyValue struct {
	Key   string
	Value string
}

func GetKey(c *gin.Context) {
	id := c.Param("ky")
	value := raft.RetrieveKVApi().Lookup(id)

	c.JSON(http.StatusOK, value)
	return
}

func PutKey(c *gin.Context) {
	kv := &KeyValue{}
	b, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	err = json.Unmarshal(b, kv)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	raft.RetrieveKVApi().Propose(kv.Key, kv.Value)

	c.JSON(http.StatusOK, "ok")
	return
}

func AddNode(c *gin.Context) {
	id := c.Param("nd")
	iid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	b, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// TOOD: 判斷，如果收到的 raft 新增節點 id 與目前自身節點 id 相同，則不需要執行此步驟
	raft.RetrieveKVApi().AddNode(uint64(iid), string(b))
	c.JSON(http.StatusOK, "ok")
	return
}

func DelNode(c *gin.Context) {
	id := c.Param("nd")
	iid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	raft.RetrieveKVApi().DelNode(uint64(iid))
	c.JSON(http.StatusOK, "ok")
	return
}
