package api

import (
	"context"

	"github.com/ozonmp/cnm-purchase-api/internal/repo"
	cnm_purchase_api "github.com/ozonmp/cnm-purchase-api/pkg/cnm-purchase-api"
	pb "github.com/ozonmp/cnm-purchase-api/pkg/cnm-purchase-api"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (o *purchaseAPI) ListPurchasesV1(
	ctx context.Context,
	req *pb.ListPurchasesV1Request,
) (*pb.ListPurchasesV1Response, error) {
	log.Debug().Msg("request is [" + req.String() + "]")

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("ListPurchasesV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	purchases, nextCursor, err := o.repo.ListPurchases(ctx, req.GetLimit(), req.GetCursor())

	if err != nil {
		log.Error().Err(err).Msg("error during ListPurchasesV1")

		if errors.Is(err, repo.ErrNoPurchases) {
			return nil, status.Error(codes.Internal, "no purchases")
		}

		return nil, status.Error(codes.Internal, "")
	}

	purchasesResp := make([]*cnm_purchase_api.Purchase, 0, len(purchases))

	for _, v := range purchases {
		purchasesResp = append(purchasesResp, &cnm_purchase_api.Purchase{Id: v.ID, TotalSum: v.TotalSum})
	}

	return &pb.ListPurchasesV1Response{Items: purchasesResp, Cursor: nextCursor}, nil
}
