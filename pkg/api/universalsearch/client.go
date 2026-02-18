package universalsearch

import (
	"context"
	"fmt"

	"github.com/synology-community/go-synology/pkg/api"
	"github.com/synology-community/go-synology/pkg/api/universalsearch/methods"
	"github.com/synology-community/go-synology/pkg/models"
)

type Client struct {
	client api.Api
}

func New(client api.Api) Api {
	return &Client{client: client}
}

func (c *Client) Search(ctx context.Context, keyword string, from int, size int) (*models.SearchData, error) {
	searchWeightList := `[{"field":"SYNOMDSearchFileName","weight":1}]`
	queryID := fmt.Sprintf(`"%d"`, from)
	agent := `"sus"`

	return api.Post[models.SearchData](c.client, ctx, &models.SearchRequest{
		Keyword:          keyword,
		SearchWeightList: searchWeightList,
		QueryID:          queryID,
		Agent:            agent,
		From:             from,
		Size:             size,
	}, methods.Search)
}
