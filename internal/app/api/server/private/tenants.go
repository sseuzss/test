package private

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"soldr/internal/app/api/models"
	srverrors "soldr/internal/app/api/server/errors"
	"soldr/internal/app/api/utils"
)

type tenants struct {
	Tenants []models.Tenant `json:"tenants"`
	Total   uint64          `json:"total"`
}

var tenantsSQLMappers = map[string]interface{}{
	"id":     "`{{table}}`.id",
	"hash":   "`{{table}}`.hash",
	"status": "`{{table}}`.status",
	"data": "CONCAT(`{{table}}`.hash, ' | ', " +
		"`{{table}}`.status)",
}

// GetTenants is a function to return tenants list view on dashboard
// @Summary Retrieve tenants list
// @Tags Tenants
// @Produce json
// @Param request query utils.TableQuery true "query table params"
// @Success 200 {object} utils.successResp{data=tenants} "tenants list received successful"
// @Failure 400 {object} utils.errorResp "invalid query request data"
// @Failure 403 {object} utils.errorResp "getting tenants not permitted"
// @Failure 500 {object} utils.errorResp "internal error on getting tenants"
// @Router /tenants/ [get]
func GetTenants(c *gin.Context) {
	var (
		err   error
		gDB   *gorm.DB
		query utils.TableQuery
		resp  tenants
	)

	if err = c.ShouldBindQuery(&query); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error binding query")
		utils.HTTPError(c, srverrors.ErrTenantsInvalidRequest, err)
		return
	}

	if gDB = utils.GetGormDB(c, "gDB"); gDB == nil {
		utils.HTTPError(c, srverrors.ErrInternalDBNotFound, nil)
		return
	}

	query.Init("tenants", tenantsSQLMappers)

	rid, _ := utils.GetUint64(c, "rid")
	tid, _ := utils.GetUint64(c, "tid")

	switch rid {
	case models.RoleSAdmin:
	case models.RoleAdmin:
		fallthrough
	case models.RoleUser:
		query.SetFilters([]func(db *gorm.DB) *gorm.DB{
			func(db *gorm.DB) *gorm.DB {
				return db.Where("id = ?", tid)
			},
		})
	default:
		utils.FromContext(c).WithError(nil).Errorf("error filtering user role services: unexpected role")
		utils.HTTPError(c, srverrors.ErrInternal, nil)
		return
	}

	if resp.Total, err = query.Query(gDB, &resp.Tenants); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error finding tenants")
		utils.HTTPError(c, srverrors.ErrInternal, err)
		return
	}

	for i := 0; i < len(resp.Tenants); i++ {
		if err = resp.Tenants[i].Valid(); err != nil {
			utils.FromContext(c).WithError(err).Errorf("error validating tenant data '%s'", resp.Tenants[i].Hash)
			utils.HTTPError(c, srverrors.ErrTenantsInvalidData, err)
			return
		}
	}

	utils.HTTPSuccess(c, http.StatusOK, resp)
}

// GetTenant returns tenant by hash
// @Summary Retrieve tenant by hash
// @Tags Tenants
// @Produce json
// @Param hash path string true "hash in hex format (md5)" minlength(32) maxlength(32)
// @Success 200 {object} utils.successResp{data=models.Tenant} "tenant received successful"
// @Failure 403 {object} utils.errorResp "getting tenant not permitted
// @Failure 404 {object} utils.errorResp "tenant not found"
// @Failure 500 {object} utils.errorResp "internal error on getting tenant"
// @Router /tenants/{hash} [get]
func GetTenant(c *gin.Context) {
	var (
		err  error
		gDB  *gorm.DB
		hash string = c.Param("hash")
		resp models.Tenant
	)

	if gDB = utils.GetGormDB(c, "gDB"); gDB == nil {
		utils.HTTPError(c, srverrors.ErrInternalDBNotFound, err)
		return
	}

	rid, _ := utils.GetUint64(c, "rid")
	tid, _ := utils.GetUint64(c, "tid")
	scope := func(db *gorm.DB) *gorm.DB {
		switch rid {
		case models.RoleSAdmin:
			return db.Where("hash = ?", hash)
		case models.RoleAdmin:
			fallthrough
		case models.RoleUser:
			return db.Where("hash = ? AND id = ?", hash, tid)
		default:
			db.AddError(errors.New("unexpected user role"))
			return db
		}
	}

	if err = gDB.Scopes(scope).Take(&resp).Error; err != nil {
		utils.FromContext(c).WithError(err).Errorf("error finding tenant by hash")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HTTPError(c, srverrors.ErrTenantsNotFound, err)
		} else {
			utils.HTTPError(c, srverrors.ErrInternal, err)
		}
		return
	} else if err = resp.Valid(); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error validating tenant data '%s'", resp.Hash)
		utils.HTTPError(c, srverrors.ErrTenantsInvalidData, err)
		return
	}

	utils.HTTPSuccess(c, http.StatusOK, resp)
}

// CreateTenant is a function to create new tenant
// @Summary Create new tenant
// @Tags Tenants
// @Accept json
// @Produce json
// @Param json body models.Tenant true "tenant model to create from"
// @Success 201 {object} utils.successResp{data=models.Tenant} "tenant created successful"
// @Failure 400 {object} utils.errorResp "invalid tenant request data"
// @Failure 403 {object} utils.errorResp "creating tenant not permitted"
// @Failure 500 {object} utils.errorResp "internal error on creating tenant"
// @Router /tenants/ [post]
func CreateTenant(c *gin.Context) {
	var (
		err    error
		gDB    *gorm.DB
		tenant models.Tenant
	)

	if err = c.ShouldBindJSON(&tenant); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error binding JSON")
		utils.HTTPError(c, srverrors.ErrTenantsInvalidRequest, err)
		return
	}

	if gDB = utils.GetGormDB(c, "gDB"); gDB == nil {
		utils.HTTPError(c, srverrors.ErrInternalDBNotFound, nil)
		return
	}

	tenant.ID = 0
	tenant.Hash = utils.MakeTenantHash(tenant.Description)

	if err = gDB.Create(&tenant).Error; err != nil {
		utils.FromContext(c).WithError(err).Errorf("error creating tenant")
		utils.HTTPError(c, srverrors.ErrInternal, err)
		return
	}

	utils.HTTPSuccess(c, http.StatusCreated, tenant)
}

