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
	// fields 参数指定返回的元数据，需显式请求时间相关字段才会在结果中返回
	fields := `["SYNOMDExtension","SYNOMDFSName","SYNOMDFSSize","SYNOMDIsDir","SYNOMDOwnerUserID","SYNOMDPath","SYNOMDSharePath","SYNOMDContentModificationDate","SYNOMDFSCreationDate","SYNOMDLastUsedDate"]`

	return api.Post[models.SearchData](c.client, ctx, &models.SearchRequest{
		Keyword:          keyword,
		SearchWeightList: searchWeightList,
		QueryID:          queryID,
		Agent:            agent,
		Fields:           fields,
		From:             from,
		Size:             size,
	}, methods.Search)
}
