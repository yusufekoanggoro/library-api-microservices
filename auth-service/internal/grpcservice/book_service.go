package grpcservice

import (
	"auth-service/proto/book"
	"context"
	"log"

	"google.golang.org/grpc"
)

type BookGRPCClient struct {
	client book.BookServiceClient
}

func NewBookGRPCClient(bookServiceAddr string) *BookGRPCClient {
	conn, err := grpc.Dial(bookServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Book Service: %v", err)
	}

	client := book.NewBookServiceClient(conn)
	return &BookGRPCClient{client: client}
}

func (c *BookGRPCClient) SaveUser(ctx context.Context, req *book.UserData) (*book.BookResponse, error) {
	res, err := c.client.ReceiveUser(ctx, req)
	if err != nil {
		return &book.BookResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return res, nil
}

func (c *BookGRPCClient) DeleteUser(ctx context.Context, userId uint) (*book.BookResponse, error) {
	req := &book.DeleteData{Id: uint32(userId)}
	res, err := c.client.DeleteUser(ctx, req)
	if err != nil {
		return &book.BookResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return res, nil
}
