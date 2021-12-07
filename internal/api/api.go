package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/ozonmp/cnm-purchase-api/internal/repo"

	pb "github.com/ozonmp/cnm-purchase-api/pkg/cnm-purchase-api"
)

var (
	totalPurchaseNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cnm_purchase_api_purchase_not_found_total",
		Help: "Total number of purchases that were not found",
	})
)

type purchaseAPI struct {
	pb.UnimplementedCnmPurchaseApiServiceServer
	repo repo.Repo
}

// NewPurchaseAPI returns api of cnm-purchase-api service
func NewPurchaseAPI(r repo.Repo) pb.CnmPurchaseApiServiceServer {
	return &purchaseAPI{repo: r}
}
