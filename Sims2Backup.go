package main

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"slices"
	"syscall"
	"time"
	
	"golang.org/x/sys/windows/registry"
)

//settings.txt
const SETTINGS string = `#Backup every...
backup_freq = 7 #days

#Number of backups to keep (older backups will be deleted)
number_of_backups = 3

#Path to game launcher (optional)
launcher_path = 

#Neighborhoods to NOT backup (seperate with commas)
exceptions = Tutorial`

type Settings struct {
	freq int
	nBackups int
	launcherPath string
	exceptions []string
}

func main() {
	savePath, backupPath, err := getPaths()
	settings := Settings{}
	
	if err == nil {
		settings, err = parseSettings(backupPath)
		
		if err == nil {
			err = createBackups(savePath, backupPath, settings)
		}
	}
	
	if err != nil {
		fmt.Print("\nPress enter to start the game...")
		fmt.Scanln()
		fmt.Print("\n")
	}
	
	err = launchGame(settings.launcherPath)
	
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Print("Game executable not found...")
		} else {
			fmt.Print(err)
		}
		
		fmt.Scanln()
		fmt.Print("\n")
	}
}

func getPaths() (string, string, error) {
	keys, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\WOW6432Node\\EA GAMES\\The Sims 2", registry.QUERY_VALUE)
	if err != nil {
		fmt.Println("Failed to find game save location")
		return "", "", err
	}
	defer keys.Close()
	
	saveDir, _, err := keys.GetStringValue("DisplayName")
	if err != nil {
		fmt.Println("Failed to find game save location")
		return "", "", err
	}
	
	userDirPath, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to find game save location")
		return "", "", err
	}
	
	savePath := filepath.Join(userDirPath, "My Documents", "EA Games", saveDir, "Neighborhoods")
	backupPath := filepath.Join(userDirPath, "My Documents", "EA Games", "Sims 2 Backups")
	
	return savePath, backupPath, nil
}

func parseSettings(backupPath string) (Settings, error) {
	settings := Settings{7, 3, "", []string{"Tutorial"}}
	
	buf, err := os.ReadFile(filepath.Join(backupPath, "settings.txt"))
	
	if errors.Is(err, fs.ErrNotExist) {
		err := os.Mkdir(backupPath, os.ModeDir)
		if err != nil && !errors.Is(err, fs.ErrExist) {
			fmt.Println("Failed to create the Sims 2 Backups folder")
			return settings, err
		}
		
		err = os.WriteFile(filepath.Join(backupPath, "settings.txt"), []byte(SETTINGS), 666)
		if err != nil {
			fmt.Println("Failed to create settings.txt")
		}
		
		return settings, err
		
	} else if err != nil {
		fmt.Println("Failed to read settings")
		return settings, err
	}
	
	s := string(buf)
	lines := strings.Split(s, "\n")
	
	for _, line := range lines {
		commentIndex := strings.Index(line, "#")
		
		if commentIndex != -1 {
			line = line[:commentIndex]
		}
		
		left, right, found := strings.Cut(line, "=")
		
		if found {
			left = strings.TrimSpace(left)
			right = strings.TrimSpace(right)
			
			if left == "backup_freq" {
				settings.freq, err = strconv.Atoi(right)
				
				if err != nil {
					fmt.Println("Failed to read backup frequency")
					return settings, err
				}
				
			} else if left == "number_of_backups" {
				settings.nBackups, err = strconv.Atoi(right)
				
				if err != nil {
					fmt.Println("Failed to read the number of backups")
					return settings, err
				}
				
			} else if left == "launcher_path" {
				settings.launcherPath = right
			
			} else if left == "exceptions" {
				settings.exceptions = strings.Split(right, ",")
			}
		}
	}
	
	return settings, nil
}

