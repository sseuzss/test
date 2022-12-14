package private

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"soldr/internal/app/api/models"
	srverrors "soldr/internal/app/api/server/errors"
	"soldr/internal/app/api/utils"
)

type users struct {
	Users []models.UserRoleTenant `json:"users"`
	Total uint64                  `json:"total"`
}

var usersSQLMappers = map[string]interface{}{
	"id":        "`{{table}}`.id",
	"hash":      "`{{table}}`.hash",
	"mail":      "`{{table}}`.mail",
	"name":      "`{{table}}`.name",
	"role_id":   "`{{table}}`.role_id",
	"status":    "`{{table}}`.status",
	"tenant_id": "`{{table}}`.tenant_id",
	"data": "CONCAT(`{{table}}`.hash, ' | ', " +
		"`{{table}}`.mail, ' | ', " +
		"`{{table}}`.name, ' | ', " +
		"`{{table}}`.status)",
}

// GetCurrentUser is a function to return account information
// @Summary Retrieve current user information
// @Tags Users
// @Produce json
// @Success 200 {object} utils.successResp{data=models.UserRoleTenant} "user info received successful"
// @Failure 403 {object} utils.errorResp "getting current user not permitted"
// @Failure 404 {object} utils.errorResp "current user not found"
// @Failure 500 {object} utils.errorResp "internal error on getting current user"
// @Router /user/ [get]
func GetCurrentUser(c *gin.Context) {
	var (
		err  error
		gDB  *gorm.DB
		resp models.UserRoleTenant
	)

	if gDB = utils.GetGormDB(c, "gDB"); gDB == nil {
		utils.HTTPError(c, srverrors.ErrInternalDBNotFound, nil)
		return
	}

	uid, _ := utils.GetUint64(c, "uid")

	err = gDB.Take(&resp.User, "id = ?", uid).
		Related(&resp.Role).
		Related(&resp.Tenant).Error
	if err != nil {
		utils.FromContext(c).WithError(err).Errorf("error finding current user")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HTTPError(c, srverrors.ErrUsersNotFound, err)
		} else {
			utils.HTTPError(c, srverrors.ErrInternal, err)
		}
		return
	} else if err = resp.Valid(); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error validating user data '%s'", resp.Hash)
		utils.HTTPError(c, srverrors.ErrUsersInvalidData, err)
		return
	}

	utils.HTTPSuccess(c, http.StatusOK, resp)
}

// ChangePasswordCurrentUser is a function to update account password
// @Summary Update password for current user (account)
// @Tags Users
// @Accept json
// @Produce json
// @Param json body models.Password true "container to validate and update account password"
// @Success 200 {object} utils.successResp "account password updated successful"
// @Failure 400 {object} utils.errorResp "invalid account password form data"
// @Failure 403 {object} utils.errorResp "updating account password not permitted"
// @Failure 404 {object} utils.errorResp "current user not found"
// @Failure 500 {object} utils.errorResp "internal error on updating account password"
// @Router /user/password [put]
func ChangePasswordCurrentUser(c *gin.Context) {
	var (
		encPass []byte
		err     error
		form    models.Password
		gDB     *gorm.DB
		user    models.UserPassword
	)

	if err = c.ShouldBindJSON(&form); err != nil || form.Valid() != nil {
		if err == nil {
			err = form.Valid()
		}
		utils.FromContext(c).WithError(err).Errorf("error binding JSON")
		utils.HTTPError(c, srverrors.ErrChangePasswordCurrentUserInvalidPassword, err)
		return
	}

	if gDB = utils.GetGormDB(c, "gDB"); gDB == nil {
		utils.HTTPError(c, srverrors.ErrInternalDBNotFound, err)
		return
	}

	uid, _ := utils.GetUint64(c, "uid")
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", uid)
	}

	if err = gDB.Scopes(scope).Take(&user).Error; err != nil {
		utils.FromContext(c).WithError(err).Errorf("error finding current user")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HTTPError(c, srverrors.ErrUsersNotFound, err)
		} else {
			utils.HTTPError(c, srverrors.ErrInternal, err)
		}
		return
	} else if err = user.Valid(); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error validating user data '%s'", user.Hash)
		utils.HTTPError(c, srverrors.ErrUsersInvalidData, err)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.CurrentPassword)); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error checking password for current user")
		utils.HTTPError(c, srverrors.ErrChangePasswordCurrentUserInvalidCurrentPassword, err)
		return
	}

	if encPass, err = utils.EncryptPassword(form.Password); err == nil {
		user.Password = string(encPass)
	} else {
		utils.FromContext(c).WithError(err).Errorf("error making new password for current user")
		utils.HTTPError(c, srverrors.ErrChangePasswordCurrentUserInvalidNewPassword, err)
		return
	}

	if err = gDB.Scopes(scope).Select("password").Save(&user).Error; err != nil {
		utils.FromContext(c).WithError(err).Errorf("error updating password for current user")
		utils.HTTPError(c, srverrors.ErrInternal, err)
		return
	}

	utils.HTTPSuccess(c, http.StatusOK, struct{}{})
}

