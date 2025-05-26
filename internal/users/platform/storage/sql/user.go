package sql

const (
	sqlUserTable = "users"
)

type sqlUser struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	ApiKey    string `db:"api_key"`
	CreatedAt string `db:"created_at"`
}
