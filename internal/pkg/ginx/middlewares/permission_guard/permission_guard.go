package permission_guard

import (
	"net/http"

	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/service/staff"
	"github.com/emirpasic/gods/v2/sets"
	"github.com/emirpasic/gods/v2/sets/hashset"
	"github.com/gin-gonic/gin"

	"bitbucket.org/starlinglabs/cst-wstyle-integration/pkg/ginx"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/pkg/ginx/response"
)

type StaffPermissionMiddlewareBuilder struct {
	publicPaths sets.Set[string]
	staffSvc    staff.StaffService
}

func NewStaffPermissionHandler(staffSvc staff.StaffService) *StaffPermissionMiddlewareBuilder {
	return &StaffPermissionMiddlewareBuilder{
		publicPaths: hashset.New[string](),
		staffSvc:    staffSvc,
	}
}

func (l *StaffPermissionMiddlewareBuilder) IgnorePaths(path string) *StaffPermissionMiddlewareBuilder {
	l.publicPaths.Add(path)
	return l
}

func (m *StaffPermissionMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if m.publicPaths.Contains(ctx.Request.URL.Path) {
			return
		}
		rawVal, ok := ctx.Get("shopline")
		if !ok {
			res := response.NewResponse(403000, http.StatusForbidden, "Not Authorized.")
			response.SendResponse(ctx, res, nil)
			return
		}
		claims, ok := rawVal.(ginx.ShoplineClaims)
		if !ok {
			res := response.NewResponse(403000, http.StatusForbidden, "Not Authorized.")
			response.SendResponse(ctx, res, nil)
			return
		}

		isAllowed, err := m.staffSvc.CheckPermission(ctx, claims.Data.StaffId, claims.Data.MerchantId, ctx.Request.URL.Path)
		if !isAllowed || err != nil {
			res := response.NewResponse(403000, http.StatusForbidden, "Not Authorized.")
			response.SendResponse(ctx, res, nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
