package model

import "math"

// Meta represents posts metadata model
type Meta struct {
	Posts   int
	PerPage int
	Pages   int
}

// NewMeta instantiate posts metadata with pre-calculated fields
func NewMeta(posts, perPage int) *Meta {
	return &Meta{
		Posts:   posts,
		PerPage: perPage,
		Pages:   int(math.Ceil(float64(posts) / float64(perPage))),
	}
}
