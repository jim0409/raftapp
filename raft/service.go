package raft

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.etcd.io/etcd/raft/v3/raftpb"
)

type RaftNode struct {
	id       int
	join     bool
	kvport   int
	clusters []string

	proc       chan string
	confc      chan raftpb.ConfChange
	leaderaddr string
	wt         int // wait time to close
}

func (r *RaftNode) RunRaftNode() {
	var kvs *kvstore
	getSnapshot := func() ([]byte, error) { return kvs.getSnapshot() }

	// id 使用 uint 應該沒問題，要修改 newRaftNode
	commitC, errorC, snapshotterReady := newRaftNode(r.id, r.clusters, r.join, getSnapshot, r.proc, r.confc)

	kvs = newKVStore(<-snapshotterReady, r.proc, commitC, errorC)

	if r.join {
		r.regist()
	}

	// the key-value http handler will propose updates to raft
	// serveHttpKVAPI(kvs, r.kvport, r.confc, errorC)
	InitKVAPI(kvs, r.kvport, r.confc, errorC)
}

func (r *RaftNode) Close() {
	if r.join {
		r.unregist() // may consider to unregistr ?
		time.Sleep(time.Duration(r.wt) * time.Second)
	}
	close(r.proc)  // prposeC
	close(r.confc) // confChangeC
}

func InitRaftNode(id int, kvport int, clusters []string, join bool, leaderaddr string, wt int) *RaftNode {
	return &RaftNode{
		id:         id,
		kvport:     kvport,
		clusters:   clusters,
		join:       join,
		proc:       make(chan string),
		confc:      make(chan raftpb.ConfChange),
		leaderaddr: leaderaddr,
		wt:         wt,
	}
}

func (r *RaftNode) regist() {
	peeradrr := r.clusters[len(r.clusters)-1]
	body := strings.NewReader(peeradrr)
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/node/%d", r.leaderaddr, r.id), body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func (r *RaftNode) unregist() {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%v/node/%d", r.leaderaddr, r.id), nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
