package builders

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
)

type BoardBuilder struct {
	Board *models.Board
}

func NewBoardBuilder() *BoardBuilder {
	board := &models.Board{
		ProjectId: 1,
		OwnerId:   1,
		DefaultPermissions: &models.Permission{
			Read:  true,
			Write: true,
			Admin: false,
		},
		Datetimes: &models.Datetimes{
			Created:  1,
			Updated:  1,
			Accessed: 1,
		},
		Title: "Default Title",
	}
	return &BoardBuilder{Board: board}
}

func (b *BoardBuilder) Build() *models.Board {
	return b.Board
}

func (b *BoardBuilder) WithTitle(title string) *BoardBuilder {
	b.Board.Title = title
	return b
}

func (b *BoardBuilder) WithoutPerm() *BoardBuilder {
	b.Board.DefaultPermissions = nil
	return b
}

func (b *BoardBuilder) WithOwner(id int) *BoardBuilder {
	b.Board.OwnerId = id
	return b
}

func (b *BoardBuilder) WithProject(id int) *BoardBuilder {
	b.Board.ProjectId = id
	return b
}

func (b *BoardBuilder) WithPerm(r, w, a bool) *BoardBuilder {
	perm := &models.Permission{
		Read:  r,
		Write: w,
		Admin: a,
	}
	b.Board.DefaultPermissions = perm
	return b
}

func (b *BoardBuilder) WithDate(c, u, a int64) *BoardBuilder {
	date := &models.Datetimes{
		Created:  c,
		Updated:  u,
		Accessed: a,
	}
	b.Board.Datetimes = date
	return b
}
