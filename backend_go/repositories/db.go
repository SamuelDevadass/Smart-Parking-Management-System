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

// -------------------LIST WINGS------------------------
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

// -----------------------GET CENTRE-ID FOR WING-------------------------
func GetCentreForWings(ctx context.Context, wing string) (int, error) {
	var centre int = 0
	err := DB.QueryRow(ctx, `SELECT centre_id FROM has_wing_floor 
										WHERE wing = ($1)`, wing).Scan(&centre)
	if err != nil {
		log.Println("Query row failed")
	}
	return centre, err
}

// GET SPOT AVAILABILITY
func GetSpotAvailability(ctx context.Context, wing string) (map[string]int, error) {
	spots_map_dict := make(map[string]int)
	var total_spots, free_spots int = 0, 0
	err := DB.QueryRow(ctx, `SELECT COUNT(*) AS total_spots_two_wheeler,
                        		COUNT(*) FILTER (WHERE availability = TRUE) AS free_spots_two_wheeler
                        		FROM has_parking_spot WHERE wing = $1 AND size = $2`,
		wing, "Two Wheeler").Scan(total_spots, free_spots) //order matters here: (total, free)
	if err != nil {
		log.Println("Query failed to fetch spots\n", err)
		return nil, err
	}
	spots_map_dict["total_spots_two_wheeler"] = total_spots
	spots_map_dict["free_spots_two_wheeler"] = free_spots
	err = DB.QueryRow(ctx, `SELECT COUNT(*) AS total_spots_four_wheeler,
                        		COUNT(*) FILTER (WHERE availability = TRUE) AS free_spots_four_wheeler
                        		FROM has_parking_spot WHERE wing = $1 AND size = $2`,
		wing, "Four Wheeler").Scan(total_spots, free_spots) //order matters here: (total, free)
	if err != nil {
		log.Println("Query failed to fetch spots\n", err)
		return nil, err
	}
	spots_map_dict["total_spots_four_wheeler"] = total_spots
	spots_map_dict["free_spots_four_wheeler"] = free_spots

	return spots_map_dict, nil
}
