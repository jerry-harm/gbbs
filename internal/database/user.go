package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func hashPass(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err.Error())
	}
	return string(hash)
}

func comparePass(hashed string, password string) bool {
	hashedbyte := []byte(hashed)
	passwordbyte := []byte(password)

	err := bcrypt.CompareHashAndPassword(hashedbyte, passwordbyte)
	if err != nil {
		log.Println(err.Error())
		return false
	} else {
		return true
	}
}

func Register(name string, password string, email string) error {
	user := User{Name: name, PasswordHash: hashPass(password), Email: email}
	err := DB.Create(&user).Error
	return err

}

func Login(name string, password string) bool {
	var user User
	log.Println("pass ", password, " user ", name)
	if err := DB.Where("name = ?", name).First(&user); err.Error != nil {
		return false
	}
	return comparePass(user.PasswordHash, password)
}
