package builders

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
)

type ListBuilder struct {
	List *models.TaskList
}

func NewListBuilder() *ListBuilder {
	list := &models.TaskList{
		BoardId:  1,
		Title:    "Default List Title",
		Position: 1,
	}
	return &ListBuilder{List: list}
}

func (b *ListBuilder) Build() *models.TaskList {
	return b.List
}

func (b *ListBuilder) WithTitle(title string) *ListBuilder {
	b.List.Title = title
	return b
}

func (b *ListBuilder) WithBoard(id int) *ListBuilder {
	b.List.BoardId = id
	return b
}

func (b *ListBuilder) WithPos(id int) *ListBuilder {
	b.List.Position = id
	return b
}
