package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Caption  string
	ImageUrl string
	Creator  string
	Comments []Comment
}