// GetUsers returns users list
// @Summary Retrieve users list by filters
// @Tags Users
// @Produce json
// @Param request query utils.TableQuery true "query table params"
// @Success 200 {object} utils.successResp{data=users} "users list received successful"
// @Failure 400 {object} utils.errorResp "invalid query request data"
// @Failure 403 {object} utils.errorResp "getting users not permitted"
// @Failure 500 {object} utils.errorResp "internal error on getting users"
// @Router /users/ [get]
func GetUsers(c *gin.Context) {
	var (
		err    error
		gDB    *gorm.DB
		query  utils.TableQuery
		resp   users
		rids   []uint64
		roles  []models.Role
		tenans []models.Tenant
		tids   []uint64
	)

	if err = c.ShouldBindQuery(&query); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error binding query")
		utils.HTTPError(c, srverrors.ErrUsersInvalidRequest, err)
		return
	}

	if gDB = utils.GetGormDB(c, "gDB"); gDB == nil {
		utils.HTTPError(c, srverrors.ErrInternalDBNotFound, nil)
		return
	}

	query.Init("users", usersSQLMappers)

	rid, _ := utils.GetUint64(c, "rid")
	tid, _ := utils.GetUint64(c, "tid")
	uid, _ := utils.GetUint64(c, "uid")

	switch rid {
	case models.RoleSAdmin:
	case models.RoleAdmin:
		query.SetFilters([]func(db *gorm.DB) *gorm.DB{
			func(db *gorm.DB) *gorm.DB {
				return db.Where("tenant_id = ?", tid)
			},
		})
	case models.RoleUser:
		query.SetFilters([]func(db *gorm.DB) *gorm.DB{
			func(db *gorm.DB) *gorm.DB {
				return db.Where("id = ? AND tenant_id = ?", uid, tid)
			},
		})
	default:
		utils.FromContext(c).WithError(nil).Errorf("error filtering user role services: unexpected role")
		utils.HTTPError(c, srverrors.ErrInternal, err)
		return
	}

	if resp.Total, err = query.Query(gDB, &resp.Users); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error finding users")
		utils.HTTPError(c, srverrors.ErrInternal, err)
		return
	}

	for _, user := range resp.Users {
		rids = append(rids, user.RoleID)
		tids = append(tids, user.TenantID)
	}

	if err = gDB.Find(&roles, "id IN (?)", rids).Error; err != nil {
		utils.FromContext(c).WithError(err).Errorf("error finding linked roles")
		utils.HTTPError(c, srverrors.ErrInternal, err)
		return
	}

	if err = gDB.Find(&tenans, "id IN (?)", tids).Error; err != nil {
		utils.FromContext(c).WithError(err).Errorf("error finding linked tenants")
		utils.HTTPError(c, srverrors.ErrInternal, err)
		return
	}

	for i := range resp.Users {
		tenantID := resp.Users[i].TenantID
		for _, tenant := range tenans {
			if tenantID == tenant.ID {
				resp.Users[i].Tenant = tenant
				break
			}
		}

		roleID := resp.Users[i].RoleID
		for _, role := range roles {
			if roleID == role.ID {
				resp.Users[i].Role = role
				break
			}
		}
	}

	for i := 0; i < len(resp.Users); i++ {
		if err = resp.Users[i].Valid(); err != nil {
			utils.FromContext(c).WithError(err).Errorf("error validating user data '%s'", resp.Users[i].Hash)
			utils.HTTPError(c, srverrors.ErrUsersInvalidData, err)
			return
		}
	}

	utils.HTTPSuccess(c, http.StatusOK, resp)
}

