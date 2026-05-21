package matcher

import (
	"challenge/model"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
)

type ResponseParser[T any] interface {
	Parse(body io.ReadCloser, model T) error
}

type JsonParser struct{}

func (p JsonParser) Parse(body io.ReadCloser, output any) (any error) {
	err := json.NewDecoder(body).Decode(output)
	if err != nil {
		return fmt.Errorf("%w", ErrInvalidJSONResponse)
	}
	return nil
}

type Msg struct {
	XMLName xml.Name `xml:"msg"`
	ID      *ID      `xml:"id,omitempty"`
	Done    *Done    `xml:"done,omitempty"`
}

type ID struct {
	Value string `xml:"value,attr"`
}

type Done struct{}

type XMLParser struct{}

func (p XMLParser) Parse(body io.ReadCloser, output any) (any error) {
	var m Msg
	in, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to body for xml: %w", err)
	}
	if err := xml.Unmarshal(in, &m); err != nil {
		return fmt.Errorf("%w : %w", ErrInvalidXMLResponse, err)
	}
	out := output.(*model.Input)
	switch {
	case m.ID != nil:
		out.ID = m.ID.Value
		out.Status = "ok"
	case m.Done != nil:
		out.Status = "done"
	}
	return nil
}
