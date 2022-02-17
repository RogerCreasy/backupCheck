package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var config Config

type Config struct {
	Webhook string `json:"webhook"`
}

var backupFiles []BackupFile

// data struct for config files and servers to check json
type BackupFile struct {
	Server string `json:"server"`
	Path   string `json:"path"`
}

// populate config var from the system configuration
func setConfiguration() {
	//import the list of files to check
	files, err := ioutil.ReadFile("/etc/backupCheck/backupFiles.json")

	if err != nil {
		log.Fatal("Error opening files config file ", err)
	}

	err = json.Unmarshal(files, &backupFiles)
	if err != nil {
		log.Fatal("Error during files Unmarshal(): ", err)
	}

	// import the config file
	fileContent, err := ioutil.ReadFile("/etc/backupCheck/config.json")
	if err != nil {
		log.Fatal("Error opening files config file ", err)
	}

	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		log.Fatal("Error during config Unmarshal(): ", err)
	}
}

func sendMessage(server string, fileName string, modDate string, warning string) {
	message := map[string]string{"text": "*Warning:* " + warning + "\n" + server + "\n" + fileName + "\n" + modDate}
	payload, err := json.Marshal(message)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(config.Webhook, "application/json",
		bytes.NewBuffer(payload))

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	if res["json"] != nil {
		log.Fatal(res["json"])
	}
}

func findMostRecentFile(files []string) (MostRecent string) {
	iteration := 0
	var ModTime time.Time
	for _, filePath := range files {
		file, err := os.Stat(filePath)
		if err != nil {
			log.Fatal(err)
		}
		if iteration == 0 {
			ModTime = file.ModTime()
			MostRecent = filePath
		}
		if file.ModTime().After(ModTime) {
			ModTime = file.ModTime()
			MostRecent = filePath
		}

		iteration++
	}

	return MostRecent
}

func main() {
	setConfiguration()
	Path := ""
	for _, file := range backupFiles {
		if strings.Contains(file.Path, "*") {
			files, err := filepath.Glob(file.Path)
			if err != nil {
				log.Fatal(err)
			}

			Path = findMostRecentFile(files)
		} else {
			Path = file.Path
		}

		fileName := "File: " + filepath.Base(Path)
		server := "Server: " + file.Server

		if fileStat, err := os.Stat(Path); err == nil {
			now := time.Now()
			comparisonTime := now.AddDate(0, 0, -7)
			fileMod := fileStat.ModTime()
			modDate := "Last Modified: " + fileStat.ModTime().String()
			warning := "Backup older than 7 days"
			old := fileMod.Before(comparisonTime) // Is the mod time more than 7 days ago?

			if old {
				sendMessage(server, fileName, modDate, warning)
			}
		} else if os.IsNotExist(err) {
			fileName := "File: " + filepath.Base(file.Path)
			modDate := "Last Modified: N/A"
			warning := "file not found"
			sendMessage(server, fileName, modDate, warning)
		} else if err != nil {
			log.Fatal(err)
		}
	}
}
