package matcher

import (
	"bytes"
	"challenge/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidJSONResponse = errors.New("upstream: invlaid json response")
	ErrInvalidXMLResponse  = errors.New("upstream: invlaid xml response")
	ErrorRequestFailed     = errors.New("http request failed")
)

type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type FixtureClient struct {
	httpClient HTTPClient
	jsonParser ResponseParser[any]
	xmlParser  ResponseParser[any]
}

func NewFixtureClient(httpClient HTTPClient) *FixtureClient {
	return &FixtureClient{
		httpClient: httpClient,
		jsonParser: JsonParser{},
		xmlParser:  XMLParser{},
	}
}

func (f *FixtureClient) GetB() (*model.Input, error) {
	response, err := f.httpClient.Do(f.buildGetRequest("/source/b"))
	if err != nil {
		return nil, fmt.Errorf("%w : %w", ErrorRequestFailed, err)
	}

	defer response.Body.Close()
	output := &model.Input{}
	err = f.xmlParser.Parse(response.Body, output)
	return output, err
}

func (f *FixtureClient) GetA() (*model.Input, error) {
	response, err := f.httpClient.Do(f.buildGetRequest("/source/a"))
	if err != nil {
		return nil, fmt.Errorf("%w : %w", ErrorRequestFailed, err)
	}
	defer response.Body.Close()

	output := &model.Input{}
	err = f.jsonParser.Parse(response.Body, output)
	return output, err
}

func (f *FixtureClient) PostToSink(out model.Output) error {
	payloadBuf, err := json.Marshal(out)
	if err != nil {
		return err
	}
	request, _ := http.NewRequest(http.MethodPost, "/sink/a", bytes.NewBuffer(payloadBuf))

	response, err := f.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusAccepted {
		return errors.New("reqqust is unsuccessful")
	}
	return nil
}

func (*FixtureClient) buildGetRequest(path string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, path, nil)
	return request
}
