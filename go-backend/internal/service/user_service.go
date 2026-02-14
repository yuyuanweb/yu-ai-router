package service

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/yupi/airouter/go-backend/internal/common"
	"github.com/yupi/airouter/go-backend/internal/constant"
	"github.com/yupi/airouter/go-backend/internal/errno"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"github.com/yupi/airouter/go-backend/internal/model/vo"
	"github.com/yupi/airouter/go-backend/internal/repository"
)

const (
	minAccountLength  = 4
	minPasswordLength = 8
	unlimitedQuota    = int64(-1)
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) UserRegister(userAccount, userPassword, checkPassword string) (int64, error) {
	if hasBlank(userAccount, userPassword, checkPassword) {
		return 0, errno.NewWithMessage(errno.ParamsError, "参数为空")
	}
	if len(userAccount) < minAccountLength {
		return 0, errno.NewWithMessage(errno.ParamsError, "账号长度过短")
	}
	if len(userPassword) < minPasswordLength || len(checkPassword) < minPasswordLength {
		return 0, errno.NewWithMessage(errno.ParamsError, "密码长度过短")
	}
	if userPassword != checkPassword {
		return 0, errno.NewWithMessage(errno.ParamsError, "两次输入的密码不一致")
	}

	count, err := s.userRepo.CountByAccount(userAccount)
	if err != nil {
		return 0, errno.New(errno.SystemError)
	}
	if count > 0 {
		return 0, errno.NewWithMessage(errno.ParamsError, "账号重复")
	}

	user := &entity.User{
		UserAccount:  userAccount,
		UserPassword: s.GetEncryptPassword(userPassword),
		UserName:     constant.DefaultUserName,
		UserRole:     constant.DefaultRole,
	}
	id, err := s.userRepo.Create(user)
	if err != nil {
		return 0, errno.NewWithMessage(errno.OperationError, "注册失败，数据库错误")
	}
	return id, nil
}

func (s *UserService) UserLogin(userAccount, userPassword string, c *gin.Context) (*vo.LoginUserVO, error) {
	if hasBlank(userAccount, userPassword) {
		return nil, errno.NewWithMessage(errno.ParamsError, "参数为空")
	}
	if len(userAccount) < minAccountLength {
		return nil, errno.NewWithMessage(errno.ParamsError, "账号长度过短")
	}
	if len(userPassword) < minPasswordLength {
		return nil, errno.NewWithMessage(errno.ParamsError, "密码长度过短")
	}

	encryptPassword := s.GetEncryptPassword(userPassword)
	user, err := s.userRepo.GetByAccountAndPassword(userAccount, encryptPassword)
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	if user == nil {
		return nil, errno.NewWithMessage(errno.ParamsError, "用户不存在或密码错误")
	}

	session := sessions.Default(c)
	session.Set(constant.UserLoginState, strconv.FormatInt(user.ID, 10))
	if err = session.Save(); err != nil {
		return nil, errno.New(errno.SystemError)
	}
	loginUserVO := s.GetLoginUserVO(user)
	return &loginUserVO, nil
}

func (s *UserService) GetLoginUser(c *gin.Context) (*entity.User, error) {
	session := sessions.Default(c)
	sessionValue := session.Get(constant.UserLoginState)
	if sessionValue == nil {
		return nil, errno.New(errno.NotLoginError)
	}

	userIDStr, ok := sessionValue.(string)
	if !ok || userIDStr == "" {
		return nil, errno.New(errno.NotLoginError)
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		return nil, errno.New(errno.NotLoginError)
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	if user == nil {
		return nil, errno.New(errno.NotLoginError)
	}
	return user, nil
}

func (s *UserService) UserLogout(c *gin.Context) error {
	session := sessions.Default(c)
	sessionValue := session.Get(constant.UserLoginState)
	if sessionValue == nil {
		return errno.NewWithMessage(errno.OperationError, "用户未登录")
	}
	session.Delete(constant.UserLoginState)
	if err := session.Save(); err != nil {
		return errno.New(errno.SystemError)
	}
	return nil
}

func (s *UserService) CreateUser(req dto.UserAddRequest) (int64, error) {
	user := &entity.User{
		UserName:     req.UserName,
		UserAccount:  req.UserAccount,
		UserAvatar:   req.UserAvatar,
		UserProfile:  req.UserProfile,
		UserRole:     req.UserRole,
		UserPassword: s.GetEncryptPassword(constant.DefaultUserPassword),
	}
	id, err := s.userRepo.Create(user)
	if err != nil {
		return 0, errno.New(errno.OperationError)
	}
	return id, nil
}

func (s *UserService) GetUserByID(id int64) (*entity.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errno.New(errno.SystemError)
	}
	if user == nil {
		return nil, errno.New(errno.NotFoundError)
	}
	return user, nil
}

func (s *UserService) DeleteUser(id int64) (bool, error) {
	ok, err := s.userRepo.SoftDeleteByID(id)
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	return ok, nil
}

