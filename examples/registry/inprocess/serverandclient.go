package main

import (
	"context"
	"flag"
	"log"

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/examples"
	"github.com/smallnest/rpcx/server"
)

var (
	addr     = flag.String("addr", "localhost:8972", "server address")
	etcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
	basePath = flag.String("base", "/rpcx_test", "prefix path")
)

func main() {
	flag.Parse()

	s := server.NewServer()
	addRegistryPlugin(s)

	s.RegisterName("Arith", new(examples.Arith), "")

	go func() {
		s.Serve("tcp", *addr)
	}()

	d := client.NewInprocessDiscovery()
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &examples.Args{
		A: 10,
		B: 20,
	}

	for i := 0; i < 10000; i++ {

		reply := &examples.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)

	}
}

func addRegistryPlugin(s *server.Server) {

	r := client.InprocessClient
	s.Plugins.Add(r)
}
