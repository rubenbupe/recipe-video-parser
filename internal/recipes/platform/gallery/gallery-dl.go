package gallery

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type DownloadResult struct {
	FilePath    string
	Url         string
	Extension   string
	MimeType    string
	Description string
}

func DownloadFile(url, id, downloadDir, configFile string) (DownloadResult, error) {
	args := []string{"--write-metadata", "-D", downloadDir, "-f", fmt.Sprintf("%s.{extension}", id)}
	if configFile != "" {
		args = append(args, "-c", configFile)
	}
	args = append(args, url)
	cmd := exec.Command("gallery-dl", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return DownloadResult{}, fmt.Errorf("failed to download video: %w, details: %s", err, output)
	}

	filePath := strings.TrimSpace(strings.TrimPrefix(string(output), "# "))

	if filePath == "" {
		return DownloadResult{}, fmt.Errorf("video file not found after download")
	}

	jsonPath := filePath + ".json"
	jsonFile, err := os.ReadFile(jsonPath)
	if err != nil {
		RemoveFile(filePath) // Clean up the downloaded file if metadata read fails
		return DownloadResult{}, fmt.Errorf("failed to read metadata json: %w", err)
	}

	desc := ""
	extension := ""
	var meta map[string]interface{}
	if err := json.Unmarshal(jsonFile, &meta); err == nil {
		// contents is an array of {desc: string}. Join all descs with newlines
		if d, ok := meta["contents"].([]interface{}); ok {
			var descs []string
			for _, item := range d {
				if m, ok := item.(map[string]interface{}); ok {
					if descValue, ok := m["desc"].(string); ok {
						descs = append(descs, descValue)
					}
				}
			}
			desc = strings.Join(descs, "\n")
		} else if d, ok := meta["desc"].(string); ok {
			desc = d
		} else if d, ok := meta["description"].(string); ok {
			desc = d
		}

		if d, ok := meta["extension"].(string); ok {
			extension = d
		} else {
			RemoveFile(filePath) // Clean up the downloaded file if metadata read fails
			return DownloadResult{}, fmt.Errorf("failed to extract extension from metadata")
		}

	}

	return DownloadResult{
		FilePath:    filePath,
		Url:         url,
		Extension:   extension,
		MimeType:    "video/" + extension,
		Description: desc,
	}, nil
}

func RemoveFile(filePath string) error {
	err1 := os.Remove(filePath)
	err2 := os.Remove(filePath + ".json")
	if err1 != nil && !os.IsNotExist(err1) {
		return fmt.Errorf("failed to remove file %s: %w", filePath, err1)
	}
	if err2 != nil && !os.IsNotExist(err2) {
		return fmt.Errorf("failed to remove file %s.json: %w", filePath, err2)
	}
	return nil
}
