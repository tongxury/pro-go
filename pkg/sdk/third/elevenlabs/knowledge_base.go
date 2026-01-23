package elevenlabs

import (
	"context"
	"fmt"
)

// ListKnowledgeBase lists all documents in the knowledge base.
func (c *Client) ListKnowledgeBase(ctx context.Context, pageSize int, cursor string) (*ListKnowledgeBaseResp, error) {
	var respData ListKnowledgeBaseResp
	req := c.request(ctx)
	if pageSize > 0 {
		req.SetQueryParam("page_size", fmt.Sprintf("%d", pageSize))
	}
	if cursor != "" {
		req.SetQueryParam("cursor", cursor)
	}
	resp, err := req.SetResult(&respData).Get("/v1/convai/knowledge-base")
	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}
	return &respData, nil
}

// GetKnowledgeDocument retrieves a specific document metadata by ID.
func (c *Client) GetKnowledgeDocument(ctx context.Context, documentID string) (*KnowledgeDocument, error) {
	var respData KnowledgeDocument
	resp, err := c.request(ctx).SetResult(&respData).Get("/v1/convai/knowledge-base/" + documentID)
	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}
	return &respData, nil
}

// GetKnowledgeDocumentContent retrieves the content of a specific document.
func (c *Client) GetKnowledgeDocumentContent(ctx context.Context, documentID string) (string, error) {
	resp, err := c.request(ctx).Get("/v1/convai/knowledge-base/" + documentID + "/content")
	if err := handleResponse(resp, err); err != nil {
		return "", err
	}
	return resp.String(), nil
}

// GetDependentAgents lists agents that use a specific knowledge base document.
func (c *Client) GetDependentAgents(ctx context.Context, documentID string, pageSize int, cursor string) (*ListDependentAgentsResp, error) {
	var respData ListDependentAgentsResp
	req := c.request(ctx)
	if pageSize > 0 {
		req.SetQueryParam("page_size", fmt.Sprintf("%d", pageSize))
	}
	if cursor != "" {
		req.SetQueryParam("cursor", cursor)
	}
	resp, err := req.SetResult(&respData).Get("/v1/convai/knowledge-base/" + documentID + "/dependent-agents")
	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}
	return &respData, nil
}

// DeleteKnowledgeDocument deletes a document from the knowledge base.
func (c *Client) DeleteKnowledgeDocument(ctx context.Context, documentID string) error {
	resp, err := c.request(ctx).Delete("/v1/convai/knowledge-base/" + documentID)
	return handleResponse(resp, err)
}

// AddToKnowledgeBase adds a new document to the knowledge base (via URL or direct text).
// For file uploads, a separate multipart implementation would be needed.
func (c *Client) AddToKnowledgeBase(ctx context.Context, name string, url string) (string, error) {
	var result struct {
		DocumentID string `json:"id"`
	}
	reqBody := map[string]string{
		"name": name,
		"url":  url,
	}
	resp, err := c.request(ctx).SetBody(reqBody).SetResult(&result).Post("/v1/convai/knowledge-base")
	if err := handleResponse(resp, err); err != nil {
		return "", err
	}
	return result.DocumentID, nil
}
