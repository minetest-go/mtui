package modmanager

import (
	"mtui/db"
	"mtui/types"
)

type ModManager struct {
	world_dir      string
	repo           *db.ModRepository
	handlercontext *HandlerContext
	handlers       map[types.SourceType]SourceTypeHandler
}

func New(world_dir string, repo *db.ModRepository) *ModManager {
	return &ModManager{
		world_dir: world_dir,
		repo:      repo,
		handlercontext: &HandlerContext{
			WorldDir: world_dir,
			Repo:     repo,
		},
		handlers: map[types.SourceType]SourceTypeHandler{
			types.SourceTypeGIT: &GitModHandler{},
			types.SourceTypeCDB: &ContentDBModHandler{},
		},
	}
}

func (m *ModManager) Mod(id string) (*types.Mod, error) {
	return m.repo.GetByID(id)
}

func (m *ModManager) Create(mod *types.Mod) error {
	handler := m.handlers[mod.SourceType]
	return handler.Create(m.handlercontext, mod)
}

func (m *ModManager) Status(mod *types.Mod) (*ModStatus, error) {
	handler := m.handlers[mod.SourceType]
	return handler.Status(m.handlercontext, mod)
}

func (m *ModManager) Update(mod *types.Mod, version string) error {
	handler := m.handlers[mod.SourceType]
	return handler.Update(m.handlercontext, mod, version)
}

func (m *ModManager) Remove(mod *types.Mod) error {
	handler := m.handlers[mod.SourceType]
	return handler.Remove(m.handlercontext, mod)
}
