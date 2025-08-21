package module

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"reapp/pkg/lang"
	"reapp/pkg/paginator"
	"reapp/pkg/queryfilter"
)

type ServiceYields[T IWithID[TID], TID IUintID] struct {
	CreateUseCase TBusinessLogic[T]
	UpdateUseCase TBusinessLogic[T]
	DeleteUseCase TBusinessLogic[T]
}

type Service[T IWithID[TID], TID IUintID] struct {
	repository *Repository[T, TID]
	yields     ServiceYields[T, TID]
}

func NewService[T IWithID[TID], TID IUintID](
	repository *Repository[T, TID],
	yields ServiceYields[T, TID],
) IService[T, TID] {
	return Service[T, TID]{repository, yields}
}

func (s Service[T, TID]) Create(db *gorm.DB, model *T) error {
	tx := db.Begin()
	if err := s.repository.Create(tx, model); err != nil {
		tx.Rollback()
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	if s.yields.CreateUseCase != nil {
		if err := s.yields.CreateUseCase(tx, model); err != nil {
			tx.Rollback()
			log.Printf("%s", err.Error())
			return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	return nil
}

func (s Service[T, TID]) GetByID(db *gorm.DB, id TID) (*T, error) {
	return s.repository.GetByID(db, id)
}

func (s Service[T, TID]) Update(db *gorm.DB, model *T) error {
	if _, err := s.repository.GetByID(db, (*model).GetID()); err != nil {
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(db, "response", "error"))
	}

	tx := db.Begin()
	if err := s.repository.Update(tx, model); err != nil {
		tx.Rollback()
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	if s.yields.UpdateUseCase != nil {
		if err := s.yields.UpdateUseCase(tx, model); err != nil {
			tx.Rollback()
			log.Printf("%s", err.Error())
			return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	return nil
}

func (s Service[T, TID]) Delete(db *gorm.DB, id TID) error {
	model, err := s.repository.GetByID(db, id)
	if err != nil {
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(db, "response", "error"))
	}

	if s.yields.DeleteUseCase != nil {
		if err := s.yields.DeleteUseCase(db, model); err != nil {
			log.Printf("%s", err.Error())
			return fmt.Errorf("%s", lang.TranByDB(db, "response", "error"))
		}
	}

	return s.repository.Delete(db, id)
}

func (s Service[T, TID]) List(ctx *gin.Context, db *gorm.DB, pg *paginator.Pagination[T], filterFields []queryfilter.FilterField) error {
	return s.repository.List(ctx, db, pg, filterFields)
}
