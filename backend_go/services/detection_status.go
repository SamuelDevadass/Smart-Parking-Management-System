/** Grouped By HANDLER->operation-id
 */
package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const pythonBaseURL = "http://localhost:8001"
const detectionBasePath = `C:\Users\Samuel\Desktop\SAM\CODING\PYTHON\LicensePlateRecognition\License-Plate-Detector\backend`

// --------------------------VEHICLE_ENTRY_EXIT-------------------------------------------

// ----------------mark-entry--------------------
type DetectionStatus struct {
	Status       string `json:"status"`
	LicensePlate string `json:"license_plate"`
	FolderPath   string `json:"folder_path"`
}

func GetDetectionStatus() (*DetectionStatus, error) {
	url := pythonBaseURL + "/api/detection/status"
	log.Println("Calling Python detection service:", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	log.Println("Python response status:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("python detection service returned %s", resp.Status)
	}

	var status DetectionStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode detection response: %w", err)
	}
	log.Printf("Detection status received: %+v\n", status)
	return &status, nil
}

func CreateFolderPath(license_plate string, folder_path string) string {
	if folder_path == "" || license_plate == "" {
		log.Println("Error fetching detection status")
		return ""
	}

	// Python gives us something like:
	// scans/2026-10-01_16-56-38
	// Make it relative to the Python detection service.
	sourcePath := filepath.Join(detectionBasePath, folder_path)
	date := time.Now().Format("02-01-2006")
	currentTime := time.Now().Format("15-04-05")
	vehiclePath := filepath.Join("capture_log", date,
		license_plate+"_"+currentTime)

	if err := os.MkdirAll(vehiclePath, 0755); err != nil {
		log.Println("Error creating capture directory:", err)
		return ""
	}
	log.Println("Copying detection files from:", sourcePath)
	log.Println("Copying detection files to:", vehiclePath)
	if err := copyDir(sourcePath, vehiclePath); err != nil {
		log.Println("Error copying capture files:", err)
		return ""
	}
	return vehiclePath
}

func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relativePath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relativePath)
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		return copyFile(path, targetPath)
	})
}

func copyFile(src string, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()
	output, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(output, input); err != nil {
		output.Close()
		return err
	}
	return output.Close()
}
