package portalsfi

type Client struct {
	endpoint string
}

func NewClient(endpoint string) *Client {

	//Bearer 7d488b51-86d9-473c-bf51-0ddc839711e3
	return &Client{
		endpoint: endpoint,
	}
}
