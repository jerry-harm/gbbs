package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Group struct {
	Id    uint    `gorm:"primarykey"`
	Name  string  `gorm:"unique"`
	Users []*User `gorm:"many2many:user_groups;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	Id           uint   `gorm:"primarykey"`
	Name         string `gorm:"unique"`
	PasswordHash string
	Groups       []*Group `gorm:"many2many:user_groups;"`
	Email        string

	Lastlogin sql.NullTime

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Area struct {
	Name           string
	MaxLength      uint
	WriteGroupIds  []*Group `gorm:"many2many:area_writegroups;"`
	ReadGroupIds   []*Group `gorm:"many2many:area_readgroups;"`
	ManageGroupIds []*Group `gorm:"many2many:area_managegroups;"`

	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Post struct {
	Title  string
	Text   string `gorm:"type:text"`
	AreaId uint
	Area   Area
	UserId uint
	User   User

	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Board struct {
	Text   string `gorm:"type:text"`
	UserId uint
	User   User

	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(viper.GetString("database")), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	exist := db.Migrator().HasTable(&Group{})
	err = db.AutoMigrate(&Group{}, &User{}, &Area{}, &Post{}, &Board{})
	if err != nil {
		log.Fatalln(err)
	}
	DB = db
	if !exist {
		adminGroup := Group{Name: "Admin"}
		DB.Create(&adminGroup)
		userGroup := Group{Name: "User"}
		DB.Create(&userGroup)
	}
	return db
}
