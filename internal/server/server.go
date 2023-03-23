package server

import (
	"context"
	"google.golang.org/grpc"
	"grumpy-console-companion/sotle-go/pkg/model"
	"grumpy-console-companion/sotle-go/pkg/storage"
	gcc "grumpy-console-companion/sotle-go/proto"
	"math/rand"
	"net"
	"sync"
)

type API struct {
	gcc.UnimplementedGCCServer
	address    string
	mu         sync.Mutex
	grpcServer *grpc.Server
	contract   contract
}

type contract interface {
	GetAllTopics(ctx context.Context) ([]model.Topics, error)
	GetThoughtsOnTopic(ctx context.Context, topic string) ([]model.Thoughts, error)
}

func New(address string, st storage.Storage) (*API, error) {
	var opts []grpc.ServerOption

	api := &API{
		address:    address,
		grpcServer: grpc.NewServer(opts...),
		contract:   st,
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
	if req.Thoughts != "" {
		result, err := api.contract.GetAllTopics(ctx)
		if err != nil {
			return &gcc.GetResp{
				Topic:   "error",
				Thought: "I'm shorted out, contact the creator!",
			}, err
		}

		randomTopics := result[rand.Intn(len(result))].Name
		tghResult, err := api.contract.GetThoughtsOnTopic(ctx, randomTopics)
		if err != nil {
			return &gcc.GetResp{
				Topic:   "error",
				Thought: "I'm shorted out, contact the creator!",
			}, err
		}

		return &gcc.GetResp{
			Topic:   randomTopics,
			Thought: tghResult[rand.Intn(len(tghResult))].Phrase,
		}, nil
	}

	return &gcc.GetResp{
		Topic:   "None",
		Thought: "What do you want to discuss?",
	}, nil
}
