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

#Neighborhoods to NOT backup (seperate with commas)
exceptions = Tutorial

#Path to game launcher (optional)
launcher_path = 
 
#Advanced: optional arguments to be passed to the launcher (seperate with spaces)
args = `

type Settings struct {
	freq int
	nBackups int
	exceptions []string
	launcherPath string
	args string
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
	
	err = launchGame(settings)
	
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println("Game executable not found...")
			
		} else {
			printErr("Could not launch the game's executable...", err, "")
		}
		
		fmt.Scanln()
		fmt.Print("\n")
	}
}

func getPaths() (string, string, error) {
	keys, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\WOW6432Node\\EA GAMES\\The Sims 2", registry.QUERY_VALUE)
	if err != nil {
		printErr("Failed to find game save location", err, "")
		return "", "", err
	}
	defer keys.Close()
	
	saveDir, _, err := keys.GetStringValue("DisplayName")
	if err != nil {
		printErr("Failed to find game save location", err, "")
		return "", "", err
	}
	
	userDirPath, err := os.UserHomeDir()
	if err != nil {
		printErr("Failed to find game save location", err, "")
		return "", "", err
	}
	
	savePath := filepath.Join(userDirPath, "Documents", "EA Games", saveDir, "Neighborhoods")
	backupPath := filepath.Join(userDirPath, "Documents", "EA Games", "Sims 2 Backups")
	
	return savePath, backupPath, nil
}

func parseSettings(backupPath string) (Settings, error) {
	settings := Settings{7, 3, []string{"Tutorial"}, "", ""}
	settingsPath := filepath.Join(backupPath, "settings.txt")
	
	buf, err := os.ReadFile(settingsPath)
	
	if errors.Is(err, fs.ErrNotExist) {
		err := os.Mkdir(backupPath, os.ModeDir)
		if err != nil && !errors.Is(err, fs.ErrExist) {
			printErr("Failed to create the Sims 2 Backups folder", err, backupPath)
			return settings, err
		}
		
		err = os.WriteFile(settingsPath, []byte(SETTINGS), 666)
		if err != nil {
			printErr("Failed to create settings.txt", err, settingsPath)
		}
		
		return settings, err
		
	} else if err != nil {
		printErr("Failed to read settings", err, settingsPath)
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
					printErr("Failed to read backup frequency", err, "")
					return settings, err
					
				} else if settings.freq < 0 {
					printErr("Backup frequency should be a positive number", err, "")
				}
				
			} else if left == "number_of_backups" {
				settings.nBackups, err = strconv.Atoi(right)
				
				if err != nil {
					printErr("Failed to read the number of backups", err, "")
					return settings, err
					
				} else if settings.nBackups <= 0 {
					printErr("Number of backups should be larger than zero", err, "")
				}
				
			} else if left == "exceptions" {
				settings.exceptions = strings.Split(strings.ReplaceAll(right, " ", ""), ",")
				
			} else if left == "launcher_path" {
				settings.launcherPath = right
				
			} else if(left == "args") {
				settings.args = right
			}
		}
	}
	
	return settings, nil
}

func createBackups(savePath string, backupPath string, settings Settings) error {
	//loop over neighborhoods
	hoodDirs, err := os.ReadDir(savePath)
	if err != nil {
		printErr("Failed to find neighborhoods", err, savePath)
		return err
	}
	
	hoodDirs = filter(hoodDirs, func(hoodDir fs.DirEntry) bool {
		return hoodDir.IsDir() && !slices.Contains(settings.exceptions, hoodDir.Name())
	})
	
	if len(hoodDirs) > 0 {
		err := os.Mkdir(backupPath, os.ModeDir)
		if err != nil && !errors.Is(err, fs.ErrExist) {
			printErr("Could not create the backups folder", err, backupPath)
			return err
		}
	}
	
	for _, hoodDir := range hoodDirs {
		hoodSavePath := filepath.Join(savePath, hoodDir.Name())
		hoodBackupPath := filepath.Join(backupPath, hoodDir.Name())
		
		//filter and sort old backups
		err := os.Mkdir(hoodBackupPath, os.ModeDir)
		if err != nil && !errors.Is(err, fs.ErrExist) {
			printErr("Could not create the neighborhood's backup folder", err, hoodBackupPath)
			return err
		}
		
		backups, err := os.ReadDir(hoodBackupPath)
		if err != nil {
			printErr("Could not access the neighborhood's backups", err, hoodBackupPath)
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
				printErr("Failed to create backup", err, newBackupPath)
				return err
			}
			
		//create a new backup if the last backup is old
		} else {
			lastBackup := backups[0]
			lastBackupDate, err := time.Parse("2006-01-02", lastBackup.Name()[:10])
			if err != nil {
				printErr("Failed to parse backup date", err, "")
				return err
			}
			
			if int(time.Since(lastBackupDate).Hours() / 24) >= settings.freq {
				fmt.Println("Creating Backup for " + hoodDir.Name() + "...")
				err = createBackup(hoodSavePath, newBackupPath)
				
				if err != nil {
					printErr("Failed to create backup", err, newBackupPath)
					return err
				}
				
				//delete old backups
				for i := settings.nBackups - 1; i < len(backups); i++ {
					os.Remove(filepath.Join(hoodBackupPath, backups[i].Name()))
				}
			}
		}
	}
	
	return nil
}

func launchGame(settings Settings) error {
	launchers := []string{}
	
	if settings.launcherPath == "" {
		launchers = []string{"Sims2RPC.exe", "Sims2EP9RPC.exe", "Sims2EP9.exe"}
	} else {
		launchers = []string{settings.launcherPath}
	}
	
	for _, launcher := range launchers {
		if fileExists(launcher) {
			fmt.Println("Starting Game...\n")
			
			cmd := exec.Command("cmd", "")
			cmd.SysProcAttr = &syscall.SysProcAttr {
				CmdLine: "/c start /b \"\" \"" + launcher + "\" " + settings.args,
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			
			return err
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

func printErr(info string, err error, path string) {
	fmt.Println(info)
	fmt.Print("error: ")
	fmt.Println(err)
	
	_ = path
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