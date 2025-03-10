package repository

import (
	"notex/model"
	"notex/pkg/database"
	"strings"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.GetDB(),
	}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// FindByID 通过ID查找用户
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 通过邮箱查找用户
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll 查找所有用户
func (r *UserRepository) FindAll(page, pageSize int, role, status, search string) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.Model(&model.User{})

	// 应用过滤条件
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		search = "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(email) LIKE ? OR LOWER(nickname) LIKE ?", search, search)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateLastLogin 更新最后登录时间
func (r *UserRepository) UpdateLastLogin(userID uint) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("last_login", gorm.Expr("NOW()")).Error
}

// MarkEmailAsVerified 标记邮箱为已验证
func (r *UserRepository) MarkEmailAsVerified(userID uint) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("email_verified", true).Error
}

// UpdateEmail 更新邮箱
func (r *UserRepository) UpdateEmail(userID uint, newEmail string) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"email":          newEmail,
			"email_verified": false,
		}).Error
}

// IsUsernameExistsExcept 检查用户名是否已被使用（排除指定用户）
func (r *UserRepository) IsUsernameExistsExcept(username string, excludeUserID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).
		Where("username = ? AND id != ?", username, excludeUserID).
		Count(&count).Error
	return count > 0, err
}
