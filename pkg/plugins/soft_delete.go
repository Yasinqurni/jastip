package plugins

import "gorm.io/gorm"

// Redefined soft delete plugin.
// Change this SoftDelete type according library/pakcage that used in this project.
type SoftDelete struct {
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
