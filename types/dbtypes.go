package types

import (
	"database/sql"
	"time"
)

type Section struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type User struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Avatar    string         `json:"avatar"`
	Age       sql.NullInt32  `json:"age"`
	Site      sql.NullString `json:"site"`
	AboutMe   sql.NullString `json:"about_me"`
	Email     sql.NullString `json:"email"`
	PubEmail  sql.NullBool   `json:"pub_email"`
	// not password
}

type Product struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Active      bool      `json:"active"`
	SectionID   int       `json:"section_id"`
	UserID      int       `json:"user_id"`
	Date        time.Time `json:"date"`
	Price       int       `json:"price"`
}