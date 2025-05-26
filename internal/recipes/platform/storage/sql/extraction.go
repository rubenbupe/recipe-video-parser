package sql

const (
	sqlExtractionTable = "recipe_extractions"
)

type sqlExtraction struct {
	ID        string `db:"id"`
	UserID    string `db:"user_id"`
	Data      string `db:"data"`
	Metadata  string `db:"metadata"`
	CreatedAt string `db:"created_at"`
}
