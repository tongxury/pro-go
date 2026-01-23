package elevenlabs

import (
	"context"
	"fmt"
	"store/pkg/sdk/conv"
)

// ListConversations lists past conversations.
func (c *Client) ListConversations(ctx context.Context, params *ListConversationsParams) (*ListConversationsResp, error) {
	var respData ListConversationsResp
	req := c.request(ctx)
	if params != nil {
		if params.AgentID != "" {
			req.SetQueryParam("agent_id", params.AgentID)
		}
		if params.CallStartAfter > 0 {
			req.SetQueryParam("call_start_after", fmt.Sprintf("%d", params.CallStartAfter))
		}
		if params.CallStartBefore > 0 {
			req.SetQueryParam("call_start_before", fmt.Sprintf("%d", params.CallStartBefore))
		}
		if params.PageSize > 0 {
			req.SetQueryParam("page_size", fmt.Sprintf("%d", params.PageSize))
		}
		if params.Cursor != "" {
			req.SetQueryParam("cursor", params.Cursor)
		}
	}
	resp, err := req.SetResult(&respData).Get("/v1/convai/conversations")
	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}
	return &respData, nil
}

// GetConversation retrieves a conversation detail by ID.
func (c *Client) GetConversation(ctx context.Context, conversationID string) (*ConversationDetail, error) {
	var respData ConversationDetail
	resp, err := c.request(ctx).SetResult(&respData).Get("/v1/convai/conversations/" + conversationID)
	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}
	return &respData, nil
}

// DeleteConversation deletes a conversation.
func (c *Client) DeleteConversation(ctx context.Context, conversationID string) error {
	resp, err := c.request(ctx).Delete("/v1/convai/conversations/" + conversationID)
	return handleResponse(resp, err)
}

// GetSignedURL generates a signed URL for a conversation.
func (c *Client) GetSignedURL(ctx context.Context, agentID string) (string, error) {
	var result struct {
		SignedURL string `json:"signed_url"`
	}
	resp, err := c.request(ctx).SetQueryParam("agent_id", agentID).Get("/v1/convai/conversation/get-signed-url")
	if err := handleResponse(resp, err); err != nil {
		return "", err
	}
	if err := conv.J2S(resp.Body(), &result); err != nil {
		return "", err
	}
	return result.SignedURL, nil
}

// GenerateConversationToken generates a token for private agent conversation.
func (c *Client) GenerateConversationToken(ctx context.Context, agentID string) (string, error) {
	var result struct {
		Token string `json:"token"`
	}
	resp, err := c.request(ctx).SetQueryParam("agent_id", agentID).Get("/v1/convai/conversation/token")
	if err := handleResponse(resp, err); err != nil {
		return "", err
	}
	if err := conv.J2S(resp.Body(), &result); err != nil {
		return "", err
	}
	return result.Token, nil
}
