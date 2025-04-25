package grpcservice

import (
	"author-service/proto/book"
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

func (c *BookGRPCClient) SaveAuthor(ctx context.Context, req *book.AuthorData) (*book.BookResponse, error) {
	res, err := c.client.ReceiveAuthor(ctx, req)
	if err != nil {
		return &book.BookResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return res, nil
}

func (c *BookGRPCClient) DeleteAuthor(ctx context.Context, authorId uint) (*book.BookResponse, error) {
	req := &book.DeleteData{Id: uint32(authorId)}
	res, err := c.client.DeleteAuthor(ctx, req)
	if err != nil {
		return &book.BookResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return res, nil
}
