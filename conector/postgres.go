package conector

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type PostgresConnection struct {
	con *pgx.Conn
}

func (conn *PostgresConnection) Connect(ctx context.Context, dbDsn string) error {
	var err error = nil
	conn.con, err = pgx.Connect(ctx, dbDsn)
	if err != nil {
		return err
	}
	return nil

}

func (conn *PostgresConnection) GetData(ctx context.Context, collectionName string) ([]map[string]any, error) {
	query := fmt.Sprintf("SELECT * FROM %s", collectionName)
	rows, err := conn.con.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fields := rows.FieldDescriptions()
	result := make([]map[string]any, 0)

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		rowMap := make(map[string]any)
		for i, field := range fields {
			rowMap[field.Name] = values[i]
		}
		result = append(result, rowMap)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (conn *PostgresConnection) WriteData(ctx context.Context, tablenName string, data []map[string]any) error {
	tx, err := conn.con.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, row := range data {
		columns := make([]string, 0, len(row))
		values := make([]any, 0, len(row))
		placeholders := make([]string, 0, len(row))

		i := 1
		for col, val := range row {
			columns = append(columns, col)
			values = append(values, val)
			placeholders = append(placeholders, fmt.Sprintf("$%d", i))
			i++
		}

		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tablenName, join(columns, ", "), join(placeholders, ", "))
		_, err := tx.Exec(ctx, query, values...)
		if err != nil {
			return err
		}

	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func join(elements []string, delimiter string) string {
	result := ""
	for i, element := range elements {
		if i > 0 {
			result += delimiter
		}
		result += element
	}
	return result
}

func (conn *PostgresConnection) Close(ctx context.Context) error {

	return conn.con.Close(ctx)
}
