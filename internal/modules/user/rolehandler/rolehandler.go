package rolehandler

import (
	"reapp/internal/modules/user/rolemodel"
	"reapp/internal/modules/user/roleservice"
	"reapp/pkg/base/basehandler"
	"reapp/pkg/validators"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	service roleservice.IRoleService
}

func NewRoleHandler(s roleservice.IRoleService) *RoleHandler {
	return &RoleHandler{service: s}
}

func (h *RoleHandler) List(ctx *gin.Context) {
	basehandler.Paginate(ctx, h.service, &rolemodel.RoleListQuery{})
}

func (h *RoleHandler) GetAll(ctx *gin.Context) {
	basehandler.GetAll(ctx, h.service)
}

func (h *RoleHandler) GetDetail(ctx *gin.Context) {
	basehandler.GetDetail(ctx, h.service)
}

func (h *RoleHandler) Create(ctx *gin.Context) {
	basehandler.Create(ctx, h.service, &rolemodel.Role{}, &rolemodel.Role{}, nil, nil)
}

func (h *RoleHandler) Update(ctx *gin.Context) {
	basehandler.Update(ctx, h.service, &rolemodel.Role{}, func(modelDTO any, id uint64) error {
		if dto, ok := modelDTO.(*rolemodel.Role); ok {
			dto.ScopeUnique = validators.ExceptByID(id)
		}
		return nil
	}, nil, func(model *rolemodel.Role) error {
		model.PermissionIds = []uint32{}
		return nil
	})
}

func (h *RoleHandler) Delete(ctx *gin.Context) {
	basehandler.Delete(ctx, h.service)
}
