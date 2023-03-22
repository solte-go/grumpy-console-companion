package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	gcc "grumpy-console-companion/sotle-go/proto"
	"math/rand"
	"net"
	"sync"
)

type API struct {
	gcc.UnimplementedGCCServer
	address    string
	quotes     map[string][]string
	mu         sync.Mutex
	grpcServer *grpc.Server
}

func New(address string) (*API, error) {
	var opts []grpc.ServerOption
	q := make(map[string][]string)
	q["greeting"] = []string{"The day is longer when you're bored", "Coffee is the only liquid you need"}

	api := &API{
		address:    address,
		quotes:     q,
		grpcServer: grpc.NewServer(opts...),
	}
	api.grpcServer.RegisterService(&gcc.GCC_ServiceDesc, api)
	return api, nil
}

func (api *API) Start() error {
	api.mu.Lock()
	defer api.mu.Unlock()

	listen, err := net.Listen("tcp", api.address)
	if err != nil {
		return err
	}
	return api.grpcServer.Serve(listen)
}

func (api *API) Stop() {
	api.mu.Lock()
	defer api.mu.Unlock()

	api.grpcServer.Stop()
}

func (api *API) GetGCC(ctx context.Context, req *gcc.GetReq) (*gcc.GetResp, error) {
	var (
		author string
		quotes []string
	)
	if author == "" {
		for author, quotes = range api.quotes {
			break
		}
	} else {
		author = req.Thoughts
		var ok bool
		quotes, ok = api.quotes[req.Thoughts]
		if !ok {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("author %q not found", req.Thoughts))
		}
	}
	return &gcc.GetResp{
		Topic:   author,
		Thought: quotes[rand.Intn(len(quotes))],
	}, nil
}
