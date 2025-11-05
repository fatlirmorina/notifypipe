package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Client wraps the Docker client
type Client struct {
	cli *client.Client
	ctx context.Context
}

// NewClient creates a new Docker client
func NewClient(socketPath string) (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		cli: cli,
		ctx: context.Background(),
	}, nil
}

// ListContainers lists all containers
func (c *Client) ListContainers() ([]types.Container, error) {
	return c.cli.ContainerList(c.ctx, container.ListOptions{All: true})
}

// GetContainer gets a container by ID
func (c *Client) GetContainer(id string) (types.ContainerJSON, error) {
	return c.cli.ContainerInspect(c.ctx, id)
}

// Close closes the Docker client
func (c *Client) Close() error {
	return c.cli.Close()
}

// GetClient returns the underlying Docker client
func (c *Client) GetClient() *client.Client {
	return c.cli
}

// GetContext returns the context
func (c *Client) GetContext() context.Context {
	return c.ctx
}
