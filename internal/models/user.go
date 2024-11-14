package models

type User struct {
	Base
	Username     string  `json:"username" gorm:"type:varchar(100)"`
	Email        string  `json:"email" gorm:"type:varchar(100);unique_index"`
	Organization *string `json:"organization" grom:"type:varchar(100)"`
	Password     string  `json:"-"` //password will be omitted during marshalling to struct
}
