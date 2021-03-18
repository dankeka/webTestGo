package types

import "time"

type Section struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type Announcement struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Active      bool      `json:"active"`
	SectionID   int       `json:"section_id"`
	UserID      int       `json:"user_id"`
	Date        time.Time `json:"date"`
}