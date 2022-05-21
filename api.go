package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Client interface {
	ZoneAPI
	RecordAPI
}

type ClientConfiguration struct {
	APIToken   string
	apiUrl     string
	HttpClient *http.Client
}

type ErrorResponse struct {
	Code    int
	Message string
}

func (er ErrorResponse) Error() string {
	return fmt.Sprintf("%v: %s", er.Code, er.Message)
}

type requestBuilder struct {
	*http.Request
	client *clientImpl
}

func (rb requestBuilder) AddQueryParams(params url.Values) requestBuilder {
	if len(params) == 0 {
		return rb
	}
	if rb.URL.RawQuery == "" {
		rb.URL.RawQuery = params.Encode()
	} else {
		rb.URL.RawQuery += "&" + params.Encode()
	}
	return rb
}

func (rb requestBuilder) Send() (*http.Response, error) {
	rb.Header.Add("Auth-API-Token", rb.client.token)

	response, err := rb.client.client.Do(rb.Request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		responseErr := ErrorResponse{
			Code:    response.StatusCode,
			Message: response.Status,
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return response, responseErr
		}

		raw := make(map[string]interface{})
		err = json.Unmarshal(body, &raw)
		if err != nil && raw["error"] != nil {
			errorString, ok := raw["error"].(string)
			if ok {
				responseErr.Message = errorString
			}
		}
		return response, responseErr
	}

	return response, nil
}

func (rb requestBuilder) WritePlain(body io.ReadCloser) requestBuilder {
	rb.Header.Add("Content-Type", "text/plain")
	rb.Body = body
	return rb
}

func (rb requestBuilder) ReadPlain() (string, error) {
	rb.Header.Add("Accept", "text/plain")

	response, err := rb.Send()
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (rb requestBuilder) WriteJSON(body interface{}) requestBuilder {
	rb.Header.Add("Content-Type", "application/json")
	encoded, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	rb.Body = ioutil.NopCloser(bytes.NewBuffer(encoded))
	return rb
}

func (rb requestBuilder) ReadJSON(result interface{}) error {
	rb.Header.Add("Accept", "application/json")

	response, err := rb.Send()
	if err != nil {
		return err
	}
	defer func() {
		ioutil.ReadAll(response.Body)
		response.Body.Close()
	}()
	return json.NewDecoder(response.Body).Decode(&result)
}

func (rb requestBuilder) JSON(body interface{}, result interface{}) error {
	if body != nil {
		rb.WriteJSON(body)
	}
	return rb.ReadJSON(result)
}

type clientImpl struct {
	token  string
	url    string
	client *http.Client
}

func NewClient(config ClientConfiguration) ZoneAPI {
	client := http.DefaultClient
	if config.HttpClient != nil {
		client = config.HttpClient
	}

	url := config.apiUrl
	if url == "" {
		url = "https://dns.hetzner.com/api/v1"
	}

	return &clientImpl{
		token:  config.APIToken,
		url:    url,
		client: client,
	}
}

func addPagedQueryParams(query url.Values, request PagedRequest) {
	if request.Page != 0 {
		query.Add("page", strconv.Itoa(request.Page))
	}
	if request.PerPage != 0 {
		query.Add("per_page", strconv.Itoa(request.PerPage))
	}
}

func (c clientImpl) request(method, path string) requestBuilder {
	absoluteUrl := fmt.Sprintf("%s/%s", c.url, path)
	req, _ := http.NewRequest(method, absoluteUrl, nil)
	return requestBuilder{req, &c}
}
