package model

// Tag represents post tag model
type Tag struct {
	Tag   string `storm:"id"`
	Paths []string
}
