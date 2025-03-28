package main

import (
	"github.com/NainVictorin1/homework2/Internal/data"
)

type TemplateData struct {
	Title      string
	HeaderText string
	FormErrors map[string]string
	FormData   map[string]string

	Todos     []*data.Todo
	Journals  []*data.Journal
	Feedbacks []*data.Feedback
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:      "Default Title",
		HeaderText: "Default HeaderText",
		FormErrors: map[string]string{},
		FormData:   map[string]string{},
	}
}
