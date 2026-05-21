package matcher

import (
	"challenge/model"
	"time"
)

type Matcher struct {
	inputA        chan model.Input
	inputB        chan model.Input
	orphanTimeout time.Duration
}

func NewMatcher(inA, inB chan model.Input, orphanTimeout time.Duration) *Matcher {
	return &Matcher{
		inputA:        inA,
		inputB:        inB,
		orphanTimeout: orphanTimeout,
	}
}

func (m *Matcher) Match() <-chan model.Output {
	output := make(chan model.Output)

	go func() {
		defer close(output)
		pendingA := make(map[string]time.Time)
		pendingB := make(map[string]time.Time)
		openA, openB := true, true
		ticker := time.NewTicker(m.orphanTimeout / 2)
		defer ticker.Stop()
		for openA || openB {
			select {
			case a, ok := <-m.inputA:
				if !ok {
					m.inputA = nil
					openA = false
					continue
				}
				if _, found := pendingB[a.ID]; found {
					delete(pendingB, a.ID)
					output <- model.Output{ID: a.ID, Kind: model.Joined}
				} else {
					pendingA[a.ID] = time.Now()
				}
			case b, ok := <-m.inputB:
				if !ok {
					m.inputB = nil
					openB = false
					continue
				}
				if _, found := pendingA[b.ID]; found {
					delete(pendingA, b.ID)
					output <- model.Output{ID: b.ID, Kind: model.Joined}
				} else {
					pendingB[b.ID] = time.Now()
				}
			case now := <-ticker.C:
				for id, arrived := range pendingA {
					if now.Sub(arrived) >= m.orphanTimeout {
						delete(pendingA, id)
						output <- model.Output{ID: id, Kind: model.Orphaned}
					}
				}
				for id, arrived := range pendingB {
					if now.Sub(arrived) >= m.orphanTimeout {
						delete(pendingB, id)
						output <- model.Output{ID: id, Kind: model.Orphaned}
					}
				}
			}
		}
		for id := range pendingA {
			output <- model.Output{ID: id, Kind: model.Orphaned}
		}
		for id := range pendingB {
			output <- model.Output{ID: id, Kind: model.Orphaned}
		}
	}()

	return output
}
