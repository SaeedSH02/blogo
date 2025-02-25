package models

import "time"



type User struct {
	ID		 uint 	`gorm:"primaryKey" json:"id"`
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null" json:"password"`
	Articles  []Article `gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}



type Article struct {
	ID 	uint     `json:"id" gorm:"primaryKey"`
	Title string `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	AuthorID  uint      `json:"author_id" gorm:"not null"`
	Author    User      `gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
