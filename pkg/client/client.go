package client

type Config struct {
	server string
}

type Client struct {
	config Config
}

func NewClient(server string) *Client {
	return &Client{
		config: Config{
			server: server,
		},
	}
}
