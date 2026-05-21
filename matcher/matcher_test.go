package matcher

import (
	"challenge/model"
	"reflect"
	"testing"
	"time"
)

func TestMatcherReceiveInputs(t *testing.T) {
	inA := make(chan model.Input)
	inB := make(chan model.Input)
	matcher := NewMatcher(inA, inB, 2*time.Second)

	t.Run("output received message sent on input A", func(t *testing.T) {

		go func() {
			defer close(inA)
			inA <- model.Input{ID: "1"}
			inA <- model.Input{ID: "2"}
			inA <- model.Input{ID: "3"}
		}()

		go func() {
			defer close(inB)
			inB <- model.Input{ID: "1"}
			inB <- model.Input{ID: "3"}
			inB <- model.Input{ID: "4"}
		}()

		res := make([]model.Output, 0)
		for r := range matcher.Match() {
			res = append(res, r)
		}
		expect := []model.Output{
			{ID: "1", Kind: model.Joined},
			{ID: "3", Kind: model.Joined},
			{ID: "2", Kind: model.Orphaned},
			{ID: "4", Kind: model.Orphaned},
		}
		if !reflect.DeepEqual(res, expect) {
			t.Errorf("actual %q, expected %q", res, expect)
		}
	})

}
