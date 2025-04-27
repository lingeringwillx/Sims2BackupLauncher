## Sims 2 Backup Launcher

This a simple script that backups your neighborhoods before starting the game.

**Requires:** Windows 10 or higher.

**Instructions:**

1- Copy the contents of the **TSBin** to the folder of the last expansion. Typically **SP9** or **EP9**.

2- Copy the **Sims 2 Backups** folder to the Documents folder.

3- Go to **Documents\Sims 2 Backups** and open the **settings.txt** file and configure it to your liking:

- **backup_freq** is how often your neighborhoods would be backed in days.

- **number_of_backups** is how many backups to keep. Older backups will be deleted.

- **exceptions** is a list of the neighborhoods that you don't want to backup, seperated by commas.

- **launcher_path** is the path to the game's launcher

- **args** is any arguments to pass to the launcher.

- **save_path** is the path to the game's save folder.

- **backup_path** is where the files would be backed up.

4- Launch **Sims2BackupLauncher.exe**. It will backup your neighborhoods and then launch the file in *launcher_path*.

**Note:** Don't rename the backups or their folder. The script depends on their names to figure out which one is the oldest.
