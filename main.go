package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"raft-app/conf"
	"raft-app/db"
	"raft-app/raft"
	"raft-app/router"
	"syscall"

	"github.com/gin-gonic/gin"
)

var (
	raftnode     *raft.RaftNode
	path         = flag.String("config", "./conf/app.dev.ini", "config location")
	checkcommit  = flag.Bool("version", false, "burry code for check version")
	gitcommitnum string
	defaultWT    = 5
	localip      string
)

func declareIp() error {
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {

			return err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP.To4()
				if !ip.IsLoopback() && ip != nil {
					localip = ip.String()
					return nil
				}
			}
		}
	}
	return nil
}

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

	// 如果有指定的 PeerAddr 則使用指定，否則使用本地第一張網卡的ip
	if cfg.PeerAddr == "" {
		if err = declareIp(); err != nil {
			return err
		}
		cfg.PeerAddr = fmt.Sprintf("http://%v:2379", localip)
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

	route := gin.Default()
	router.ApiRouter(route)

	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HttpPort),
		Handler: route,
	}

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(fmt.Sprintf("http listen : %v\n", err))
			panic(err)
		}
	}()

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