// GetUser is a function to return user by hash
// @Summary Retrieve user by hash
// @Tags Users
// @Produce json
// @Param hash path string true "hash in hex format (md5)" minlength(32) maxlength(32)
// @Success 200 {object} utils.successResp{data=models.UserRoleTenant} "user received successful"
// @Failure 403 {object} utils.errorResp "getting user not permitted
// @Failure 404 {object} utils.errorResp "user not found"
// @Failure 500 {object} utils.errorResp "internal error on getting user"
// @Router /users/{hash} [get]
func GetUser(c *gin.Context) {
	var (
		err  error
		gDB  *gorm.DB
		hash string = c.Param("hash")
		resp models.UserRoleTenant
	)

	if gDB = utils.GetGormDB(c, "gDB"); gDB == nil {
		utils.HTTPError(c, srverrors.ErrInternalDBNotFound, nil)
		return
	}

	rid, _ := utils.GetUint64(c, "rid")
	tid, _ := utils.GetUint64(c, "tid")
	uid, _ := utils.GetUint64(c, "uid")
	scope := func(db *gorm.DB) *gorm.DB {
		switch rid {
		case models.RoleSAdmin:
			return db.Where("hash = ?", hash)
		case models.RoleAdmin:
			return db.Where("tenant_id = ? AND hash = ?", tid, hash)
		case models.RoleUser:
			return db.Where("tenant_id = ? AND hash = ? AND id = ?", tid, hash, uid)
		default:
			db.AddError(errors.New("unexpected user role"))
			return db
		}
	}

	if err = gDB.Scopes(scope).Take(&resp.User).Error; err != nil {
		utils.FromContext(c).WithError(err).Errorf("error finding user by hash")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HTTPError(c, srverrors.ErrUsersNotFound, err)
		} else {
			utils.HTTPError(c, srverrors.ErrInternal, err)
		}
		return
	}
	if err = gDB.Model(&resp.User).Related(&resp.Role).Related(&resp.Tenant).Error; err != nil {
		utils.FromContext(c).WithError(err).Errorf("error finding related models by user hash")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HTTPError(c, srverrors.ErrGetUserModelsNotFound, err)
		} else {
			utils.HTTPError(c, srverrors.ErrInternal, err)
		}
		return
	}
	if err = resp.Valid(); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error validating user data '%s'", resp.Hash)
		utils.HTTPError(c, srverrors.ErrUsersInvalidData, err)
		return
	}

	utils.HTTPSuccess(c, http.StatusOK, resp)
}

// CreateUser is a function to create new user
// @Summary Create new user
// @Tags Users
// @Accept json
// @Produce json
// @Param json body models.UserPassword true "user model to create from"
// @Success 201 {object} utils.successResp{data=models.UserRoleTenant} "user created successful"
// @Failure 400 {object} utils.errorResp "invalid user request data"
// @Failure 403 {object} utils.errorResp "creating user not permitted"
// @Failure 500 {object} utils.errorResp "internal error on creating user"
// @Router /users/ [post]
func CreateUser(c *gin.Context) {
	var (
		encPassword []byte
		err         error
		gDB         *gorm.DB
		resp        models.UserRoleTenant
		user        models.UserPassword
	)

	if err = c.ShouldBindJSON(&user); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error binding JSON")
		utils.HTTPError(c, srverrors.ErrUsersInvalidRequest, err)
		return
	}

	if gDB = utils.GetGormDB(c, "gDB"); gDB == nil {
		utils.HTTPError(c, srverrors.ErrInternalDBNotFound, nil)
		return
	}

	rid, _ := utils.GetUint64(c, "rid")
	tid, _ := utils.GetUint64(c, "tid")

	switch rid {
	case models.RoleSAdmin:
	case models.RoleAdmin:
		fallthrough
	case models.RoleUser:
		if user.RoleID < rid {
			user.RoleID = rid
		}
		user.TenantID = tid
	default:
		utils.FromContext(c).WithError(nil).Errorf("error filtering user role services: unexpected role")
		utils.HTTPError(c, srverrors.ErrInternal, nil)
		return
	}

	user.ID = 0
	user.Hash = utils.MakeUserHash(user.Name)
	if err = user.Valid(); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error validating user")
		utils.HTTPError(c, srverrors.ErrCreateUserInvalidUser, err)
		return
	}

	if encPassword, err = utils.EncryptPassword(user.Password); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error encoding password")
		utils.HTTPError(c, srverrors.ErrInternal, err)
		return
	} else {
		user.Password = string(encPassword)
	}

	if err = gDB.Create(&user).Error; err != nil {
		utils.FromContext(c).WithError(err).Errorf("error creating user")
		utils.HTTPError(c, srverrors.ErrInternal, err)
		return
	}

	err = gDB.Take(&resp.User, "hash = ?", user.Hash).
		Related(&resp.Role).
		Related(&resp.Tenant).Error
	if err != nil {
		utils.FromContext(c).WithError(err).Errorf("error finding user by hash")
		utils.HTTPError(c, srverrors.ErrInternal, err)
		return
	} else if err = resp.Valid(); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error validating user data '%s'", resp.Hash)
		utils.HTTPError(c, srverrors.ErrUsersInvalidData, err)
		return
	}

	utils.HTTPSuccess(c, http.StatusCreated, resp)
}

