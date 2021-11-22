//nolint:dupl, revive
package api

import (
	"context"

	pb "github.com/ozonmp/cnm-purchase-api/pkg/cnm-purchase-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (o *purchaseAPI) CreatePurchaseV1(
	ctx context.Context,
	req *pb.CreatePurchaseV1Request,
) (*pb.CreatePurchaseV1Response, error) {
	log.Debug().Msg("request is [" + req.String() + "]")

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("CreatePurchaseV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	purchaseID, err := o.repo.CreatePurchase(ctx, req.GetTotalSum())
	if err != nil {
		log.Error().Err(err).Msg("error during CreatePurchaseV1")

		return nil, status.Error(codes.Internal, "")
	}

	return &pb.CreatePurchaseV1Response{PurchaseId: purchaseID}, nil
}
