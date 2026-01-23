package elevenlabs

import (
	"context"
	"fmt"
)

type ListToolsResp struct {
	Tools []Tool `json:"tools"`
	Pagination
}

// ListTools lists all available tools.
func (c *Client) ListTools(ctx context.Context, pageSize int, cursor string) (*ListToolsResp, error) {
	var respData ListToolsResp
	req := c.request(ctx)
	if pageSize > 0 {
		req.SetQueryParam("page_size", fmt.Sprintf("%d", pageSize))
	}
	if cursor != "" {
		req.SetQueryParam("cursor", cursor)
	}
	resp, err := req.SetResult(&respData).Get("/v1/convai/tools")
	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}
	return &respData, nil
}

// GetTool retrieves a specific tool by ID.
func (c *Client) GetTool(ctx context.Context, toolID string) (*Tool, error) {
	var respData Tool
	resp, err := c.request(ctx).SetResult(&respData).Get("/v1/convai/tools/" + toolID)
	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}
	return &respData, nil
}

// CreateTool creates a new tool.
func (c *Client) CreateTool(ctx context.Context, tool *Tool) (string, error) {
	var result struct {
		ToolID string `json:"tool_id"`
	}
	resp, err := c.request(ctx).SetBody(tool).SetResult(&result).Post("/v1/convai/tools")
	if err := handleResponse(resp, err); err != nil {
		return "", err
	}
	return result.ToolID, nil
}

// UpdateTool updates an existing tool.
func (c *Client) UpdateTool(ctx context.Context, toolID string, tool *Tool) error {
	resp, err := c.request(ctx).SetBody(tool).Patch("/v1/convai/tools/" + toolID)
	return handleResponse(resp, err)
}

// DeleteTool deletes a tool.
func (c *Client) DeleteTool(ctx context.Context, toolID string) error {
	resp, err := c.request(ctx).Delete("/v1/convai/tools/" + toolID)
	return handleResponse(resp, err)
}
