package api

import (
	"context"

	pb "github.com/ozonmp/cnm-purchase-api/pkg/cnm-purchase-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (o *purchaseAPI) RemovePurchaseV1(
	ctx context.Context,
	req *pb.RemovePurchaseV1Request,
) (*pb.RemovePurchaseV1Response, error) {
	log.Debug().Msg("request is [" + req.String() + "]")

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("RemovePurchaseV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	found, err := o.repo.RemovePurchase(ctx, req.GetPurchaseId())
	if err != nil {
		log.Error().Err(err).Msg("error during RemovePurchaseV1")

		return nil, status.Error(codes.Internal, "")
	}

	return &pb.RemovePurchaseV1Response{Found: found}, nil
}
