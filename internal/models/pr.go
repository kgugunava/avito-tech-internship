package models

type PullRequestDB struct {
    PRID      int        `db:"pr_id"`
    Name      string     `db:"name"`
    AuthorID  int        `db:"author_id"`
    Status    string     `db:"status"` // OPEN / MERGED
}