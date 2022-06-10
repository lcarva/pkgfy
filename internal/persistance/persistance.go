package persistance

import (
	"time"

	"github.com/lcarva/pkgfy/internal/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PackageORM struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"primarykey"`
	Alias     string
	URL       string
}

func Save(dbFile string, pkg core.Package) (err error) {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&PackageORM{})
	if err != nil {
		return
	}

	return db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&PackageORM{
		Name:  pkg.Name,
		Alias: pkg.Alias,
		URL:   pkg.URL,
	}).Error
}

func List(dbFile string) (pkgs []core.Package, err error) {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&PackageORM{})
	if err != nil {
		return
	}

	pkgsORM := []PackageORM{}
	err = db.Order("name asc").Find(&pkgsORM).Error
	if err != nil {
		return
	}
	for _, pkgORM := range pkgsORM {
		pkgs = append(pkgs, core.Package{
			Name:  pkgORM.Name,
			Alias: pkgORM.Alias,
			URL:   pkgORM.URL,
		})
	}
	return
}
