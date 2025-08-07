package roleservice

import (
	"fmt"
	"log"
	"reapp/internal/modules/user/rolemodel"
	"reapp/internal/modules/user/rolerepository"
	"reapp/pkg/filterscopes"
	"reapp/pkg/lang"
	"reapp/pkg/paginator"

	"gorm.io/gorm"
)

type IRoleService interface {
	Create(db *gorm.DB, role *rolemodel.Role) error
	Update(db *gorm.DB, role *rolemodel.Role) error
	GetByID(db *gorm.DB, id uint64) (*rolemodel.Role, error)
	GetDetail(db *gorm.DB, id uint64) (*rolemodel.Role, error)
	Delete(db *gorm.DB, id uint64) error
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
	tx := db.Begin()
	if err := s.repository.Create(tx, role); err != nil {
		tx.Rollback()
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	if err := s.repository.AddPermissions(tx, role.ID, role.PermissionIds); err != nil {
		tx.Rollback()
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	return nil
}

func (s *RoleService) Update(db *gorm.DB, role *rolemodel.Role) error {
	if _, err := s.repository.GetByID(db, role.ID); err != nil {
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(db, "response", "error"))
	}

	tx := db.Begin()
	if err := s.repository.Update(tx, role); err != nil {
		tx.Rollback()
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	if err := s.repository.RemovePermissions(tx, role.ID); err != nil {
		tx.Rollback()
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	if err := s.repository.AddPermissions(tx, role.ID, role.PermissionIds); err != nil {
		tx.Rollback()
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	return nil
}

func (s *RoleService) GetByID(db *gorm.DB, id uint64) (*rolemodel.Role, error) {
	return s.repository.GetByID(db, uint16(id))
}

func (s *RoleService) GetDetail(db *gorm.DB, id uint64) (*rolemodel.Role, error) {
	return s.repository.GetDetail(db, uint16(id))
}

func (s *RoleService) Delete(db *gorm.DB, id uint64) error {
	if id < 2 {
		return fmt.Errorf("%s", lang.TranByDB(db, "response", "error"))
	}

	if _, err := s.repository.GetByID(db, uint16(id)); err != nil {
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(db, "response", "not_found"))
	}

	count, err := s.repository.RoleUserCount(db, uint16(id))
	if err != nil {
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(db, "response", "error"))
	}
	if count > 0 {
		return fmt.Errorf("%s", lang.TranByDB(db, "response", "error"))
	}

	return s.repository.Delete(db, uint16(id))
}

func (s *RoleService) GetAll(db *gorm.DB) ([]rolemodel.Role, error) {
	return s.repository.GetAll(db)
}

func (s *RoleService) List(db *gorm.DB, pg *paginator.Pagination, filters []filterscopes.QueryFilter) error {
	return s.repository.List(db, pg, filters)
}
