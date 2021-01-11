package builders

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
)

type LabelBuilder struct {
	Label *models.Label
}

func NewLabelBuilder() *LabelBuilder {
	label := &models.Label{
		BoardId: 1,
		Name:    "Default Label Name",
		Color:   255,
	}
	return &LabelBuilder{Label: label}
}

func (b *LabelBuilder) Build() *models.Label {
	return b.Label
}

func (b *LabelBuilder) WithName(name string) *LabelBuilder {
	b.Label.Name = name
	return b
}

func (b *LabelBuilder) WithBoard(id int) *LabelBuilder {
	b.Label.BoardId = id
	return b
}

func (b *LabelBuilder) WithColor(color uint32) *LabelBuilder {
	b.Label.Color = color
	return b
}
