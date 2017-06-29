package model

type Tag struct {
	Tag   string `storm:"id"`
	Paths []string
}
