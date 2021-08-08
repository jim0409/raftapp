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
	if r.id != 1 {
		r.join = true
	}

	var kvs *kvstore
	getSnapshot := func() ([]byte, error) { return kvs.getSnapshot() }

	// id 使用 uint 應該沒問題，要修改 newRaftNode
	commitC, errorC, snapshotterReady := newRaftNode(r.id, r.clusters, r.join, getSnapshot, r.proc, r.confc)

	kvs = newKVStore(<-snapshotterReady, r.proc, commitC, errorC)

	if err := r.regist(); err != nil {
		fmt.Println("something wrong while registing! ", err)
	}

	// the key-value http handler will propose updates to raft
	InitKVAPI(kvs, r.kvport, r.confc, errorC)
}

func (r *RaftNode) Close() {
	// may consider to unregistr ?
	if err := r.unregist(); err != nil {
		panic(err)
	}
	time.Sleep(time.Duration(r.wt) * time.Second)
	close(r.proc)  // prposeC
	close(r.confc) // confChangeC
}

func InitRaftNode(id int, kvport int, clusters []string, leaderaddr string, wt int) *RaftNode {
	return &RaftNode{
		id:         id,
		kvport:     kvport,
		clusters:   clusters,
		join:       false,
		proc:       make(chan string),
		confc:      make(chan raftpb.ConfChange),
		leaderaddr: leaderaddr,
		wt:         wt,
	}
}

func (r *RaftNode) regist() error {
	if r.id == 1 {
		return nil
	}

	peeraddr := r.clusters[r.id-1]
	body := strings.NewReader(peeraddr)
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/node/%d", r.leaderaddr, r.id), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (r *RaftNode) unregist() error {
	// id 1 still possible to leave
	// if r.id == 1 {
	// 	return nil
	// }

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%v/node/%d", r.leaderaddr, r.id), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
