package database

import (
	"github.com/charmbracelet/log"
	"golang.org/x/crypto/bcrypt"
)

func hashPass(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Info(err.Error())
	}
	return string(hash)
}

func comparePass(hashed string, password string) bool {
	hashedbyte := []byte(hashed)
	passwordbyte := []byte(password)

	err := bcrypt.CompareHashAndPassword(hashedbyte, passwordbyte)
	if err != nil {
		log.Debug(err.Error())
		return false
	} else {
		return true
	}
}

func Register(name string, password string, email string) error {
	user := User{Name: name, PasswordHash: hashPass(password), Email: email}
	err := DB.Create(&user).Error
	if err != nil {
		return err
	}
	// first user be admin
	if user.Id == 1 {
		var adminGroup Group
		DB.Where("name = ?", "Admin").First(&adminGroup)
		user.Groups = append(user.Groups, &adminGroup)
	}

	var userGroup Group
	DB.Where("name = ?", "User").First(&userGroup)

	user.Groups = append(user.Groups, &userGroup)
	err = DB.Save(user).Error
	return err
}

func Login(name string, password string) bool {
	var user User
	if err := DB.Where("name = ?", name).First(&user); err.Error != nil {
		return false
	}
	return comparePass(user.PasswordHash, password)
}

func (u User) Board(text string) {
	if len(text) == 0 {
		return
	}
	board := Board{
		Text:   text,
		UserId: u.Id,
		User:   u,
	}
	DB.Create(&board)
}

func (u User) Post(area Area, title string, text string) {
	if len(text) == 0 && len(title) == 0 {
		return
	}
	create := false
	for _, ug := range u.Groups {
		if ug.Name == "Admin" {
			create = true
			break
		}
		for _, ag := range area.WriteGroups {
			if ug.Id == ag.Id {
				create = true
				break
			}
		}
	}
	if create {
		post := Post{
			Title:  title,
			Text:   text,
			AreaId: area.ID,
			UserId: u.Id,
			User:   u,
		}
		DB.Create(&post)
	}
}
