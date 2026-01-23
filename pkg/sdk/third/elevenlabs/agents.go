package elevenlabs

import (
	"context"
	"fmt"
)

// ListAgents lists all agents.
func (c *Client) ListAgents(ctx context.Context, params *ListAgentsParams) (*ListAgentsResp, error) {
	var respData ListAgentsResp
	req := c.request(ctx)
	if params != nil {
		if params.PageSize > 0 {
			req.SetQueryParam("page_size", fmt.Sprintf("%d", params.PageSize))
		}
		if params.Search != "" {
			req.SetQueryParam("search", params.Search)
		}
		if params.Archived {
			req.SetQueryParam("archived", "true")
		}
		if params.ShowOnlyOwned {
			req.SetQueryParam("show_only_owned_agents", "true")
		}
		if params.SortBy != "" {
			req.SetQueryParam("sort_by", params.SortBy)
		}
		if params.SortDirection != "" {
			req.SetQueryParam("sort_direction", params.SortDirection)
		}
		if params.Cursor != "" {
			req.SetQueryParam("cursor", params.Cursor)
		}
	}
	resp, err := req.SetResult(&respData).Get("/v1/convai/agents")
	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}
	return &respData, nil
}

// GetAgent retrieves a specific agent by ID.
func (c *Client) GetAgent(ctx context.Context, agentID string) (*Agent, error) {
	var respData Agent
	resp, err := c.request(ctx).SetResult(&respData).Get("/v1/convai/agents/" + agentID)
	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}
	return &respData, nil
}

// CreateAgent creates a new agent.
func (c *Client) CreateAgent(ctx context.Context, reqBody *CreateAgentRequest) (string, error) {
	var result struct {
		AgentID string `json:"agent_id"`
	}
	resp, err := c.request(ctx).SetBody(reqBody).SetResult(&result).Post("/v1/convai/agents/create")
	if err := handleResponse(resp, err); err != nil {
		return "", err
	}
	return result.AgentID, nil
}

// UpdateAgent updates an existing agent.
func (c *Client) UpdateAgent(ctx context.Context, agentID string, reqBody map[string]any) error {
	resp, err := c.request(ctx).SetBody(reqBody).Patch("/v1/convai/agents/" + agentID)
	return handleResponse(resp, err)
}

// DeleteAgent deletes an agent.
func (c *Client) DeleteAgent(ctx context.Context, agentID string) error {
	resp, err := c.request(ctx).Delete("/v1/convai/agents/" + agentID)
	return handleResponse(resp, err)
}
