package api

import (
	"context"

	pb "github.com/ozonmp/cnm-purchase-api/pkg/cnm-purchase-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (o *purchaseAPI) ListPurchasesV1(
	ctx context.Context,
	req *pb.ListPurchasesV1Request,
) (*pb.ListPurchasesV1Response, error) {
	log.Debug().Msg("request is [" + req.String() + "]")

	return nil, status.Error(codes.Unimplemented, "method is unimplemented")
}
