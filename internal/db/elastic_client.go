package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	esv8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"github.com/rodrigogmartins/log-processor/internal/service"
)

type ElasticSearchClient struct {
	Client *esv8.Client
}

func NewElasticSearchClient(addresses []string) (*ElasticSearchClient, error) {
	cfg := esv8.Config{
		Addresses: addresses,
	}
	client, err := esv8.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ElasticSearchClient{
		Client: client,
	}, nil
}

func (c *ElasticSearchClient) Index(ctx context.Context, index string, id string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.Client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document ID %s: %s", id, res.String())
	}

	return nil
}

func (c *ElasticSearchClient) SearchLogs(ctx context.Context, index string, query map[string]interface{}, size int) ([]service.Log, error) {
	body := map[string]interface{}{
		"query": query,
		"size":  size,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.Search(
		c.Client.Search.WithContext(ctx),
		c.Client.Search.WithIndex(index),
		c.Client.Search.WithBody(bytes.NewReader(data)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching logs: %s", res.String())
	}

	var r struct {
		Hits struct {
			Hits []struct {
				Source service.Log `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	dataResp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(dataResp, &r); err != nil {
		return nil, err
	}

	logs := make([]service.Log, len(r.Hits.Hits))
	for i, hit := range r.Hits.Hits {
		logs[i] = hit.Source
	}

	return logs, nil
}