func (s *UserService) UpdateUser(req dto.UserUpdateRequest) (bool, error) {
	if req.ID == nil {
		return false, errno.New(errno.ParamsError)
	}
	user := &entity.User{
		ID:          req.ID.Int64(),
		UserName:    req.UserName,
		UserAvatar:  req.UserAvatar,
		UserProfile: req.UserProfile,
		UserRole:    req.UserRole,
	}
	ok, err := s.userRepo.UpdateByID(user)
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	if !ok {
		return false, errno.New(errno.OperationError)
	}
	return true, nil
}

func (s *UserService) ListUserVOByPage(req dto.UserQueryRequest) (common.PageResponse[vo.UserVO], error) {
	pageNum := req.PageNum
	if pageNum <= 0 {
		pageNum = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	sortOrder := req.SortOrder
	if sortOrder == "" {
		sortOrder = "descend"
	}
	req.PageNum = pageNum
	req.PageSize = pageSize
	req.SortOrder = sortOrder

	users, total, err := s.userRepo.ListByPage(req)
	if err != nil {
		return common.PageResponse[vo.UserVO]{}, errno.New(errno.SystemError)
	}
	userVOList := s.GetUserVOList(users)
	return common.BuildPageResponse(userVOList, pageNum, pageSize, total), nil
}

func (s *UserService) GetLoginUserVO(user *entity.User) vo.LoginUserVO {
	return vo.LoginUserVO{
		ID:          user.ID,
		UserAccount: user.UserAccount,
		UserName:    user.UserName,
		UserAvatar:  user.UserAvatar,
		UserProfile: user.UserProfile,
		UserRole:    user.UserRole,
		UserStatus:  user.UserStatus,
		TokenQuota:  user.TokenQuota,
		UsedTokens:  user.UsedTokens,
		CreateTime:  user.CreateTime,
		UpdateTime:  user.UpdateTime,
	}
}

func (s *UserService) GetUserVO(user *entity.User) vo.UserVO {
	return vo.UserVO{
		ID:          user.ID,
		UserAccount: user.UserAccount,
		UserName:    user.UserName,
		UserAvatar:  user.UserAvatar,
		UserProfile: user.UserProfile,
		UserRole:    user.UserRole,
		UserStatus:  user.UserStatus,
		TokenQuota:  user.TokenQuota,
		UsedTokens:  user.UsedTokens,
		CreateTime:  user.CreateTime,
	}
}

func (s *UserService) GetRemainingQuota(userID int64) (int64, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return 0, errno.New(errno.SystemError)
	}
	if user == nil {
		return 0, errno.New(errno.NotFoundError)
	}
	if user.TokenQuota == unlimitedQuota {
		return unlimitedQuota, nil
	}
	remaining := user.TokenQuota - user.UsedTokens
	if remaining < 0 {
		return 0, nil
	}
	return remaining, nil
}

func (s *UserService) SetUserQuota(userID, tokenQuota int64) (bool, error) {
	if userID <= 0 {
		return false, errno.New(errno.ParamsError)
	}
	ok, err := s.userRepo.UpdateQuota(userID, tokenQuota)
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	if !ok {
		return false, errno.New(errno.OperationError)
	}
	return true, nil
}

func (s *UserService) ResetUserQuota(userID int64) (bool, error) {
	if userID <= 0 {
		return false, errno.New(errno.ParamsError)
	}
	ok, err := s.userRepo.ResetUsedTokens(userID)
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	if !ok {
		return false, errno.New(errno.OperationError)
	}
	return true, nil
}

func (s *UserService) DisableUser(userID int64) (bool, error) {
	if userID <= 0 {
		return false, errno.New(errno.ParamsError)
	}
	ok, err := s.userRepo.UpdateStatus(userID, "disabled")
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	if !ok {
		return false, errno.New(errno.OperationError)
	}
	return true, nil
}

func (s *UserService) EnableUser(userID int64) (bool, error) {
	if userID <= 0 {
		return false, errno.New(errno.ParamsError)
	}
	ok, err := s.userRepo.UpdateStatus(userID, "active")
	if err != nil {
		return false, errno.New(errno.SystemError)
	}
	if !ok {
		return false, errno.New(errno.OperationError)
	}
	return true, nil
}

func (s *UserService) GetUserVOList(users []entity.User) []vo.UserVO {
	if len(users) == 0 {
		return make([]vo.UserVO, 0)
	}
	result := make([]vo.UserVO, 0, len(users))
	for _, item := range users {
		itemCopy := item
		result = append(result, s.GetUserVO(&itemCopy))
	}
	return result
}

func (s *UserService) GetEncryptPassword(userPassword string) string {
	sum := md5.Sum([]byte(userPassword + constant.PasswordSalt))
	return hex.EncodeToString(sum[:])
}

func hasBlank(values ...string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return true
		}
	}
	return false
}
