package db

import "gorm.io/gorm"

type Node struct {
	gorm.Model
	Port int
	Addr string
}

type ImpNode interface {
	InsertDbRecord(int, string) (int, error)
	UpdateDbRecord(int, int, string) error
	GetClusterIps() ([]string, error)
	ReturnNodeInfo(int) (*Node, error)
}

var nodeTable = "nodes"

func (n *Node) TableName() string {
	return nodeTable
}

// InsertDBRecord : 在加入新的節點時，會主動去註冊節點的溝通 port 以及 peer-connectin address
func (db *Operation) InsertDbRecord(port int, addr string) (int, error) {
	n := &Node{
		Port: port,
		Addr: addr,
	}
	return int(n.ID), db.DB.Table(nodeTable).Create(n).Error
}

func (db *Operation) UpdateDbRecord(id int, port int, addr string) error {
	return db.DB.Table(nodeTable).Where(`id = ?`, id).Updates(&Node{
		Port: port,
		Addr: addr,
	}).Error
}

func (db *Operation) GetClusterIps() ([]string, error) {
	urls := []string{}
	if err := db.DB.Table(nodeTable).Select(`addr`).Where(`deleted_at is NULL`).Order(`id`).Scan(&urls).Error; err != nil {
		return nil, err
	}

	return urls, nil
}

func (db *Operation) ReturnNodeInfo(id int) (*Node, error) {
	n := &Node{}
	if err := db.DB.Table(nodeTable).Select(`*`).Where(`id = ? and deleted_at is NULL`, id).Scan(n).Error; err != nil {
		return nil, err
	}
	return n, nil
}
