package model

type Input struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

type OutputKind string

const (
	Joined   OutputKind = "joined"
	Orphaned OutputKind = "orphaned"
)

type Output struct {
	Kind OutputKind `json:"kind"`
	ID   string     `json:"id"`
}
