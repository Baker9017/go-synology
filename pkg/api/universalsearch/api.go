package universalsearch

import (
	"context"

	"github.com/synology-community/go-synology/pkg/models"
)

type Api interface {
	Search(ctx context.Context, keyword string, from int, size int) (*models.SearchData, error)
}