func createBackups(savePath string, backupPath string, settings Settings) error {
	//loop over neighborhoods
	hoodDirs, err := os.ReadDir(savePath)
	if err != nil {
		fmt.Println("Failed to find neighborhoods")
		return err
	}
	
	hoodDirs = filter(hoodDirs, func(hoodDir fs.DirEntry) bool {
		return hoodDir.IsDir() && !slices.Contains(settings.exceptions, hoodDir.Name())
	})
	
	if len(hoodDirs) > 0 {
		err := os.Mkdir(backupPath, os.ModeDir)
		if err != nil && !errors.Is(err, fs.ErrExist) {
			fmt.Println("Could not create the backups folder")
			return err
		}
	}
	
	for _, hoodDir := range hoodDirs {
		hoodSavePath := filepath.Join(savePath, hoodDir.Name())
		hoodBackupPath := filepath.Join(backupPath, hoodDir.Name())
		
		//filter and sort old backups
		err := os.Mkdir(hoodBackupPath, os.ModeDir)
		if err != nil && !errors.Is(err, fs.ErrExist) {
			fmt.Println("Could not create the neighborhood's backup folder")
			return err
		}
		
		backups, err := os.ReadDir(hoodBackupPath)
		if err != nil {
			fmt.Println("Could not access the neighborhood's backups")
			return err
		}
		
		backups = filter(backups, func(backup fs.DirEntry) bool {
			_, err := time.Parse("2006-01-02", backup.Name()[:10])
			return err == nil && !backup.IsDir()
		})
		
		slices.SortFunc(backups, func(a fs.DirEntry, b fs.DirEntry) int {
			return -cmp.Compare(a.Name(), b.Name())
		})
		
		newBackupPath := filepath.Join(hoodBackupPath, time.Now().Format("2006-01-02") + ".7z")
		
		//create first backup
		if len(backups) == 0 {
			fmt.Println("Creating Backup for " + hoodDir.Name() + "...")
			err := createBackup(hoodSavePath, newBackupPath)
			
			if err != nil {
				fmt.Println("Failed to create backup")
				return err
			}
			
		//create a new backup if the last backup is old
		} else {
			lastBackup := backups[0]
			lastBackupDate, err := time.Parse("2006-01-02", lastBackup.Name()[:10])
			if err != nil {
				fmt.Println("Failed to parse backup date")
				return err
			}
			
			if int(time.Since(lastBackupDate).Hours() / 24) >= settings.freq {
				fmt.Println("Creating Backup for " + hoodDir.Name() + "...")
				err = createBackup(hoodSavePath, newBackupPath)
				
				if err != nil {
					fmt.Println("Failed to create backup")
					return err
				}
				
				//delete old backups
				for i := settings.nBackups; i < len(backups); i++ {
					os.Remove(filepath.Join(hoodBackupPath, backups[i].Name()))
				}
			}
		}
	}
	
	return nil
}

func launchGame(launcherPath string) error {
	launchers := []string{}
	
	if launcherPath == "" {
		launchers = []string{"Sims2RPC.exe", "Sims2EP9RPC.exe", "Sims2EP9.exe"}
	} else {
		launchers = []string{launcherPath}
	}
	
	for _, launcher := range launchers {
		fullPath, err := filepath.Abs(strings.Trim(launcher, "\""))
		if(err != nil) {
			return err
		}
		
		if fileExists(fullPath) {
			fmt.Println("Starting Game...\n")
			
			cmd := exec.Command("cmd", "")
			cmd.SysProcAttr = &syscall.SysProcAttr {
				CmdLine: "/c start /b \"\" \"" + fullPath + "\"",
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			
			return nil
		}
	}
	
	return fs.ErrNotExist
}

func createBackup(source string, destination string) error {
	output, err := exec.Command(".\\7zr.exe", "a", "-mx1", destination, source).CombinedOutput()
	
	if err != nil {
		os.Remove(destination)
		fmt.Println(string(output))
	}
	
	return err
}

//generic filter function
func filter[T any](slice []T, f func(T) bool) []T {
    var n []T
    for _, e := range slice {
        if f(e) {
            n = append(n, e)
        }
    }
    return n
}

func fileExists(fileName string) bool {
    _, err := os.Stat(fileName)
    return !errors.Is(err, fs.ErrNotExist)
}