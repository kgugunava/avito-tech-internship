package models

type UserDB struct {
    UserID   int    `db:"user_id"`
    Name     string `db:"name"`
    IsActive bool   `db:"is_active"`
    TeamID   *int   `db:"team_id"`
}