using IniParser;
using IniParser.Model;
using System.Diagnostics;
using System.IO.Compression;
using System.Collections.Immutable;
using System.Globalization;

namespace BackupLauncher
{
    internal class Settings
    {
        public int backupFreq;
        public int nBackups;
        public string[] exceptions;
        public string launcherPath;
        public string args;
        public string savePath;
        public string backupPath;

        public Settings(string path)
        {
            FileIniDataParser parser = new FileIniDataParser();
            IniData data = parser.ReadFile(path);

            this.backupFreq = int.Parse(data["BackupSettings"]["BackupFrequency"]);
            this.nBackups = int.Parse(data["BackupSettings"]["NumberOfBackups"]);
            this.exceptions = data["BackupSettings"]["Exceptions"].Split(",");
            this.launcherPath = data["Paths"]["LauncherPath"];
            this.args = data["Paths"]["Arguments"];
            this.savePath = Path.Join(data["Paths"]["SavePath"], "Neighborhoods");
            this.backupPath = data["Paths"]["BackupPath"];
        }
    }
    internal class Program
    {
        static void Main(string[] args)
        {
            string AppDataPath = Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData);
            string settingsPath = Path.Join(AppDataPath, "BackupLauncher", "settings.ini");

            if (!File.Exists(settingsPath)) {
                ProcessStartInfo proc = new ProcessStartInfo("BackupLauncherSettings.exe", "");
                proc.UseShellExecute = true;
                Process.Start(proc);
                Environment.Exit(1);
            }

            Settings settings = new Settings(settingsPath);

            if(!Directory.Exists(settings.savePath))
            {
                Console.WriteLine("Save folder " + settings.launcherPath + "was not found");
                Console.ReadKey();
                Environment.Exit(1);
            }

            if (settings.backupPath == "")
            {
                string documentsPath = Environment.GetFolderPath(Environment.SpecialFolder.Personal);
                settings.backupPath = Path.Join(documentsPath, "Sims 2 Backups");

                if (!Directory.Exists(settings.backupPath))
                {
                    Directory.CreateDirectory(settings.backupPath);
                }
            }

            if (!Directory.Exists(settings.backupPath))
            {
                Console.WriteLine("Backup folder " + settings.launcherPath + "was not found");
                Console.ReadKey();
                Environment.Exit(1);
            }

            bool ok = CreateBackups(settings);

            if (!ok)
            {
                Console.ReadKey();
                Environment.Exit(1);
            }

            if (!File.Exists(settings.launcherPath))
            {
                Console.WriteLine("The launcher " + settings.launcherPath + "was not found");
                Console.ReadKey();
                Environment.Exit(1);
            }

            ok = LaunchGame(settings.launcherPath, settings.args);

            if (!ok)
            {
                Console.ReadKey();
                Environment.Exit(1);
            }
        }

        static bool CreateBackups(Settings settings)
        {
            //filter exceptions
            List<string> hoodPaths = Directory.GetDirectories(settings.savePath).ToList();
            hoodPaths = hoodPaths.Where(hoodPath => !settings.exceptions.Contains(Path.GetFileName(hoodPath))).ToList();

            //backup neighborhoods
            foreach(string hoodPath in hoodPaths)
            {
                string backupPath = Path.Join(settings.backupPath, Path.GetFileName(hoodPath));

                if (!Directory.Exists(backupPath))
                {
                    Directory.CreateDirectory(backupPath);
                }

                //get and sort backups
                List<string> backups = Directory.GetFiles(backupPath).ToList();
                backups.Sort();

                //create backup
                string newBackupPath = Path.Join(backupPath, DateTime.Now.ToString("yyyy-MM-dd")) + ".zip";

                DateTime lastBackupDate;

                if(backups.Count > 0)
                {
                    lastBackupDate = DateTime.ParseExact(Path.GetFileName(backups.Last()).Substring(0,10), "yyyy-MM-dd", CultureInfo.InvariantCulture);
                }
                else
                {
                    lastBackupDate = DateTime.MinValue;
                }
                    
                if ((DateTime.Now - lastBackupDate).Days >= settings.backupFreq)
                {
                    Console.WriteLine("Creating backup for " + Path.GetFileName(hoodPath) + "...");
                    bool ok = CreateBackup(hoodPath, newBackupPath);

                    if (ok)
                    {
                        backups.Add(newBackupPath);
                    }
                    else
                    {
                        Console.WriteLine("Failed to create backup");
                        return false;
                    }
                }

                //delete old backups
                for (int i = 0; i <  backups.Count - settings.nBackups; i++)
                {
                    File.Delete(backups[i]);
                }

            }

            Console.WriteLine("");
            return true;
        }

        static bool CreateBackup(string folderPath, string backupPath)
        {
            try
            {
                ZipFile.CreateFromDirectory(folderPath, backupPath, CompressionLevel.Fastest, true);
                return true;
            }
            catch(IOException)
            {
                return false;
            }
        }

        static bool LaunchGame(string launcherPath, string args)
        {
            Console.WriteLine("Starting the game...");
            Thread.Sleep(1000);

            try
            {
                ProcessStartInfo proc = new ProcessStartInfo(launcherPath, args);
                proc.UseShellExecute = true;
                Process.Start(proc);
                return true;
            }
            catch
            {
                Console.WriteLine("Failed to launch the game");
                return false;
            }
        }
    }
}
