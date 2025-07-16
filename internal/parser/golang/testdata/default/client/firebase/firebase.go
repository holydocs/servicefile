package firebase

import "github.com/samber/do"

/*
service:requests Firebase
description: Handles push notifications
technology:firebase
proto:http
*/
type Client struct{}

func NewClient(i *do.Injector) *Client {
	return &Client{}
}

func (c *Client) Request() error {
	return nil
}
