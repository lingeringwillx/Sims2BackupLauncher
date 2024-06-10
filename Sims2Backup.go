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
args = 

#Optional game save path (The Sims 2 folder), if the launcher is unable to find your save on its own
save_path = 

#Optional backup path if you want the backup to be saved in a different location
backup_path = `

type Settings struct {
	freq int
	nBackups int
	exceptions []string
	launcherPath string
	args string
	savePath string
	backupPath string
	altBackupPath string
}

func main() {
	documentsPath, err := getDocumentsPath()
	settings := Settings{}
	var savePath, backupPath string
	var err2 error
	
	if err == nil {
		settings, err = parseSettings(documentsPath)
	}
	
	if err == nil {
		savePath, backupPath, err = getPaths(documentsPath, settings)
	}
	
	if err == nil {
		err = createBackups(savePath, backupPath, settings)
		
		if settings.altBackupPath != "" {
			err2 = createBackups(savePath, settings.altBackupPath, settings)
		}
	}
	
	time.Sleep(time.Second)
	
	if err != nil && err2 != nil {
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

func getDocumentsPath() (string, error) {
	keys, err := registry.OpenKey(registry.CURRENT_USER, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\User Shell Folders", registry.QUERY_VALUE)
	if err != nil {
		printErr("Failed to find documents folder location", err, "")
		return "", err
	}
	defer keys.Close()
	
	documentsPath, _, err := keys.GetStringValue("Personal")
	if err != nil {
		printErr("Failed to find documents folder location", err, "")
		return "", err
	}
	
	//resolve enviroment variables in path
	splitDocumentsPath := strings.Split(documentsPath, string(filepath.Separator))
	
	for i, part := range splitDocumentsPath {
		if len(part) > 2 && strings.HasPrefix(part, "%") && strings.HasSuffix(part, "%") {
			resolvedPath, found := os.LookupEnv(part[1:len(part) - 1])
			
			if found {
				splitDocumentsPath[i] = resolvedPath
			} else {
				fmt.Println("error: Failed to resolve enviroment variable in path", documentsPath)
				return "", errors.New("")
			}
		}
	}
	
	splitDocumentsPath[0] += string(filepath.Separator)
	documentsPath = filepath.Join(splitDocumentsPath...)
	return documentsPath, nil
}

func getPaths(documentsPath string, settings Settings) (string, string, error) {
	savePath := ""
	backupPath := ""
	
	if settings.savePath == "" {
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
		
		savePath = filepath.Join(documentsPath, "EA Games", saveDir, "Neighborhoods")
		
	} else {
		savePath = filepath.Join(settings.savePath, "Neighborhoods")
	}
	
	if settings.backupPath == "" {
		backupPath = filepath.Join(documentsPath, "EA Games", "Sims 2 Backups")
	} else {
		backupPath = settings.backupPath
	}
	
	return savePath, backupPath, nil
}

func parseSettings(documentsPath string) (Settings, error) {
	settings := Settings{7, 3, []string{"Tutorial"}, "", "", "", "", ""}
	settingsFolderPath := filepath.Join(documentsPath, "EA Games", "Sims 2 Backups")
	settingsPath := filepath.Join(settingsFolderPath, "settings.txt")
	
	buf, err := os.ReadFile(settingsPath)
	
	if errors.Is(err, fs.ErrNotExist) {
		err := os.Mkdir(settingsFolderPath, os.ModeDir)
		if err != nil && !errors.Is(err, fs.ErrExist) {
			printErr("Failed to create the Sims 2 Backups folder", err, settingsFolderPath)
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
				settings.exceptions = strings.Split(right, ",")
				
				for i, hood := range settings.exceptions {
					settings.exceptions[i] = strings.TrimSpace(hood)
				}
				
			} else if left == "launcher_path" {
				settings.launcherPath = right
				
			} else if left == "args" {
				settings.args = right
				
			} else if left == "save_path" {
				settings.savePath = right
				
			} else if left == "backup_path" {
				settings.backupPath = right
				
			} else if left == "alt_backup_path" {
				settings.altBackupPath = right
			}
		}
	}
	
	return settings, nil
}

func createBackups(savePath string, backupPath string, settings Settings) error {
	//loop over neighborhoods
	hoodDirs, err := os.ReadDir(savePath)
	if err != nil {
		printErr("Failed to find the neighborhoods", err, savePath)
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
		
		fmt.Println("Backup path:", backupPath)
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
	
	fmt.Println()
	
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
