package builders

import (
	"github.com/architectv/networking-course-project/backend/pkg/models"
)

type PermsBuilder struct {
	Perms *models.Permission
}

func NewPermsBuilder() *PermsBuilder {
	perms := &models.Permission{
		Read:  true,
		Write: true,
		Admin: false,
	}
	return &PermsBuilder{Perms: perms}
}

func (p *PermsBuilder) Build() *models.Permission {
	return p.Perms
}

func (p *PermsBuilder) WithoutPerm() *PermsBuilder {
	p.Perms = nil
	return p
}

func (p *PermsBuilder) WithPerm(r, w, a bool) *PermsBuilder {
	perm := &models.Permission{
		Read:  r,
		Write: w,
		Admin: a,
	}
	p.Perms = perm
	return p
}
