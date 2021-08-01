package raft

import "go.etcd.io/etcd/raft/v3/raftpb"

var (
	handler *httpKVAPI
	// once    sync.Once
)

func RetrieveKVApi() *httpKVAPI {
	// once.Do(func() {
	// 	handler = &httpKVAPI{}
	// })
	return handler
}

type httpKVAPI struct {
	store       *kvstore
	confChangeC chan<- raftpb.ConfChange
}

func (h *httpKVAPI) Lookup(key string) string {
	if v, ok := h.store.Lookup(key); ok {
		return v
	}
	return ""
}

func (h *httpKVAPI) Propose(key string, value string) {
	h.store.Propose(key, value)
}

func (h *httpKVAPI) AddNode(nodeId uint64, url string) {
	cc := raftpb.ConfChange{
		Type:    raftpb.ConfChangeAddNode,
		NodeID:  nodeId,
		Context: []byte(url),
	}
	h.confChangeC <- cc
}

func (h *httpKVAPI) DelNode(nodeId uint64) {
	cc := raftpb.ConfChange{
		Type:   raftpb.ConfChangeRemoveNode,
		NodeID: nodeId,
	}
	h.confChangeC <- cc
}

func InitKVAPI(kv *kvstore, port int, confChangeC chan<- raftpb.ConfChange, errorC <-chan error) {
	handler = &httpKVAPI{
		store:       kv,
		confChangeC: confChangeC,
	}
}
