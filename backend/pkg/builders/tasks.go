package builders

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
)

type TaskBuilder struct {
	Task *models.Task
}

func NewTaskBuilder() *TaskBuilder {
	task := &models.Task{
		ListId:      1,
		Title:       "Default List Title",
		Description: "Default List Desription",
		Position:    1,
		Datetimes: &models.Datetimes{
			Created:  1,
			Updated:  1,
			Accessed: 1,
		},
	}
	return &TaskBuilder{Task: task}
}

func (b *TaskBuilder) Build() *models.Task {
	return b.Task
}

func (b *TaskBuilder) WithTitle(title string) *TaskBuilder {
	b.Task.Title = title
	return b
}

func (b *TaskBuilder) WithList(id int) *TaskBuilder {
	b.Task.ListId = id
	return b
}

func (b *TaskBuilder) WithPos(id int) *TaskBuilder {
	b.Task.Position = id
	return b
}

func (b *TaskBuilder) WithDate(c, u, a int64) *TaskBuilder {
	date := &models.Datetimes{
		Created:  c,
		Updated:  u,
		Accessed: a,
	}
	b.Task.Datetimes = date
	return b
}
