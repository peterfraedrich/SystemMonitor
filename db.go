package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenDB(path string, drop bool) (*gorm.DB, error) {
	d, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if drop {
		sql := `PRAGMA writable_schema = 1; DELETE FROM sqlite_master; PRAGMA writable_schema = 0; VACUUM; PRAGMA integrity_check;`
		d.Exec(sql)
	}
	d.AutoMigrate(
		&SystemInfo{},
		&BasicMetrics{},
		&ProcessMetrics{},
		&EventsLog{},
		&ErrorLog{},
	)
	return d, nil
}
