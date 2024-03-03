package main

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"strconv"
	"strings"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"time"
	
	"golang.org/x/sys/windows/registry"
)

func main() {
	//section: get save and backup directories
	keys, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\WOW6432Node\\EA GAMES\\The Sims 2", registry.QUERY_VALUE)
	if err != nil {
		backupFailed("Failed to find game save location")
	}
	defer keys.Close()
	
	saveDir, _, err := keys.GetStringValue("DisplayName")
	if err != nil {
		backupFailed("Failed to find game save location")
	}
	
	userDirPath, err := os.UserHomeDir()
	if err != nil {
		backupFailed("Failed to find game save location")
	}
	
	savePath := filepath.Join(userDirPath, "My Documents", "EA Games", saveDir, "Neighborhoods")
	backupPath := filepath.Join(userDirPath, "My Documents", "EA Games", "Sims 2 Backups")
	
	buf, err := os.ReadFile(filepath.Join(backupPath, "settings.txt"))
	if err != nil {
		backupFailed("Failed to read settings")
	}
	
	//section: parse settings
	freq := 7
	nBackups := 3
	var exceptions []string
	
	s := string(buf)
	lines := strings.Split(s, "\n")
	
	for _, line := range lines {
		commentIndex := strings.Index(line, "#")
		
		if commentIndex != -1 {
			line = line[:commentIndex]
		}
		
		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, " ", "")
		
		left, right, found := strings.Cut(line, "=")
		
		if found {
			if left == "backup_freq" {
				freq, err = strconv.Atoi(right)
				
				if err != nil {
					backupFailed("Failed to read backup frequency")
				}
				
			} else if left == "number_of_backups" {
				nBackups, err = strconv.Atoi(right)
				
				if err != nil {
					backupFailed("Failed to read the number of backups")
				}
				
			} else if left == "exceptions" {
				exceptions = strings.Split(right, ",")
			}
		}
	}
	
	//section: loop over neighborhoods
	hoodDirs, err := os.ReadDir(savePath)
	if err != nil {
		backupFailed("Failed to find neighborhoods")
	}
	
	if len(hoodDirs) > 0 {
		err := os.Mkdir(backupPath, os.ModeDir)
		if err != nil && !errors.Is(err, fs.ErrExist) {
			backupFailed("Could not create backups folder")
		}
	}
	
	for _, hoodDir := range hoodDirs {
		hoodSavePath := filepath.Join(savePath, hoodDir.Name())
		hoodBackupPath := filepath.Join(backupPath, hoodDir.Name())
		
		if hoodDir.IsDir() && !slices.Contains(exceptions, hoodDir.Name()){
			//section: filter and sort old backups
			err := os.Mkdir(hoodBackupPath, os.ModeDir)
			if err != nil && !errors.Is(err, fs.ErrExist) {
				backupFailed("Could not create the neighborhood's backup folder")
			}
			
			backups, err := os.ReadDir(hoodBackupPath)
			if err != nil {
				backupFailed("Could not access the neighborhood's backups")
			}
			
			var filteredBackups []fs.DirEntry
			
			for _, backup := range backups {
				_, err := time.Parse("2006-01-02", backup.Name()[:10])
				if err == nil && backup.Name()[10:] == ".7z" {
					filteredBackups = append(filteredBackups, backup)
				}
			}
			
			backups = filteredBackups
			
			slices.SortFunc(backups, func(a fs.DirEntry, b fs.DirEntry) int {
				return cmp.Compare(a.Name(), b.Name())
			})
			
			slices.Reverse(backups)
			
			//section: create backup
			newBackupPath := filepath.Join(hoodBackupPath, time.Now().Format("2006-01-02") + ".7z")
			
			//create first backup
			if len(backups) == 0 {
				fmt.Println("Creating Backup for " + hoodDir.Name() + "...")
				err := createBackup(hoodSavePath, newBackupPath)
				
				if err != nil {
					backupFailed("Failed to create backup")
				}
				
			//create new backup if the last backup is old
			} else {
				lastBackup := backups[0]
				lastBackupDate, err := time.Parse("2006-01-02", lastBackup.Name()[:10])
				if err != nil {
					backupFailed("Failed to parse backup date")
				}
				
				neighborhoodPackage := filepath.Join(hoodSavePath, hoodDir.Name() + "_Neighborhood.package")
				info, err := os.Stat(neighborhoodPackage)
				if err != nil {
					backupFailed("Could not retrieve neighborhood info")
				}
				
				if int(time.Since(lastBackupDate).Hours() / 24) < freq || info.ModTime().Sub(lastBackupDate).Hours() < 24 {
					continue
				}
				
				fmt.Println("Creating Backup for " + hoodDir.Name() + "...")
				err = createBackup(hoodSavePath, newBackupPath)
				
				if err != nil {
					backupFailed("Failed to create backup")
				}
				
				//section: delete old backups
				for i := nBackups; i < len(backups); i++ {
					os.Remove(filepath.Join(hoodBackupPath, backups[i].Name()))
				}
			}
		}
	}
	
	//section: launch game
	launchGame()
}

func backupFailed(err string) {
	fmt.Println(err)
	fmt.Print("\nPress enter to start the game...")
	fmt.Scanln()
	fmt.Print("\n")
	launchGame()
}

func launchGame() {
	launcherNames := [3]string{"Sims2RPC.exe", "Sims2EP9RPC.exe", "Sims2EP9.exe"}
	
	for _, name := range launcherNames {
		if fileExists(name) {
			fullPath, err := filepath.Abs(name)
			if(err != nil) {
				fmt.Println(err)
			}
			
			fmt.Println("Starting Game...\n")
			
			cmd := exec.Command(fullPath, "")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start()
			
			os.Exit(0)
		}
	}
	
	fmt.Print("Game executable not found")
	fmt.Scanln()
	fmt.Print("\n")
	os.Exit(1)
}

func createBackup(source string, destination string) error {
	output, err := exec.Command("powershell", ".//7zr.exe", "a", "\"" + destination + "\"", "\"" + source + "\"").CombinedOutput()
	
	if err != nil {
		os.Remove(destination)
		fmt.Println(string(output))
	}
	
	return err
}

func fileExists(fileName string) bool {
    _, err := os.Stat(fileName)
    return !errors.Is(err, fs.ErrNotExist)
}