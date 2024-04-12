## Sims 2 Backup Launcher

This a simple script that backups your neighborhoods before starting the game.

**Requires:** Windows 10 or higher.

**Instructions:**

1- Copy the contents of the **TSBin** to **Fun with Pets\SP9\TSBin** (or Mansions and Gardens for the CD version).

2- Copy the **Sims 2 Backups** folder to **Documents\EA Games**.

3- Go to **Documents\EA Games\Sims 2 Backups** and open the **settings.txt** file and configure it to your liking:

- **backup_freq** is how often your neighborhoods would be backed in days.

- **number_of_backups** is how many backups to keep (older backups will be deleted).

- **exceptions** is a list of the neighborhoods that you don't want to backup, seperated by commas.

4- Launch **Sims2Backup.exe**. It will backup your neighborhoods and then launch *launcher_path* if it's specified in **settings.txt**, otherwise it will use the Sims2RPC launcher if you have it, or the normal game executable if you don't.

**Note:** Don't rename the backups or their folder. The script depends on their names to figure out which one is the oldest.
