package payu

import (
	"fmt"
)

type (
	CreatePaymentResp struct {
		*Payment
		Links []Links `json:"links"`
	}

	CreateOrderResp struct {
		Status      *Status `json:"status,omitempty"`
		RedirectUri string  `json:"redirectUri,omitempty"`
		OrderId     string  `json:"orderId,omitempty"`
		ExtOrderId  string  `json:"extOrderId,omitempty"`
	}

	Status struct {
		StatusCode  string `json:"statusCode,omitempty"`
		Severity    string `json:"severity,omitempty"`
		Code        string `json:"code,omitempty"`
		CodeLiteral string `json:"codeLiteral,omitempty"`
		StatusDesc  string `json:"statusDesc,omitempty"`
	}

	ExecutePaymentResp struct {
		Intent       PaymentIntent `json:"intent"`
		Payer        *Payer        `json:"payer"`
		Transactions []Transaction `json:"transactions"`
		Links        []Links       `json:"links"`
	}

	ListPaymentsResp struct {
		Payments []Payment `json:"payments"`
	}
)

// CreatePayment creates a payment in Paypal
func (c *Client) CreatePayment(p Payment) (*CreatePaymentResp, error) {
	req, err := NewRequest("POST", fmt.Sprintf("%s/api/v2_1/orders", c.APIBase), p)
	if err != nil {
		return nil, err
	}

	v := &CreatePaymentResp{}

	err = c.SendWithAuth(req, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
func (c *Client) CreateOrder(p Order) (*CreateOrderResp, error) {
	req, err := NewRequest("POST", fmt.Sprintf("%s/api/v2_1/orders", c.APIBase), p)
	if err != nil {
		return nil, err
	}

	v := &CreateOrderResp{}

	err = c.SendWithAutWithouthRespError(req, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// ExecutePayment completes an approved Paypal payment that has been approved by the payer
func (c *Client) ExecutePayment(paymentID, payerID string, transactions []Transaction) (*ExecutePaymentResp, error) {
	req, err := NewRequest("POST", fmt.Sprintf("%s/payments/payment/%s/execute", c.APIBase, paymentID), struct {
		PayerID      string        `json:"payer_id"`
		Transactions []Transaction `json:"transactions"`
	}{
		payerID,
		transactions,
	})
	if err != nil {
		return nil, err
	}

	v := &ExecutePaymentResp{}

	err = c.SendWithAuth(req, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// GetPayment fetches a payment in Paypal
func (c *Client) GetPayment(id string) (*Payment, error) {
	req, err := NewRequest("GET", fmt.Sprintf("%s/payments/payment/%s", c.APIBase, id), nil)
	if err != nil {
		return nil, err
	}

	v := &Payment{}

	err = c.SendWithAuth(req, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// ListPayments retrieve payments resources from Paypal
func (c *Client) ListPayments(filter map[string]string) ([]Payment, error) {
	req, err := NewRequest("GET", fmt.Sprintf("%s/payments/payment/", c.APIBase), nil)
	if err != nil {
		return nil, err
	}

	if filter != nil {
		q := req.URL.Query()

		for k, v := range filter {
			q.Set(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}

	var v ListPaymentsResp

	err = c.SendWithAuth(req, &v)
	if err != nil {
		return nil, err
	}

	return v.Payments, nil
}
