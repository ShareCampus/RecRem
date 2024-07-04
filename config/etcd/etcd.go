package etcd

import (
	"context"
	"log"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	EtcdIns *clientv3.Client
	once    sync.Once
)

// InitEtcd initializes the etcd client as a singleton.
func InitEtcd() {
	once.Do(func() {
		cfg := clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "root",
			Password:    "@Wh060030",
		}

		cli, err := clientv3.New(cfg)
		if err != nil {
			log.Fatalf("Failed to connect to etcd: %v", err)
		}

		// Test the connection
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, err = cli.Get(ctx, "test")
		if err != nil && err != rpctypes.ErrKeyNotFound {
			log.Fatalf("Failed to connect to etcd: %v", err)
		}

		EtcdIns = cli
	})
}

// GetEtcdClient returns the etcd client instance.
func GetEtcdClient() *clientv3.Client {
	if EtcdIns == nil {
		log.Fatal("Etcd client is not initialized. Call InitEtcd() first.")
	}
	return EtcdIns
}
