package database

import "gorm.io/gorm"

type Movie struct {
	gorm.Model
	Title     string `json:"title"`
	BechdelID int    `gorm:"unique" json:"id"`
	IMDBID    string `gorm:"uniqueIndex" json:"imdbid"`
	Rating    int    `json:"rating"`
}
