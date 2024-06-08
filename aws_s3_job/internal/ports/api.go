package ports

import (
	"context"
)

type Api interface {
	ProcessS3Files(ctx context.Context, bucketName string) error
}