// PatchTenant updates tenant by hash
// @Summary Update tenant
// @Tags Tenants
// @Produce json
// @Param json body models.Tenant true "tenant model to update"
// @Param hash path string true "tenant hash in hex format (md5)" minlength(32) maxlength(32)
// @Success 200 {object} utils.successResp{data=models.Tenant} "tenant updated successful"
// @Failure 400 {object} utils.errorResp "invalid tenant request data"
// @Failure 403 {object} utils.errorResp "updating tenant not permitted"
// @Failure 404 {object} utils.errorResp "tenant not found"
// @Failure 500 {object} utils.errorResp "internal error on updating tenant"
// @Router /tenant/{hash} [put]
func PatchTenant(c *gin.Context) {
	var (
		err    error
		gDB    *gorm.DB
		hash   string = c.Param("hash")
		tenant models.Tenant
	)

	if err = c.ShouldBindJSON(&tenant); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error binding JSON")
		utils.HTTPError(c, srverrors.ErrTenantsInvalidRequest, err)
		return
	} else if hash != tenant.Hash {
		utils.FromContext(c).WithError(nil).Errorf("mismatch tenant hash to requested one")
		utils.HTTPError(c, srverrors.ErrTenantsInvalidRequest, nil)
		return
	} else if err = tenant.Valid(); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error validating tenant JSON")
		utils.HTTPError(c, srverrors.ErrTenantsInvalidRequest, err)
		return
	}

	if gDB = utils.GetGormDB(c, "gDB"); gDB == nil {
		utils.HTTPError(c, srverrors.ErrInternalDBNotFound, err)
		return
	}

	rid, _ := utils.GetUint64(c, "rid")
	tid, _ := utils.GetUint64(c, "tid")
	scope := func(db *gorm.DB) *gorm.DB {
		switch rid {
		case models.RoleSAdmin:
			return db.Where("hash = ?", hash)
		case models.RoleAdmin:
			fallthrough
		case models.RoleUser:
			return db.Where("hash = ? AND id = ?", hash, tid)
		default:
			db.AddError(errors.New("unexpected user role"))
			return db
		}
	}

	public_info := []interface{}{"description", "status"}
	err = gDB.Scopes(scope).Select("", public_info...).Save(&tenant).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		utils.FromContext(c).WithError(nil).Errorf("error updating tenant by hash '%s', tenant not found", hash)
		utils.HTTPError(c, srverrors.ErrTenantsNotFound, err)
		return
	} else if err != nil {
		utils.FromContext(c).WithError(err).Errorf("error updating tenant by hash '%s'", hash)
		utils.HTTPError(c, srverrors.ErrInternal, err)
		return
	}

	utils.HTTPSuccess(c, http.StatusOK, tenant)
}

// DeleteTenant is a function to delete tenant by hash
// @Summary Delete tenant by hash
// @Tags Tenants
// @Produce json
// @Param hash path string true "hash in hex format (md5)" minlength(32) maxlength(32)
// @Success 200 {object} utils.successResp "tenant deleted successful"
// @Failure 403 {object} utils.errorResp "deleting tenant not permitted"
// @Failure 404 {object} utils.errorResp "tenant not found"
// @Failure 500 {object} utils.errorResp "internal error on deleting tenant"
// @Router /tenants/{hash} [delete]
func DeleteTenant(c *gin.Context) {
	var (
		err    error
		gDB    *gorm.DB
		hash   string = c.Param("hash")
		tenant models.Tenant
	)

	if gDB = utils.GetGormDB(c, "gDB"); gDB == nil {
		utils.HTTPError(c, srverrors.ErrInternalDBNotFound, err)
		return
	}

	rid, _ := utils.GetUint64(c, "rid")
	tid, _ := utils.GetUint64(c, "tid")
	scope := func(db *gorm.DB) *gorm.DB {
		switch rid {
		case models.RoleSAdmin:
			return db.Where("hash = ?", hash)
		case models.RoleAdmin:
			fallthrough
		case models.RoleUser:
			return db.Where("hash = ? AND id = ?", hash, tid)
		default:
			db.AddError(errors.New("unexpected user role"))
			return db
		}
	}

	if err = gDB.Scopes(scope).Take(&tenant).Error; err != nil {
		utils.FromContext(c).WithError(err).Errorf("error finding tenant by hash")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HTTPError(c, srverrors.ErrTenantsNotFound, err)
		} else {
			utils.HTTPError(c, srverrors.ErrInternal, err)
		}
		return
	} else if err = tenant.Valid(); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error validating tenant data '%s'", tenant.Hash)
		utils.HTTPError(c, srverrors.ErrTenantsInvalidData, err)
		return
	}

	if err = gDB.Delete(&tenant).Error; err != nil {
		utils.FromContext(c).WithError(err).Errorf("error deleting tenant by hash '%s'", hash)
		utils.HTTPError(c, srverrors.ErrInternal, err)
		return
	}

	utils.HTTPSuccess(c, http.StatusOK, struct{}{})
}
