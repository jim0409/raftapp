package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"raft-app/conf"
	"raft-app/db"
	"raft-app/raft"
	"syscall"
)

/*
1. start service with InitNodeConfig
2. throw input variables ndoe_addr, node_port into mysql
3. query mysql records, if records length is greater than 1, add join(true) flag
4. else return `id` as the started node `id`
*/

var (
	raftnode     *raft.RaftNode
	path         = flag.String("config", "./conf/app.dev.ini", "config location")
	checkcommit  = flag.Bool("version", false, "burry code for check version")
	gitcommitnum string
	defaultWT    = 5
)

func checkComimit() {
	fmt.Println(gitcommitnum)
}

func Init() error {
	flag.Parse()
	join := false

	// if there is a needed to check git commit num ... print it out
	if *checkcommit {
		checkComimit()
		os.Exit(1)
	}

	// read config and pass variables ...
	cfg, err := conf.InitConfig(*path)
	if err != nil {
		return err
	}

	opdb, err := db.NewDBConfiguration(cfg.DbUser, cfg.DbPassword, "mysql", cfg.DbName, cfg.DbPort, cfg.DbHost).NewDBConnection()
	if err != nil {
		return err
	}

	if cfg.ID == 0 {
		aid, err := opdb.InsertDbRecord(cfg.HttpPort, cfg.PeerAddr)
		if err != nil {
			return err
		}
		cfg.ID = aid

	} else {
		node, err := opdb.ReturnNodeInfo(cfg.ID)
		if err != nil {
			return err
		}
		cfg.HttpPort = node.Port
		cfg.PeerAddr = node.Addr
	}

	if cfg.HttpPort == 0 {
		return fmt.Errorf("no http port was provided!")
	}

	if cfg.PeerAddr == "" {
		return fmt.Errorf("no peer addr was provided!")
	}

	if cfg.LeaderAddr == "" {
		return fmt.Errorf("lack of leader address")
	}

	if cfg.WaitToClose == 0 {
		cfg.WaitToClose = defaultWT
	}

	clusters, err := opdb.GetClusterIps()
	if err != nil {
		return err
	}

	if len(clusters) > 1 {
		join = true
	}

	raftnode = raft.InitRaftNode(cfg.ID, cfg.HttpPort, clusters, join, cfg.LeaderAddr, cfg.WaitToClose)

	return nil
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic err: ", err)
		}
	}()

	err := Init()
	if err != nil {
		panic(err)
	}

	raftnode.RunRaftNode()
	defer raftnode.Close()

	gracefulShutdown()
}

// gracefulShutdown: handle the worker connection
func gracefulShutdown() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	<-done
}
