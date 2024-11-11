package query

import (
	"context"
	"love-remittance-be-apps/core/config"
	"love-remittance-be-apps/lib/model"

	"github.com/jackc/pgx/v5"
)

type Query struct {
}

func Row(sql string, args ...any) (pgx.Row, *model.ErrorData) {
	//DB
	db, dbErr := config.GetDBRemit()
	if dbErr != nil {
		return nil, &model.ErrorData{
			Title:       "Query Failed",
			Description: "Unable to connect Database",
		}
	}

	defer db.Close(context.Background())
	row := db.QueryRow(context.Background(), sql, args...)
	return row, nil
}

func ExecUpdate(sql string, args ...any) (int, *model.ErrorData) {
	//DB
	db, dbErr := config.GetDBRemit()
	if dbErr != nil {
		return 0, &model.ErrorData{
			Title:       "Query Failed",
			Description: "Unable to connect Database",
		}
	}

	defer db.Close(context.Background())
	commandTag, err := db.Exec(context.Background(), sql, args...)
	if err != nil {
		return 0, &model.ErrorData{
			Description: "Unable to connect Database",
			Title:       err.Error(),
		}
	}
	if commandTag.RowsAffected() != 1 {
		return 0, &model.ErrorData{
			Description: "No row found to update",
			Title:       "Query Failed",
		}
	}
	return int(commandTag.RowsAffected()), nil
}

func ExecInsert(sql string, args ...any) (int, *model.ErrorData) {
	//DB
	db, dbErr := config.GetDBRemit()
	if dbErr != nil {
		return 0, &model.ErrorData{
			Title:       "Query Failed",
			Description: "Unable to connect Database",
		}
	}

	defer db.Close(context.Background())
	commandTag, err := db.Exec(context.Background(), sql, args...)
	if err != nil {
		return 0, &model.ErrorData{
			Description: "Unable to connect Database",
			Title:       "Query Failed",
		}
	}
	if commandTag.RowsAffected() != 1 {
		return 0, &model.ErrorData{
			Description: "No row found to insert",
			Title:       "Query Failed",
		}
	}

	return int(commandTag.RowsAffected()), nil
}

func Rows(sql string, args ...any) (pgx.Rows, *model.ErrorData) {
	//DB
	db, dbErr := config.GetDBRemit()
	if dbErr != nil {
		return nil, &model.ErrorData{
			Description: "Unable to connect Database",
			Title:       "Query Failed",
			Field:       dbErr.Error(),
		}
	}

	defer db.Close(context.Background())
	rows, err := db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, &model.ErrorData{
			Description: "Query Failed",
			Title:       "Query Failed",
			Field:       err.Error(),
		}
	}
	return rows, nil
}

func RowRC(sql string, args ...any) (pgx.Row, *model.ErrorData) {
	//DB
	db, dbErr := config.GetDBLovePaycode()
	if dbErr != nil {
		return nil, &model.ErrorData{
			Title:       "Query Failed",
			Description: "Unable to connect Database",
		}
	}

	defer db.Close(context.Background())
	row := db.QueryRow(context.Background(), sql, args...)
	return row, nil
}
