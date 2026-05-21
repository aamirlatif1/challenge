package matcher_test

import (
	"bytes"
	"challenge/matcher"
	"challenge/model"
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestClientForRequestA(t *testing.T) {
	httpClient := &FakeHTTPClient{}
	client := matcher.NewFixtureClient(httpClient)

	t.Run("get A request compose properly", func(t *testing.T) {
		httpClient.configure(`{ "status": "ok", "id": "XXXXX" }`, http.StatusOK, nil)
		client.GetA()
		if httpClient.request.Method != http.MethodGet {
			t.Error("requst is not GET type")
		}
		if httpClient.request.URL.Path != "/source/a" {
			t.Errorf("\n got %q expected %q", httpClient.request.URL.Path, "/source/a")
		}
	})

	t.Run("json resonse is parsed", func(t *testing.T) {
		httpClient.configure(`{ "status": "ok", "id": "id1" }`, http.StatusOK, nil)

		a, err := client.GetA()

		if err != nil {
			t.Fatal("error shoul be nil", err)
		}

		if a.ID != "id1" {
			t.Errorf("\nid did not mtached, got %q expected %q", a.ID, "id1")
		}
		if a.Status != "ok" {
			t.Errorf("\nid did not mtached, got %q expected %q", a.Status, "ok")
		}
	})

	t.Run("parsed json done response", func(t *testing.T) {
		httpClient.configure(`{ "status": "done" }`, http.StatusOK, nil)

		a, err := client.GetA()

		if err != nil {
			t.Fatal("error shoul be nil", err)
		}

		if a.Status != "done" {
			t.Errorf("\nid did not mtached, got %q expected %q", a.Status, "done")
		}
	})

	t.Run("invalid json return error", func(t *testing.T) {
		httpClient.configure(`{{-}}`, http.StatusOK, nil)

		_, err := client.GetA()

		if err == nil || !errors.Is(err, matcher.ErrInvalidJSONResponse) {
			t.Error("should return invalid json response error")
		}

	})

	t.Run("request error to source A", func(t *testing.T) {
		httpClient.configure(`{{-}}`, http.StatusOK, errors.New("Invalid request"))

		_, err := client.GetA()

		if err == nil || !errors.Is(err, matcher.ErrorRequestFailed) {
			t.Error("should return error request failed")
		}

	})

}

func TestGetRequestForB(t *testing.T) {
	httpClient := &FakeHTTPClient{}
	client := matcher.NewFixtureClient(httpClient)

	t.Run("get A request compose properly", func(t *testing.T) {
		httpClient.configure(``, http.StatusOK, nil)
		client.GetB()
		if httpClient.request.Method != http.MethodGet {
			t.Error("requst is not GET type")
		}
		if httpClient.request.URL.Path != "/source/b" {
			t.Errorf("\n got %q expected %q", httpClient.request.URL.Path, "/source/b")
		}
	})

	t.Run("xml resonse is parsed", func(t *testing.T) {
		httpClient.configure(`<?xml version="1.0" encoding="UTF-8"?><msg><id value="id2"/></msg>`, http.StatusOK, nil)

		a, err := client.GetB()

		if err != nil {
			t.Fatal("error should be nil", err)
		}

		if a.ID != "id2" {
			t.Errorf("\nid did not mtached, got %q expected %q", a.ID, "id2")
		}
		if a.Status != "ok" {
			t.Errorf("\nid did not mtached, got %q expected %q", a.Status, "ok")
		}
	})

	t.Run("invalid xml return error", func(t *testing.T) {
		httpClient.configure(`<-->`, http.StatusOK, nil)

		_, err := client.GetB()

		if err == nil || !errors.Is(err, matcher.ErrInvalidXMLResponse) {
			t.Fatal("should return invalid response error")
		}

	})

}

func TestPostToSink(t *testing.T) {
	httpClient := &FakeHTTPClient{}
	client := matcher.NewFixtureClient(httpClient)

	t.Run("post to sink compose properly", func(t *testing.T) {
		httpClient.configure(``, http.StatusOK, nil)
		out := model.Output{
			ID:   "id1",
			Kind: model.Joined,
		}
		client.PostToSink(out)
		if httpClient.request.Method != http.MethodPost {
			t.Error("requst is not GET type")
		}
		if httpClient.request.URL.Path != "/sink/a" {
			t.Errorf("\n got %q expected %q", httpClient.request.URL.Path, "/sink/a")
		}
	})

	t.Run("post sent successfully", func(t *testing.T) {
		httpClient.configure(``, http.StatusOK, nil)
		out := model.Output{
			ID:   "id1",
			Kind: model.Joined,
		}
		client.PostToSink(out)

		expected := `{"kind":"joined","id":"id1"}`
		if httpClient.request.Method != http.MethodPost {
			t.Error("requst is not GET type")
		}
		if httpClient.body != expected {
			t.Errorf("\nbody did not matched, got %q expected %q", httpClient.body, expected)
		}
	})

}

/////////--------

type FakeHTTPClient struct {
	request  *http.Request
	response *http.Response
	body     string
	err      error
}

func (c *FakeHTTPClient) configure(resStr string, statusCode int, err error) {
	if err == nil {
		c.response = &http.Response{
			Body:       io.NopCloser(bytes.NewBufferString(resStr)),
			StatusCode: statusCode,
		}
	}
	c.err = err
}

func (c *FakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.request = req
	if c.request.Method == http.MethodPost {
		bytedata, _ := io.ReadAll(c.request.Body)
		c.body = string(bytedata)
	}
	return c.response, c.err
}