// PatchUser is a function to update user by hash
// @Summary Update user
// @Tags Users
// @Produce json
// @Param json body models.UserPassword true "user model to update"
// @Param hash path string true "user hash in hex format (md5)" minlength(32) maxlength(32)
// @Success 200 {object} utils.successResp{data=models.UserRoleTenant} "user updated successful"
// @Failure 400 {object} utils.errorResp "invalid user request data"
// @Failure 403 {object} utils.errorResp "updating user not permitted"
// @Failure 404 {object} utils.errorResp "user not found"
// @Failure 500 {object} utils.errorResp "internal error on updating user"
// @Router /users/{hash} [put]
func PatchUser(c *gin.Context) {
	var (
		err  error
		gDB  *gorm.DB
		hash = c.Param("hash")
		resp models.UserRoleTenant
		user models.UserPassword
	)

	if err = c.ShouldBindJSON(&user); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error binding JSON")
		utils.HTTPError(c, srverrors.ErrUsersInvalidRequest, err)
		return
	} else if hash != user.Hash {
		utils.FromContext(c).WithError(nil).Errorf("mismatch user hash to requested one")
		utils.HTTPError(c, srverrors.ErrUsersInvalidRequest, nil)
		return
	} else if err = user.User.Valid(); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error validating user JSON")
		utils.HTTPError(c, srverrors.ErrUsersInvalidRequest, err)
		return
	} else if err = user.Valid(); user.Password != "" && err != nil {
		utils.FromContext(c).WithError(err).Errorf("error validating user password")
		utils.HTTPError(c, srverrors.ErrUsersInvalidRequest, err)
		return
	}

	if gDB = utils.GetGormDB(c, "gDB"); gDB == nil {
		utils.HTTPError(c, srverrors.ErrInternal, nil)
		return
	}

	rid, _ := utils.GetUint64(c, "rid")
	tid, _ := utils.GetUint64(c, "tid")
	uid, _ := utils.GetUint64(c, "uid")
	scope := func(db *gorm.DB) *gorm.DB {
		switch rid {
		case models.RoleSAdmin:
			return db.Where("hash = ?", hash)
		case models.RoleAdmin:
			return db.Where("tenant_id = ? AND hash = ?", tid, hash)
		case models.RoleUser:
			return db.Where("tenant_id = ? AND hash = ? AND id = ?", tid, hash, uid)
		default:
			db.AddError(errors.New("unexpected user role"))
			return db
		}
	}

	public_info := []interface{}{"mail", "name", "status"}
	if user.Password != "" {
		if encPassword, err := utils.EncryptPassword(user.Password); err != nil {
			utils.FromContext(c).WithError(err).Errorf("error encoding password")
			utils.HTTPError(c, srverrors.ErrInternal, err)
			return
		} else {
			user.Password = string(encPassword)
			public_info = append(public_info, "password")
		}
		err = gDB.Scopes(scope).Select("", public_info...).Save(&user).Error
	} else {
		err = gDB.Scopes(scope).Select("", public_info...).Save(&user.User).Error
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		utils.FromContext(c).WithError(nil).Errorf("error updating user by hash '%s', user not found", hash)
		utils.HTTPError(c, srverrors.ErrUsersNotFound, err)
		return
	} else if err != nil {
		utils.FromContext(c).WithError(err).Errorf("error updating user by hash '%s'", hash)
		utils.HTTPError(c, srverrors.ErrInternal, err)
		return
	}

	if err = gDB.Scopes(scope).Take(&resp.User).Error; err != nil {
		utils.FromContext(c).WithError(err).Errorf("error finding user by hash")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HTTPError(c, srverrors.ErrUsersNotFound, err)
		} else {
			utils.HTTPError(c, srverrors.ErrInternal, err)
		}
		return
	}
	if err = gDB.Model(&resp.User).Related(&resp.Role).Related(&resp.Tenant).Error; err != nil {
		utils.FromContext(c).WithError(err).Errorf("error finding related models by user hash")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HTTPError(c, srverrors.ErrPatchUserModelsNotFound, err)
		} else {
			utils.HTTPError(c, srverrors.ErrInternal, err)
		}
		return
	}
	if err = resp.Valid(); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error validating user data '%s'", resp.Hash)
		utils.HTTPError(c, srverrors.ErrInternal, err)
		return
	}

	utils.HTTPSuccess(c, http.StatusOK, resp)
}

