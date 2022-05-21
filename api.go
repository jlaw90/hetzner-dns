package dns

import (
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

func (c clientImpl) request(method, path string, query url.Values, body io.Reader, result interface{}) (*http.Response, error) {
	absoluteUrl := fmt.Sprintf("%s/%s", c.url, path)
	if len(query) > 0 {
		absoluteUrl += "?" + query.Encode()
	}
	req, err := http.NewRequest(method, absoluteUrl, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Auth-API-Token", c.token)

	response, err := c.client.Do(req)

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

	if result != nil {
		delayedDecode := json.RawMessage{}
		err = json.NewDecoder(response.Body).Decode(&delayedDecode)
		ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err == nil {
			fmt.Printf("RAW RESPONSE: %+v\n", string(delayedDecode))
			err = json.Unmarshal(delayedDecode, &result)
		}
	}

	return response, nil
}
