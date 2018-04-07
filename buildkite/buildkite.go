package buildkite

import (
	"github.com/buildkite/go-buildkite/buildkite"
	"log"
)

type Client struct {
	apiClient *buildkite.Client
	org       string
}

func NewAPIClient(token, org string) *Client {
	config, err := buildkite.NewTokenConfig(token, false)
	if err != nil {
		log.Fatal("Buildkite client config failed: %s", err)
	}
	client := buildkite.NewClient(config.Client())

	return &Client{
		apiClient: client,
		org:       org,
	}
}

func (c *Client) AvailableAgents() int {
	agents, _, err := c.apiClient.Agents.List(c.org, nil)
	if err != nil {
		log.Fatal("API builds call failed: %s", err)
	}
	log.Println("Agents available ", len(agents))
	return len(agents)
}

func (c *Client) BuildList() []buildkite.Build {
	builds, _, err := c.apiClient.Builds.ListByOrg(c.org, &buildkite.BuildsListOptions{
		State: []string{"scheduled", "running", "canceling"},
	})
	if err != nil {
		log.Fatal("API builds call failed %s", err)
	}
	return builds
}
