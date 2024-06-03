// internal/eksapi.go
package interfaces

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/eks"
)

type EKSAPI interface {
	ListClusters(ctx context.Context, params *eks.ListClustersInput, optFns ...func(*eks.Options)) (*eks.ListClustersOutput, error)
}
