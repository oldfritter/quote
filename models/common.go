package models

import (
	"time"
)

type CommonModel struct {
	Id        int       `json:"-" gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
