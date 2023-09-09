package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {

	//Sl         uint   `gorm:"autoIncrement"`
	//gorm.Model
	Id         uint   `gorm:"primaryKey;autoIncrement"`
	Username   string `gorm:"uniqueIndex;not null"`
	Email      string `gorm:"uniqueIndex"`
	Password   string `gorm:"not null"`
	Admin      bool   `gorm:"default:false;not null"`
	Permission bool   `gorm:"default:true;not null"`
	PersonalID uint   `gorm:"foreignKey"`
}

var DB *gorm.DB

func ConnectToDb() (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}
	var err error
	dsn := "user=postgres password=nashi dbname=project09 host=localhost port=5432 sslmode=disable"

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	DB.AutoMigrate(&User{})

	return DB, nil
}

func CreateUser(user *User) error {
	err := DB.Create(user).Error
	return err
}

func UpdateUser(user User, userid uint) {
	fmt.Println(" admin at update user: ", user)

	err := DB.Where("id = ? ", userid).Updates(user).Error
	if err != nil {
		fmt.Println("error while updating user : ", err)
	}
}
func UpdateUserAdminStatus(userID uint, newAdminStatus bool) error {
	err := DB.Model(&User{}).Where("id = ?", userID).UpdateColumn("Admin", newAdminStatus).Error
	if err != nil {
		fmt.Println("error while updating admin status: ", err)
		return err
	}
	return nil
}

func DeletesUser(userid uint) {
	err := DB.Delete(&User{}, userid).Error
	if err != nil {
		fmt.Println("error hile deleting a user : ", err)
	}
}

func GetUsers(pagelimit, offset int) ([]User, int, error) {
	var users []User
	result := DB.Limit(pagelimit).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	count := GetUserCount()
	return users, count, nil
}

func GetUserCount() int {
	var count int64
	err := DB.Model(&User{}).Count(&count).Error
	if err != nil {
		fmt.Println("error while getting total user count: ", err)
		return 0
	}
	return int(count)
}

func GetSearchUsers(searchWord string, limit, offset int) ([]User, int, error) {
	var users []User
	var count int64
	err := DB.Where("username ILIKE ?", searchWord+"%").Order("username").Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		fmt.Println("error while getting userfor serch :", err)
		return nil, 0, err
	}
	err = DB.Model(&User{}).Where("username ILIKE ?", searchWord+"%").Count(&count).Error
	if err != nil {
		fmt.Println("error while getting total user count: ", err)
		return users, 0, err
	}
	return users, int(count), nil
}

func GetUserDetails(userid uint) User {
	var user User
	err := DB.Where("id = ?", userid).First(&user).Error
	if err != nil {
		fmt.Println("error while getting userdetails : ", err)
	}
	return user
}

func GetUserid(username, password string) (uint, error) {
	var user *User
	err := DB.Where("username = ? AND password = ?", username, password).First(&user).Error
	fmt.Println("printing id : ", user.Id)
	if err != nil {
		fmt.Println("error while getting usersdata :", err)
		return 0, err
	}
	return user.Id, err
}

func Getrole(username, password string) bool {
	var user *User
	err := DB.Where("username = ? AND password = ? ", username, password).First(&user).Error
	if err != nil {
		fmt.Println("error while getting role : ", err)
		return false
	}
	return user.Admin
}

func Getusername(userid uint) string {
	var user *User
	err := DB.Where("id = ?", userid).First(&user).Error
	if err != nil {
		fmt.Println("error while getting username : ", err)
	}
	return user.Username
}

func CheckforUsername(username string) bool {
	var userCount int64
	fmt.Println("username in check for username :", username)
	err := DB.Model(&User{}).Where("Username = ?", username).Count(&userCount).Error
	if err != nil {
		fmt.Println("error while checking for username :", err)
		return false
	}

	return userCount > 0
}

func CheckforEmail(email string) bool {
	var emailCount int64
	err := DB.Model(&User{}).Where("Email = ?", email).Count(&emailCount).Error
	if err != nil {
		fmt.Println("error while looking for email :", err)
		return false
	}
	return emailCount > 0
}

func Verifypassword(username, password string) bool {
	var user *User
	err := DB.Where("username = ? ", username).First(&user).Error
	if err != nil {
		fmt.Println("error while verifying password : ", err)
		return false
	}
	if user.Password == password {
		return true
	}
	return false
}
