package repository

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/model"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found, let service handle 404
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindAll(filter UserFilter) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// Start query construction
	query := r.db.Model(&model.User{})

	// 1. Search Logic
	if filter.Search != "" {
		searchLike := "%" + filter.Search + "%"
		switch filter.Scope {
		case "name":
			query = query.Where("name LIKE ?", searchLike)
		case "email":
			query = query.Where("email LIKE ?", searchLike)
		case "id":
			// Exact match for ID
			if id, err := strconv.Atoi(filter.Search); err == nil {
				query = query.Where("id = ?", id)
			} else {
				// If scope is ID but value is not a number, return empty or fail safe
				query = query.Where("1 = 0")
			}
		default: // "all" or empty
			// Search in Name, Email, or ID
			sql := "name LIKE ? OR email LIKE ?"
			args := []interface{}{searchLike, searchLike}
			
			// If search term is numeric, include ID search
			if id, err := strconv.Atoi(filter.Search); err == nil {
				sql += " OR id = ?"
				args = append(args, id)
			}
			query = query.Where(sql, args...)
		}
	}

	// 2. Filter by Role
	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}

	// 3. Sorting
	if filter.SortBy != "" {
		sortParams := strings.Split(filter.SortBy, ":")
		if len(sortParams) == 2 {
			field := sortParams[0]
			order := strings.ToUpper(sortParams[1])
			
			// Validate Order
			if order != "ASC" && order != "DESC" {
				order = "ASC"
			}

			// Validate Field (Whitelisting to prevent SQL injection)
			allowedFields := map[string]string{
				"id":         "id",
				"name":       "name",
				"email":      "email",
				"role":       "role",
				"created_at": "created_at",
			}

			if dbField, ok := allowedFields[field]; ok {
				query = query.Order(fmt.Sprintf("%s %s", dbField, order))
			}
		}
	} else {
		// Default sort
		query = query.Order("id ASC")
	}

	// Count total records (before pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 4. Pagination
	offset := (filter.Page - 1) * filter.Limit
	if err := query.Offset(offset).Limit(filter.Limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepo) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepo) Delete(id uint) error {
	result := r.db.Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *userRepo) IsEmailTaken(email string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&model.User{}).Where("email = ?", email)
	
	if excludeID != 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}