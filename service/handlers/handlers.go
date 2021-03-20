package handlers

import (
	"context"

	pb "github.com/techxmind/keyservice/interface-defs"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.KeyServiceServer {
	return keyserviceService{}
}

type keyserviceService struct{}

func (s keyserviceService) Encrypt(ctx context.Context, in *pb.EncryptRequest) (*pb.Response, error) {
	var resp pb.Response
	return &resp, nil
}

func (s keyserviceService) EncryptBatch(ctx context.Context, in *pb.EncryptBatchRequest) (*pb.BatchResponse, error) {
	var resp pb.BatchResponse
	return &resp, nil
}

func (s keyserviceService) Decrypt(ctx context.Context, in *pb.DecryptRequest) (*pb.Response, error) {
	var resp pb.Response
	return &resp, nil
}

func (s keyserviceService) DecryptBatch(ctx context.Context, in *pb.DecryptBatchRequest) (*pb.BatchResponse, error) {
	var resp pb.BatchResponse
	return &resp, nil
}

func (s keyserviceService) Keys(ctx context.Context, in *pb.KeyRequest) (*pb.KeyResponse, error) {
	var resp pb.KeyResponse
	return &resp, nil
}

func (s keyserviceService) Ping(ctx context.Context, in *pb.Empty) (*pb.Response, error) {
	var resp pb.Response
	return &resp, nil
}
