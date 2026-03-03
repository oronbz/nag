package reminders

import ekreminders "github.com/BRO3886/go-eventkit/reminders"

type Client struct {
	ek *ekreminders.Client
}

func New() (*Client, error) {
	ek, err := ekreminders.New()
	if err != nil {
		return nil, err
	}
	return &Client{ek: ek}, nil
}
