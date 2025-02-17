package model

import "gorm.io/gorm"

const TableNameLiveConf = "live_confs"

// LiveConf mapped from table <live_confs>
type LiveConf struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	RoomDisplayID string `gorm:"column:room_display_id;not null" json:"room_display_id"`
	URL           string `gorm:"column:url;not null" json:"url"`
	Name          string `gorm:"column:name" json:"name"`
	ModifiedOn    int64  `gorm:"column:modified_on" json:"modified_on"`
	CreatedOn     int64  `gorm:"column:created_on" json:"created_on"`
	ModifiedBy    string `gorm:"column:modified_by" json:"modified_by"`
	CratedBy      string `gorm:"column:crated_by" json:"crated_by"`
	Enable        bool   `gorm:"column:enable;not null;" json:"enable"`
}

// TableName LiveConf's table name
func (*LiveConf) TableName() string {
	return TableNameLiveConf
}

func (conf *LiveConf) Insert(db *gorm.DB) error {
	return db.Create(conf).Error
}

func (conf *LiveConf) Update(db *gorm.DB) error {
	return db.Save(conf).Error
}

func DeleteLiveConfById(id int64) error {
	return DB.Delete(&LiveConf{ID: id}).Error
}

func GetLiveConfById(id int64) (*LiveConf, error) {
	var conf LiveConf
	return &conf, DB.Where("id = ?", id).First(&conf).Error
}

func GetAllLiveConf() ([]*LiveConf, error) {
	var confs []*LiveConf
	return confs, DB.Find(&confs).Error
}
