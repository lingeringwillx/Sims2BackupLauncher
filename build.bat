go get
go build
mkdir build build\TSBin "build\Sims 2 Backups"
move Sims2Backup.exe build\TSBin\Sims2Backup.exe
copy 7zr.exe build\TSBin\7zr.exe
copy settings.txt "build\Sims 2 Backups\settings.txt"
pause