package model

import "github.com/google/uuid"

const TableNameAuth = "auths"

// Auth mapped from table <auths>
type Auth struct {
	ID        string `gorm:"column:id;primaryKey" json:"id"`
	Name      string `gorm:"column:name;not null;unique" json:"name"`
	Email     string `gorm:"column:email;not null;unique" json:"email"`
	Role      string `gorm:"column:role;not null" json:"role"`
	Password  string `gorm:"column:password;not null;select:false" json:"-"`
	CreatedAt int64  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt int64  `gorm:"column:updated_at;not null" json:"updated_at"`
	CreatedBy string `gorm:"column:created_by;not null" json:"created_by"`
	UpdatedBy string `gorm:"column:updated_by;not null" json:"updated_by"`
}

// TableName Auth's table name
func (*Auth) TableName() string {
	return TableNameAuth
}

func (auth *Auth) Insert() error {
	auth.ID = uuid.New().String()
	return DB.Create(auth).Error
}

func (auth *Auth) Update() error {
	return DB.Model(auth).Updates(auth).Error
}

func (auth *Auth) Delete() error {
	return DB.Delete(auth).Error
}

func GetAllAuths() ([]*Auth, error) {
	var auths []*Auth
	return auths, DB.Find(&auths).Error
}

func IsEmailExists(email string) (bool, error) {
	var auth Auth
	return DB.Where("email = ?", email).First(&auth).RowsAffected > 0, nil
}

func GetAuthByEmail(email string) (*Auth, error) {
	var auth Auth
	return &auth, DB.Where("email = ?", email).First(&auth).Error
}

func GetAuthByID(id string) (*Auth, error) {
	var auth Auth
	return &auth, DB.Where("id = ?", id).First(&auth).Error
}
