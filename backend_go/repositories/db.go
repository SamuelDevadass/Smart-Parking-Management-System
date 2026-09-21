package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"api.com/models"
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

// ------------------------------GET SPOT AVAILABILITY-------------------------------
func GetSpotAvailability(ctx context.Context, wing string) (map[string]int, error) {
	spots_map_dict := make(map[string]int)
	var total_spots, free_spots int = 0, 0
	/*err := DB.QueryRow(ctx, `SELECT COUNT(*) AS total_spots_two_wheeler,
	                        		COUNT(*) FILTER (WHERE availability = TRUE) AS free_spots_two_wheeler
	                        		FROM has_parking_spot WHERE wing = $1 AND size = $2`,
			wing, "Two Wheeler").Scan(total_spots, free_spots) //order matters here: (total, free)
			// this is dngerous because pgsql and pgx will return the type for total_spots_two_wheeler as BIG INT
			// which disagrees with normal golang int syntax
			Even if you dont specify a variable name in the query it will return BIG INT
			// Typecast to int using ::int
			// ::int should appear just before the as clause or after the entire filter is defined*/

	err := DB.QueryRow(ctx, `SELECT COUNT(*)::int AS total_spots_two_wheeler,
                        		(COUNT(*) FILTER (WHERE availability = TRUE))::int AS free_spots_two_wheeler
                        		FROM has_parking_spot WHERE wing = $1 AND size = $2`,
		wing, "Two Wheeler").Scan(&total_spots, &free_spots)
	if err != nil {
		log.Println("Query failed to fetch spots\n", err)
		return nil, err
	}
	spots_map_dict["total_spots_two_wheeler"] = total_spots
	spots_map_dict["free_spots_two_wheeler"] = free_spots
	err = DB.QueryRow(ctx, `SELECT COUNT(*)::int AS total_spots_four_wheeler,
                        		(COUNT(*) FILTER (WHERE availability = TRUE))::int AS free_spots_four_wheeler
                        		FROM has_parking_spot WHERE wing = $1 AND size = $2`,
		wing, "Four Wheeler").Scan(&total_spots, &free_spots) //order matters here: (total, free)
	if err != nil {
		log.Println("Query failed to fetch spots\n", err)
		return nil, err
	}
	spots_map_dict["total_spots_four_wheeler"] = total_spots
	spots_map_dict["free_spots_four_wheeler"] = free_spots

	return spots_map_dict, nil
}

