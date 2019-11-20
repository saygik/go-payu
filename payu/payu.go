package payu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	// APIBaseSandBox points to the sandbox (for testing) version of the API
	APIBaseSandBox = "https://secure.snd.payu.com"

	// APIBaseLive points to the live version of the API
	APIBaseLive = "https://secure.snd.payu.com"
)

type (

	// Client represents a Paypal REST API Client
	Client struct {
		client   *http.Client
		ClientID string
		Secret   string
		APIBase  string
		Token    *TokenResp
	}

	// ErrorResponse is used when a response contains errors
	// maps to error object
	ErrorResponse struct {
		// HTTP response that caused this error
		Response *http.Response `json:"-"`

		Name            string        `json:"name"`
		DebugID         string        `json:"debug_id"`
		Message         string        `json:"message"`
		InformationLink string        `json:"information_link"`
		Details         []ErrorDetail `json:"details"`
	}

	// ErrorDetails map to error_details object
	ErrorDetail struct {
		Field string `json:"field"`
		Issue string `json:"issue"`
	}

	// TokenResp maps to the API response for the /oauth2/token endpoint
	TokenResp struct {
		Scope     string    `json:"scope"`        // "https://api.paypal.com/v1/payments/.* https://api.paypal.com/v1/vault/credit-card https://api.paypal.com/v1/vault/credit-card/.*",
		Token     string    `json:"access_token"` // "EEwJ6tF9x5WCIZDYzyZGaz6Khbw7raYRIBV_WxVvgmsG",
		Type      string    `json:"token_type"`   // "Bearer",
		AppID     string    `json:"app_id"`       // "APP-6XR95014BA15863X",
		ExpiresIn int       `json:"expires_in"`   // 28800
		ExpiresAt time.Time `json:"expires_at"`
	}
)

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v\nDetails: %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message, r.Details)
}

// NewClient returns a new Client struct
func NewClient(clientID, secret, APIBase string) *Client {
	return &Client{
		&http.Client{},
		clientID,
		secret,
		APIBase,
		nil,
	}
}

// NewRequest constructs a request. If payload is not empty, it will be
// marshalled into JSON
func NewRequest(method, url string, payload interface{}) (*http.Request, error) {
	var buf io.Reader
	if payload != nil {
		var b []byte
		b, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}
	return http.NewRequest(method, url, buf)
}

// GetAcessToken request a new access token from Paypal
//grant_type=client_credentials&client_id=369722&client_secret=aa3272c123df4a43b3d26e4b4794c0c1
func (c *Client) GetAccessToken() (*TokenResp, error) {
	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Add("client_secret", c.Secret)
	data.Add("grant_type", "client_credentials")

	//	buf := bytes.NewBuffer([]byte("grant_type=client_credentials"))
	buf := bytes.NewBufferString(data.Encode())
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/pl/standard/user/oauth/authorize"), buf)

	if err != nil {
		return nil, err
	}
	//	req.SetBasicAuth(c.ClientID, c.Secret)
	//	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	t := TokenResp{}
	err = c.Send(req, &t)
	if err == nil {
		t.ExpiresAt = time.Now().Add(time.Duration(t.ExpiresIn/2) * time.Second)
	}

	log.Println(t.Type, ": ", t.Token)
	return &t, err
}

// Send makes a request to the API, the response body will be
// unmarshaled into v, or if v is an io.Writer, the response will
// be written to it without decoding
func (c *Client) Send(req *http.Request, v interface{}) error {
	// Set default headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en_US")

	// Default values for headers
	if req.Header.Get("Content-type") == "" {
		req.Header.Set("Content-type", "application/json")
	}

	log.Println(req.Method, ": ", req.URL)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c := resp.StatusCode; c < 200 || c > 299 {
		errResp := &ErrorResponse{Response: resp}
		data, err := ioutil.ReadAll(resp.Body)

		if err == nil && len(data) > 0 {
			json.Unmarshal(data, errResp)
		}

		return errResp
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) SendWithoutResError(req *http.Request, v interface{}) error {
	// Set default headers
	//	req.Header.Set("Accept", "application/json")
	//	req.Header.Set("Accept-Language", "en_US")

	// Default values for headers
	if req.Header.Get("Content-type") == "" {
		req.Header.Set("Content-type", "application/json")
	}

	log.Println(req.Method, ": ", req.URL)

	var DefaultTransport http.RoundTripper = &http.Transport{}
	resp, err := DefaultTransport.RoundTrip(req)

	//	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//bodyBytes, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//bodyString := string(bodyBytes)
	//fmt.Printf("%q", bodyString)

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SendWithAuth makes a request to the API and apply OAuth2 header automatically.
// If the access token soon to be expired, it will try to get a new one before
// making the main request
func (c *Client) SendWithAuth(req *http.Request, v interface{}) error {
	if (c.Token == nil) || (c.Token.ExpiresAt.Before(time.Now())) {
		resp, err := c.GetAccessToken()
		if err != nil {
			return err
		}

		c.Token = resp
	}
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)

	return c.Send(req, v)
}
func (c *Client) SendWithAutWithouthRespError(req *http.Request, v interface{}) error {
	if (c.Token == nil) || (c.Token.ExpiresAt.Before(time.Now())) {
		resp, err := c.GetAccessToken()
		if err != nil {
			return err
		}

		c.Token = resp
	}
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)

	return c.SendWithoutResError(req, v)
}
