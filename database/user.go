package database

import (
	"errors"
	"html"
	"strings"

	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`

	// connections
	// belongs to relationship
	OptionsID uint    `gorm:"not null" json:"optionsID"`
	Options   Options `gorm:"not null" json:"options"`

	// has many relationship
	MonthlyExpenses []MonthlyExpense `json:"monthlyExpenses"`
}

func (user *User) Create() (*User, error) {
	options, err := CreateDefaultOptions()

	if err != nil {
		return &User{}, err
	}

	user.OptionsID = options.ID

	err = Database.Create(&user).Error

	if err != nil {
		options.Delete()
		return &User{}, err
	}
	return user, nil
}

func (user *User) UpdatePassword(oldPassword string, newPassword string) error {
	err := user.ValidatePassword(oldPassword)

	if err != nil {
		return err
	}

	err = Database.Model(&user).Update("password", newPassword).Error

	return err
}

func (user *User) UpdateMonthlyExpenses(monthlyExpenses []MonthlyExpense) error {
	user.MonthlyExpenses = monthlyExpenses
	err := Database.Save(&user).Error
	return err
}

func (user *User) Delete() error {
	// deleting attached optoins
	err := user.Options.Delete()
	if err != nil {
		return err
	}

	// deleting attached monthly
	var monthlyExpenses []MonthlyExpense
	err = Database.Where("user_id = ?", user.ID).Find(&monthlyExpenses).Error
	if err != nil {
		return err
	}

	for _, monthlyExpense := range monthlyExpenses {
		err = monthlyExpense.Delete()
		if err != nil {
			return err
		}
	}

	err = Database.Delete(&user).Error
	return err
}

func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func FindUserById(id uint) (*User, error) {
	var user User
	err := Database.Preload("Options").Where("id=?", id).Find(&user).Error

	if err != nil {
		return &User{}, err
	}

	if user.ID == 0 {
		return &user, errors.New("User doesn't exist")
	}

	return &user, nil
}

func GetUserWithMontlyExpenses(userId uint) (*User, error) {
	var user User
	err := Database.Preload("monthly_expenses").Where("user_id=?", userId).Find(&user).Error

	if err != nil {
		return &User{}, err
	}

	if user.ID == 0 {
		return &user, errors.New("User doesn't exist")
	}

	return &user, nil
}

func FindUserByUsername(username string) (*User, error) {
	var user User
	err := Database.Preload("Options").Where("username=?", username).Find(&user).Error

	if err != nil {
		return &User{}, err
	}

	if user.ID == uint(0) {
		return &user, errors.New("User doesn't exist")
	}

	return &user, nil
}

// ==============================================
// 					HOOKS
// ==============================================

func (user *User) BeforeSave(*gorm.DB) error {
	err := user.assertInput()
	if err != nil {
		return err
	}

	passwordHash, err := GenerateHash(user.Password)
	if err != nil {
		return err
	}

	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	user.Password = passwordHash
	return nil
}

// ==============================================
// 					HELPERS
// ==============================================

func GenerateHash(input string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	return string(hash), err
}

func (user *User) assertInput() error {
	// username
	usernameLengthCheck := len(user.Username) < 8 || len(user.Username) > 256
	passwordLengthCheck := len(user.Password) < 8 || len(user.Password) > 72
	if usernameLengthCheck && passwordLengthCheck {
		return errors.New("username and password not of right size")
	} else if usernameLengthCheck {
		return errors.New("username not of right size")
	} else if passwordLengthCheck {
		return errors.New("password not of right size")
	}

	// strength evaluation
	err := passwordvalidator.Validate(user.Username, 50)
	if err != nil {
		return errors.New("insecure username, try using a longer username")
	}

	return passwordvalidator.Validate(user.Password, 60)
}
