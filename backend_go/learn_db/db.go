package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type TestGo struct {
	ID   int    `json:"id"`
	Num  int    `json:"num"`
	Data string `json:"data"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No env file found")
	}
	connection_string := os.Getenv("DB_URL")
	ctx := context.Background()
	db, err := pgxpool.New(ctx, connection_string)
	if err != nil {
		log.Fatalf("CANT CONNECT TO DB %v", err)
	}
	defer db.Close()
	fmt.Println("CONNECTION POOL READY")

	/*_, err = db.Exec(ctx, `CREATE TABLE IF NOT EXISTS test_go(
							id serial PRIMARY KEY,
							num INTEGER,
							data TEXT)`)
	if err != nil {
		log.Fatalf("Failed to create table %v", err)
	}
	fmt.Println("TABLE test_go CREATED SUCCESSFULLY")
	fmt.Println("test_go(id:serial, num:integer, data: text)")
	log.Println("INSERTING INTO test_go")
	status := insert_into_table(ctx, db, 123, "Samuel")
	if status {
		fmt.Println("INSERTED INTO TABLE")
	} else {
		fmt.Println("ERROR INSERTING")
	}*/
	/*nums := []int{456, 789, 147}
	dats := []string{"ABC", "DEF", "GHI"}
	for index, num := range nums {
		log.Printf("Adding value %v:%v", num, dats[index])
		status := insert_into_table(ctx, db, num, dats[index])
		if status {
			log.Println("INSERTED INTO TABLE")
		} else {
			log.Fatalf("Error inserting at index %v", index)
		}

	}*/
	log.Println("Attempting to fetch 1 row")
	var id, num int
	var data string
	id, num, data = read_from_db(ctx, db)
	fmt.Printf("Received values (%v,%v,%v)", id, num, data)

	log.Println("\nAttempting to fetch multiple rows")

	var values []TestGo
	values, err = read_multiple_from_db(ctx, db)
	if err != nil {
		fmt.Println("Error reading multiple rows")
	} else {
		fmt.Println("test_Go(id, num, data)")
		for _, val := range values {
			fmt.Println("Tuple Received\t: ", val)
		}
	}

}

func insert_into_table(ctx context.Context, db *pgxpool.Pool,
	num int, data string) bool {
	_, err := db.Exec(ctx, `INSERT INTO test_go (num, data) 
							VALUES ($1, $2)`, num, data)
	if err != nil {
		log.Println("Error inserting into table")
		return false
	}
	return true

}

func read_from_db(ctx context.Context, db *pgxpool.Pool) (int, int, string) {
	var id, num int
	var data string
	err := db.QueryRow(ctx, `SELECT id, num, data FROM test_go
							WHERE id = 1`).Scan(&id, &num, &data)
	if err != nil {
		log.Fatalf("Query row failed")
	}
	fmt.Println("fetched 1 row")
	return id, num, data
}

func read_multiple_from_db(ctx context.Context, db *pgxpool.Pool) ([]TestGo, error) {
	rows, err := db.Query(ctx, `SELECT * FROM test_go ORDER BY id`)
	if err != nil {
		log.Fatalf("Error fetching all")
	}

	defer rows.Close()
	var results []TestGo
	for rows.Next() {
		var t TestGo
		if err := rows.Scan(&t.ID, &t.Num, &t.Data); err != nil {
			return nil, err
		}
		results = append(results, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
