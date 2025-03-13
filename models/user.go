package models

import "time"

type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Password  string    `json:"-"` // put it as - so that it doesnt get exposed in API responses
    CreatedAt time.Time `json:"created_at"`
}
