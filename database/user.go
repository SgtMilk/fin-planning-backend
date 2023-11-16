package database

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct{
	gorm.Model

	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`

	// connections
	// belongs to relationship
	OptionsID uint `gorm:"not null" json:"optionsID"`
	Options Options `gorm:"not null" json:"options"`

	// has many relationship
	MonthlyExpenses []MonthlyExpense `json:"monthlyExpenses"`
}

func (user *User) Create() (*User, error) {
	options, err := CreateDefaultOptions()

	if err != nil{
		return &User{}, err
	}

	user.OptionsID = options.ID

	err = Database.Create(&user).Error

	if err != nil{
		options.Delete()
		return &User{}, err
	}
	return user, nil
}

func (user *User) ValidatePassword(password string) error {
    return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func FindUserById(id uint) (User, error){
	var user User
	err := Database.Where("id=?", id).Find(&user).Error

	if err != nil{
		return User{}, err
	}

	return user, nil
}

func FindUserByUsername(username string) (User, error){
	usernameHash, err := generateHash(username)
	if err != nil {
        return User{}, err
    }

	var user User
	err = Database.Where("username=?", usernameHash).Find(&user).Error

	if err != nil{
		return User{}, err
	}

	return user, nil
}


// ==============================================
// 					HOOKS
// ==============================================

func (user *User) BeforeSave(*gorm.DB) error{
	usernameHash, err := generateHash(user.Username)
	if err != nil {
        return err
    }

	passwordHash, err := generateHash(user.Password)
	if err != nil {
        return err
    }
    
	user.Username = usernameHash
    user.Password = passwordHash
    return nil
}

// ==============================================
// 					HELPERS
// ==============================================

func generateHash(input string) (string, error){
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	return string(hash), err
}

