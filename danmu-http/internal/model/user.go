package model

const TableNameUser = "users"

// User mapped from table <users>
type User struct {
	ID        uint64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID    uint64 `gorm:"column:user_id;not null;index" json:"user_id"`
	DisplayID string `gorm:"column:display_id;not null" json:"display_id"`
	UserName  string `gorm:"column:user_name;not null" json:"user_name"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}

func GetAllUsersPage(page int, pageSize int) ([]*User, int64, error) {
	var users []*User
	var total int64
	db := DB.Model(&User{})
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply ordering and pagination
	err = db.Order("user_id desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error
	return users, total, err
}

func SearchUser(page int, pageSize int, keyword string) ([]*User, int64, error) {
	var users []*User
	var total int64
	db := DB.Model(&User{})
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	if keyword != "" {
		db.Where("user_name LIKE ? or display_id LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	err = db.Order("user_id desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
