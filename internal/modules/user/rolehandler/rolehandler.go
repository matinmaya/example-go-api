package rolehandler

import (
	"reapp/internal/modules/user/rolemodel"
	"reapp/internal/modules/user/roleservice"
	"reapp/pkg/base/basehandler"
	"reapp/pkg/services/rediservice"
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
	basehandler.Update(ctx, h.service, &rolemodel.Role{}, func(role *rolemodel.Role, id uint64) error {
		role.UniqueScope = validators.ExceptByID(id)
		return nil
	}, nil, func(ctx *gin.Context, role *rolemodel.Role) error {
		role.PermissionIds = []uint32{}
		go rediservice.ClearCacheOfPerms()
		return nil
	})
}

func (h *RoleHandler) Delete(ctx *gin.Context) {
	basehandler.Delete(ctx, h.service, func(*gin.Context) error {
		go rediservice.ClearCacheOfPerms()
		return nil
	})
}