// DeleteUser is a function to delete user by hash
// @Summary Delete user by hash
// @Tags Users
// @Produce json
// @Param hash path string true "hash in hex format (md5)" minlength(32) maxlength(32)
// @Success 200 {object} utils.successResp "user deleted successful"
// @Failure 403 {object} utils.errorResp "deleting user not permitted"
// @Failure 404 {object} utils.errorResp "user not found"
// @Failure 500 {object} utils.errorResp "internal error on deleting user"
// @Router /users/{hash} [delete]
func DeleteUser(c *gin.Context) {
	var (
		err  error
		gDB  *gorm.DB
		hash string = c.Param("hash")
		user models.UserRoleTenant
	)

	if gDB = utils.GetGormDB(c, "gDB"); gDB == nil {
		utils.HTTPError(c, srverrors.ErrInternalDBNotFound, nil)
		return
	}

	rid, _ := utils.GetUint64(c, "rid")
	tid, _ := utils.GetUint64(c, "tid")
	uid, _ := utils.GetUint64(c, "uid")
	scope := func(db *gorm.DB) *gorm.DB {
		switch rid {
		case models.RoleSAdmin:
			return db.Where("hash = ?", hash)
		case models.RoleAdmin:
			return db.Where("tenant_id = ? AND hash = ?", tid, hash)
		case models.RoleUser:
			return db.Where("tenant_id = ? AND hash = ? AND id = ?", tid, hash, uid)
		default:
			db.AddError(errors.New("unexpected user role"))
			return db
		}
	}

	if err = gDB.Scopes(scope).Take(&user.User).Error; err != nil {
		utils.FromContext(c).WithError(err).Errorf("error finding user by hash")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HTTPError(c, srverrors.ErrUsersNotFound, err)
		} else {
			utils.HTTPError(c, srverrors.ErrInternal, err)
		}
		return
	}
	if err = gDB.Model(&user.User).Related(&user.Role).Related(&user.Tenant).Error; err != nil {
		utils.FromContext(c).WithError(err).Errorf("error finding related models by user hash")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HTTPError(c, srverrors.ErrDeleteUserModelsNotFound, err)
		} else {
			utils.HTTPError(c, srverrors.ErrInternal, err)
		}
		return
	}
	if err = user.Valid(); err != nil {
		utils.FromContext(c).WithError(err).Errorf("error validating user data '%s'", user.Hash)
		utils.HTTPError(c, srverrors.ErrUsersInvalidData, err)
		return
	}

	if err = gDB.Delete(&user.User).Error; err != nil {
		utils.FromContext(c).WithError(err).Errorf("error deleting user by hash '%s'", hash)
		utils.HTTPError(c, srverrors.ErrInternal, err)
		return
	}

	utils.HTTPSuccess(c, http.StatusOK, struct{}{})
}
