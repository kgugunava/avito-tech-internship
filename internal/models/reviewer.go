package models

type ReviewerDB struct {
    PRID   int `db:"pr_id"`
    UserID int `db:"user_id"`
}