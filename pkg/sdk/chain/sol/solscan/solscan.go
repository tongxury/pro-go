package solscan

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	endpoint string
}

func NewClient() *Client {

	return &Client{
		endpoint: "https://api-v2.solscan.io",
	}
}

func (t *Client) GetTokenInfo(ctx context.Context, addressBase58 string) (*TokenInfo, error) {

	response, err := resty.New().R().SetContext(ctx).
		//SetHeader("Sol-Aut", "UmSPaiffdRR9U3EB9dls0fKN75px=RjAauGv3gfN").
		//SetHeader("Cookie", "cf_clearance=bESb85h7enw7hXirzB5HmhYWJHfEzsTFVpK_GrPAiHo-1721123898-1.0.1.1-jpS0H7WP32isLx2nlni1NtC9NvKKIiYHVv93KQHmrDA7GzVYOeKbUaeuKxLNGw1SMt434pZiEzAe9UDOBjXt2g; _ga_PS3V7B7KV0=GS1.1.1721123959.1.0.1721123959.0.0.0; _ga=GA1.1.621415896.1721123959").
		SetHeader("Origin", "https://solscan.io").
		SetQueryParam("address", addressBase58).
		Get(t.endpoint + "/v2/account")

	var resp SolscanResponse[Account]
	err = json.Unmarshal(response.Body(), &resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf(response.String())
	}

	return &resp.Data.TokenInfo, nil
}