// -----------------------------------GET AVAILABLE SPOTS--------------------------------------------------------------
func GetAvailableSpots(ctx context.Context, wing string, centre_id int, size string) ([]map[string]string, error) {
	var floor, spot_number, size_r string
	var result []map[string]string

	rows, err := DB.Query(ctx, `SELECT floor, spot_number, size FROM has_parking_spot
               					WHERE centre_id = $1 AND wing = $2 AND size = $3 AND availability = True`,
		centre_id, wing, size)
	if err != nil {
		log.Println("Failed to fetch data ...\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		single_row := make(map[string]string) //maps are passed by reference so always declare dynamically
		if err := rows.Scan(&floor, &spot_number, &size_r); err != nil {
			log.Println("Error while fetching free spots\n", err)
			return nil, err
		}
		single_row["floor"] = floor
		single_row["spot_number"] = spot_number
		single_row["size"] = size_r
		result = append(result, single_row)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error while fetching free spots \n", err)
		return nil, err
	}
	return result, nil
}

// ---------------GET VEHICLE--------------
func GetVehicle(ctx context.Context, license_plate string) (*models.VehicleDetails, error) {
	var vehicle models.VehicleDetails
	err := DB.QueryRow(ctx, `SELECT owner_id::int, model, colour, type
                        		FROM owns_vehicle WHERE license_number = $1`,
		license_plate).Scan(&vehicle.OwnerID, &vehicle.Model, &vehicle.Colour, &vehicle.Type)
	if err != nil {
		log.Println("Query failed to fetch Vehicle Details\n", err)
		return nil, err
	}

	err = DB.QueryRow(ctx, `SELECT phone FROM owner_phone WHERE owner_id = $1 LIMIT 1`, vehicle.OwnerID).Scan(&vehicle.Phone)
	if err != nil {
		log.Println("Query failed to fetch User Phone\n", err)
		return nil, err
	}

	err = DB.QueryRow(ctx, `SELECT name FROM owner WHERE owner_id = $1 LIMIT 1`, vehicle.OwnerID).Scan(&vehicle.Name)
	if err != nil {
		log.Println("Query failed to fetch User Phone\n", err)
		return nil, err
	}
	return &vehicle, nil
}

func SaveVehicle(ctx context.Context, vehicle *models.VehicleDetails, license_plate string) (bool, error) {
	_, err := DB.Exec(ctx, `INSERT INTO owner (owner_id, name) VALUES ($1, $2)
                        		ON CONFLICT (owner_id) DO NOTHING`, vehicle.OwnerID, vehicle.Name)
	if err != nil {
		log.Println("Insert failed:\n", err)
		return false, err
	}
	_, err = DB.Exec(ctx, `INSERT INTO owner_phone (owner_id, phone) VALUES ($1, $2)
                        		ON CONFLICT (owner_id) DO NOTHING`, vehicle.OwnerID, vehicle.Phone)
	if err != nil {
		log.Println("Insert failed:\n", err)
		return false, err
	}
	_, err = DB.Exec(ctx, `INSERT INTO owns_vehicle (owner_id, license_number, model, colour, type)
                        		VALUES ($1, $2, $3, $4, $5)
                        ON CONFLICT (license_number) DO NOTHING`, vehicle.OwnerID, license_plate, vehicle.Model, vehicle.Colour, vehicle.Type)
	if err != nil {
		log.Println("Insert failed:\n", err)
		return false, err
	}
	return true, nil
}

// Mark Entry
func MarkEntry(ctx context.Context, input *models.MarkEntryInput, entry_time time.Time) (bool, error) {
	_, err := DB.Exec(ctx, `INSERT INTO parking_log 
							(entry_time,license_number, centre_id, 
                     		wing, floor, spot_number, image_folder_path)
                    		VALUES ($1,$2,$3,$4,$5,$6,$7)`, entry_time, input.LicensePlate, input.CentreID, input.Wing, input.Floor, input.SpotNumber, input.FolderPath)
	if err != nil {
		log.Println("Insert failed:\n", err)
		return false, err
	}
	_, err = DB.Exec(ctx, `UPDATE has_parking_spot SET availability = False
                        	WHERE centre_id = $1 AND wing = $2 AND floor = $3 AND spot_number = $4`,
		input.CentreID, input.Wing, input.Floor, input.SpotNumber)
	if err != nil {
		log.Println("Update failed:\n", err)
		return false, err
	}
	return true, nil
}

func GetActiveSession(ctx context.Context, input *models.GetVehicleInput) (*models.GetActiveSessionResponse, error) {
	resp := &models.GetActiveSessionResponse{}
	err := DB.QueryRow(ctx, `SELECT entry_time, centre_id, wing, floor, spot_number
                        		FROM parking_log WHERE license_number = $1 
                        		AND exit_time IS NULL ORDER BY entry_time DESC LIMIT 1`, input.LicensePlate).Scan(&resp.Body.EntryTime, &resp.Body.CentreID,
		resp.Body.Wing, resp.Body.Floor, resp.Body.SpotNumber)
	if err != nil {
		log.Println("Error fetching session details\n", err)
		return nil, err
	}
	return resp, nil
}

// Get latest bill
func GetLatestBill(ctx context.Context, input *models.GetVehicleInput) (*models.GetLatestBillResponse, error) {
	resp := &models.GetLatestBillResponse{}
	err := DB.QueryRow(ctx, `SELECT entry_time, exit_time, duration, amount
                        	FROM parking_log WHERE license_number = $1 AND exit_time IS NOT NULL
                       	 	ORDER BY exit_time DESC LIMIT 1`, input.LicensePlate).
		Scan(&resp.Body.EntryTime, &resp.Body.ExitTime, &resp.Body.Duration, &resp.Body.Amount)
	if err != nil {
		log.Println("Error fetching bill details\n", err)
		return nil, err
	}
	err = DB.QueryRow(ctx, `SELECT o.name FROM owner o 
                        		JOIN owns_vehicle v ON v.owner_id = o.owner_id
                        		WHERE v.license_number = $1`, input.LicensePlate).Scan(&resp.Body.OwnerName)
	if err != nil {
		log.Println("Error fetching bill details\n", err)
		return nil, err
	}
	return resp, nil
}
