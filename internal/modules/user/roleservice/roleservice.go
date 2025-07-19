package roleservice

import (
	"fmt"
	"reapp/internal/modules/user/rolemodel"
	"reapp/internal/modules/user/rolerepository"
	"reapp/pkg/filterscopes"
	"reapp/pkg/paginator"

	"gorm.io/gorm"
)

type IRoleService interface {
	Create(db *gorm.DB, role *rolemodel.Role) error
	Update(db *gorm.DB, role *rolemodel.Role) error
	GetByID(db *gorm.DB, id uint16) (*rolemodel.Role, error)
	GetDetail(db *gorm.DB, id uint16) (*rolemodel.Role, error)
	Delete(db *gorm.DB, id uint16) error
	List(db *gorm.DB, pg *paginator.Pagination, filters []filterscopes.QueryFilter) error
	GetAll(db *gorm.DB) ([]rolemodel.Role, error)
}

type RoleService struct {
	repository *rolerepository.RoleRepository
}

func NewRoleService(r *rolerepository.RoleRepository) IRoleService {
	return &RoleService{repository: r}
}

func (s *RoleService) Create(db *gorm.DB, role *rolemodel.Role) error {
	return s.repository.Create(db, role)
}

func (s *RoleService) Update(db *gorm.DB, role *rolemodel.Role) error {
	if _, err := s.repository.GetByID(db, role.ID); err != nil {
		return fmt.Errorf("something went wrong")
	}

	return s.repository.Update(db, role)
}

func (s *RoleService) GetByID(db *gorm.DB, id uint16) (*rolemodel.Role, error) {
	return s.repository.GetByID(db, id)
}

func (s *RoleService) GetDetail(db *gorm.DB, id uint16) (*rolemodel.Role, error) {
	return s.repository.GetDetail(db, id)
}

func (s *RoleService) Delete(db *gorm.DB, id uint16) error {
	if _, err := s.repository.GetByID(db, id); err != nil {
		return fmt.Errorf("role not found")
	}

	count, err := s.repository.RoleUserCount(db, id)
	if err != nil {
		return fmt.Errorf("failed to count role users: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("can't delete role: it is still assigned to %d user(s)", count)
	}

	return s.repository.Delete(db, id)
}

func (s *RoleService) GetAll(db *gorm.DB) ([]rolemodel.Role, error) {
	return s.repository.GetAll(db)
}

func (s *RoleService) List(db *gorm.DB, pg *paginator.Pagination, filters []filterscopes.QueryFilter) error {
	return s.repository.List(db, pg, filters)
}
