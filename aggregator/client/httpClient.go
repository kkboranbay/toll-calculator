package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kkboranbay/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type HttpClient struct {
	Endpoint string
}

func NewHttpClient(endpoint string) Client {
	return &HttpClient{
		Endpoint: endpoint,
	}
}

func (c *HttpClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	invReq := types.GetInvoiceRequest{
		ObuID: int32(id),
	}

	b, err := json.Marshal(&invReq)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%s?obu=%d", c.Endpoint, "invoice", id)
	logrus.Infof("requesting get invoice -> %s", endpoint)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the service responded with non 200 status code %d", resp.StatusCode)
	}

	var inv types.Invoice
	if err := json.NewDecoder(resp.Body).Decode(&inv); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &inv, nil
}

func (c *HttpClient) Aggregate(ctx context.Context, aggReq *types.AggregatorRequest) error {
	b, err := json.Marshal(aggReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.Endpoint+"/aggregate", bytes.NewReader(b))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the service responded with non 200 status code %d", resp.StatusCode)
	}
	resp.Body.Close()
	return nil
}
