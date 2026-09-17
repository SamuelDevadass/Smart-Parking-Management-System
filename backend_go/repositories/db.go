package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Gloabl Pool Instance
var DB *pgxpool.Pool

func Init_DB(ctx context.Context, connection_string string) {
	db, err := pgxpool.New(ctx, connection_string)
	if err != nil {
		log.Fatalf("CANT CONNECT TO DB %v", err)
	}
	DB = db
	fmt.Println("CONNECTION POOL READY")
}

// ------------------LIST WINGS------------------------
func ListWings(ctx context.Context) ([]string, error) {
	var result []string
	var value string
	rows, err := DB.Query(ctx, `SELECT DISTINCT wing FROM has_wing_floor`)
	if err != nil {
		log.Fatal("Failed to fetch data ...\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&value); err != nil {
			log.Println("Error while fetching wings \n", err)
			return nil, err
		}
		result = append(result, value)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error while fetching wings \n", err)
		return nil, err
	}
	return result, nil
}

func GetCentreForWings(ctx context.Context, wing string) (string, error) {
	var centre string
	err := DB.QueryRow(ctx, `SELECT centre_id FROM has_wing_floor 
										WHERE wing = ($1)`, wing).Scan(&centre)
	if err != nil {
		log.Fatalf("Query row failed")
	}
	return centre, err
}
