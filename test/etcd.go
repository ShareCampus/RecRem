package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	cfg := clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
		Username:    "root",
		Password:    "@Wh060030",
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		fmt.Println("client connect failed")
	}
	defer cli.Close()
	timeout := time.Duration(10 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	resp, err := cli.Put(ctx, "sample_key", "sample_value")
	respg, err := cli.Get(ctx, "sample_key")
	fmt.Println(respg)
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			log.Fatalf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			log.Fatalf("client-side error: %v", err)
		default:
			log.Fatalf("bad cluster endpoints, which are not etcd servers: %v", err)
		}
		// use the response
		fmt.Println(resp)
	}

}
