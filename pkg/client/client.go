package client

import (
	"context"
	"google.golang.org/grpc"
	gcc "grumpy-console-companion/sotle-go/proto"
	"time"
)

type Client struct {
	client gcc.GCCClient
	conn   *grpc.ClientConn
}

func New(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{
		client: gcc.NewGCCClient(conn),
		conn:   conn,
	}, nil
}

func (c *Client) QOTD(ctx context.Context, wantAuthor string) (author, quite string, err error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
	}
	resp, err := c.client.GetGCC(ctx, &gcc.GetReq{Thoughts: wantAuthor})
	if err != nil {
		return "", "", err
	}
	return resp.Topic, resp.Thought, nil
}
