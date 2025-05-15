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

type Settings struct {
    freq int
    nBackups int
    exceptions []string
    launcherPath string
    args string
    savePath string
    backupPath string
}

func main() {
    documentsPath, err := getDocumentsPath()

    exitIfErr(err)

    settings := Settings{}
    settings, err = parseSettings(documentsPath)

    exitIfErr(err)

    err = createBackups(settings)

    exitIfErr(err)

    time.Sleep(time.Second)

    err = launchGame(settings)

    exitIfErr(err)
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

func parseSettings(documentsPath string) (Settings, error) {
    settings := Settings{}

    settingsPath := filepath.Join(documentsPath, "Sims 2 Backups", "settings.txt")
    buf, err := os.ReadFile(settingsPath)

    if err != nil {
        printErr("Failed to read the program's settings", err, settingsPath)
        return settings, err
    }

    s := string(buf)
    lines := strings.Split(s, "\n")

    for _, line := range lines {
        if strings.HasPrefix(line, "#") {
            continue
        }

        left, right, found := strings.Cut(line, "=")

        if found {
            left = strings.TrimSpace(left)
            right = strings.TrimSpace(right)

            if left == "BackupFrequency" {
                settings.freq, err = strconv.Atoi(right)

                if err != nil {
                    printErr("Failed to read backup frequency", err, "")
                    return settings, err

                } else if settings.freq < 0 {
                    printErr("Backup frequency should be a positive number", err, "")
                }

            } else if left == "NumberOfBackups" {
                settings.nBackups, err = strconv.Atoi(right)

                if err != nil {
                    printErr("Failed to read the number of backups", err, "")
                    return settings, err

                } else if settings.nBackups <= 0 {
                    printErr("Number of backups should be larger than zero", err, "")
                }

            } else if left == "Exceptions" {
                settings.exceptions = strings.Split(right, ",")

                for i, hood := range settings.exceptions {
                    settings.exceptions[i] = strings.TrimSpace(hood)
                }

            } else if left == "LauncherPath" {
                if right == "" {
                    fmt.Println("error: Launcher path is not specified in settings.txt")
                    fmt.Println(right)
                    return settings, errors.New("")

                } else if !fileExists(right) {
                    fmt.Println("error: The launcher specified in settings.txt was not found")
                    fmt.Println(right)
                    return settings, errors.New("")

                } else {
                    settings.launcherPath = right
                }

            } else if left == "Arguments" {
                settings.args = right

            } else if left == "SavePath" {
                if right == "" {
                    fmt.Println("error: The game's save location is not listed in settings.txt")
                    fmt.Println(right)
                    return settings, errors.New("")

                } else if !fileExists(right) {
                    fmt.Println("error: Save path was not found")
                    fmt.Println(right)
                    return settings, errors.New("")

                } else {
                    settings.savePath = filepath.Join(right, "Neighborhoods")
                }

            } else if left == "BackupPath" {
                if right == "" {
                    settings.backupPath = filepath.Join(documentsPath, "Sims 2 Backups")

                } else if !fileExists(right) {
                    fmt.Println("error: Backup path was not found")
                    fmt.Println(right)
                    return settings, errors.New("")

                } else {
                    settings.backupPath = right
                }
            }
        }
    }

    return settings, nil
}

func createBackups(settings Settings) error {
    //loop over neighborhoods
    hoodDirs, err := os.ReadDir(settings.savePath)
    if err != nil {
        printErr("Failed to find the game's neighborhoods", err, settings.savePath)
        return err
    }

    hoodDirs = filter(hoodDirs, func(hoodDir fs.DirEntry) bool {
        return hoodDir.IsDir() && !slices.Contains(settings.exceptions, hoodDir.Name())
    })

    if len(hoodDirs) > 0 {
        fmt.Println("Backup path:", settings.backupPath)
    }

    for _, hoodDir := range hoodDirs {
        hoodSavePath := filepath.Join(settings.savePath, hoodDir.Name())
        hoodBackupPath := filepath.Join(settings.backupPath, hoodDir.Name())

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

        newBackupPath := filepath.Join(hoodBackupPath, time.Now().Format("2006-01-02") + ".zip")

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
    fmt.Println("Starting Game...\n")

    //this messy command makes the cmd behave as expected for some reason
    cmd := exec.Command("cmd", "")
    cmd.SysProcAttr = &syscall.SysProcAttr {
        CmdLine: "/c start /b \"\" \"" + settings.launcherPath + "\" " + settings.args,
    }
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()

    if err != nil {
        printErr("Could not launch the game's executable...", err, settings.launcherPath)
    }

    return err
}

func createBackup(source string, destination string) error {
    output, err := exec.Command("./7za.exe", "a", destination, source).CombinedOutput()

    if err != nil {
        os.Remove(destination)
        fmt.Println(string(output))
    }

    return err
}

func printErr(info string, err error, path string) {
    fmt.Print("error: ")
    fmt.Println(info)
    fmt.Println(err)
    fmt.Println(path)
}

func exitIfErr(err error) {
    if err != nil {
        fmt.Scanln()
        fmt.Print("\n")
        os.Exit(1)
    }
}

func fileExists(fileName string) bool {
    _, err := os.Stat(fileName)
    return !errors.Is(err, fs.ErrNotExist)
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