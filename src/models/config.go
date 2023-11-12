package models

import (
	"sync"

	"gorm.io/gorm"
)

// ModuleConfig use for Automigrant Tables.
type ModuleConfig struct {
	DB *gorm.DB
}

// NewModuleConfig Return New Module Config.
func NewModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

func (c *ModuleConfig) TableMigration(wg *sync.WaitGroup) {
	// err := c.DB.Migrator().CreateConstraint(&Group{}, "Users")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = c.DB.Migrator().CreateConstraint(&Group{}, "fk_users_groups")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
