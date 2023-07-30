package expenditure

import (
	"database/sql"
	"time"
)

func GetRecords(db *sql.DB, chatId int64, from time.Time) ([]Expense, error) {
	dbRows, err := db.Query("select * from Expenses where ChatId = $1 and CreatedAt > $2",
		chatId,
		from.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer dbRows.Close()

	rows := []Expense{}
	for dbRows.Next() {
		r := Expense{}
		var createdStr string
		err := dbRows.Scan(&r.Id, &r.Category, &r.Amount, &r.ChatId, &createdStr, &r.CreatedBy)
		if err != nil {
			continue
		}
		createdAt, parseErr := time.Parse(time.RFC3339, createdStr)
		if parseErr != nil {
			continue
		}
		r.CreatedAt = createdAt
		rows = append(rows, r)
	}

	return rows, nil
}

func (r *Expense) Save(db *sql.DB) error {
	_, err := db.Exec("insert into Expenses (Category, Amount, ChatId, CreatedAt, CreatedBy) values ($1, $2, $3, $4, $5)",
		r.Category, r.Amount, r.ChatId, r.CreatedAt.Format(time.RFC3339), r.CreatedBy)
	return err
}
