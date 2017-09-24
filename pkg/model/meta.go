package model

import "math"

type Meta struct {
	Posts   int
	PerPage int
	Pages   int
}

func NewMeta(posts, perPage int) *Meta {
	return &Meta{
		Posts:   posts,
		PerPage: perPage,
		Pages:   int(math.Ceil(float64(posts) / float64(perPage))),
	}
}
