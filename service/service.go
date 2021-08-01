package service

import (
	"encoding/json"
	"net/http"
	"raft-app/raft"
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
	}
	err = json.Unmarshal(b, kv)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
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
	}

	b, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	raft.RetrieveKVApi().AddNode(uint64(iid), string(b))
	c.JSON(http.StatusOK, "ok")
	return
}

func DelNode(c *gin.Context) {
	id := c.Param("nd")
	iid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	raft.RetrieveKVApi().DelNode(uint64(iid))
	c.JSON(http.StatusOK, "ok")
	return
}
