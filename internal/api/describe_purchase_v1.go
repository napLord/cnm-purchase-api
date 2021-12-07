package api

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/cnm-purchase-api/internal/repo"
	pb "github.com/ozonmp/cnm-purchase-api/pkg/cnm-purchase-api"
)

func (o *purchaseAPI) DescribePurchaseV1(
	ctx context.Context,
	req *pb.DescribePurchaseV1Request,
) (*pb.DescribePurchaseV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribePurchaseV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	purchase, err := o.repo.DescribePurchase(ctx, req.PurchaseId)
	if err != nil {
		log.Error().Err(err).Msg("DescribePurchaseV1 -- failed")

		if errors.Is(err, repo.ErrNoPurchases) {
			return nil, status.Error(codes.Internal, "no purchases")
		}

		return nil, status.Error(codes.Internal, "")
	}

	if purchase == nil {
		log.Debug().Uint64("purchaseId", req.PurchaseId).Msg("purchase not found")
		totalPurchaseNotFound.Inc()

		return nil, status.Error(codes.NotFound, "purchase not found")
	}

	log.Debug().Msg("DescribePurchaseV1 - success")

	return &pb.DescribePurchaseV1Response{
		Value: &pb.Purchase{
			Id:       purchase.ID,
			TotalSum: purchase.TotalSum,
		},
	}, nil
}
