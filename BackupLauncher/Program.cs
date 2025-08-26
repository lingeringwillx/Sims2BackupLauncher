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
                try
                {
                    ProcessStartInfo process = new ProcessStartInfo("BackupLauncherSettings.exe", "");
                    process.UseShellExecute = true;
                    Process.Start(process);
                    Environment.Exit(1);
                }
                catch (Exception ex)
                {
                    Exit("Config file was not found.\nFailed to start BackupLauncherSettings.exe", ex);
                }
            }

            Settings settings = null;

            try
            {
                settings = new Settings(settingsPath);
            }
            catch (Exception ex)
            {
                Exit("Failed to read settings.ini", ex);
            }

            if(!Directory.Exists(settings.savePath))
            {
                Exit("Save folder " + settings.savePath + "was not found");
            }

            if (settings.backupPath == "")
            {
                string documentsPath = Environment.GetFolderPath(Environment.SpecialFolder.Personal);
                settings.backupPath = Path.Join(documentsPath, "Sims 2 Backups");

                if (!Directory.Exists(settings.backupPath))
                {
                    try
                    {
                        Directory.CreateDirectory(settings.backupPath);
                    }
                    catch (Exception ex)
                    {
                        Exit("Failed to create " + settings.backupPath, ex);
                    }
                }
            }

            if (!Directory.Exists(settings.backupPath))
            {
                Exit("Backup folder " + settings.backupPath + "was not found");
            }

            CreateBackups(settings);

            if (!File.Exists(settings.launcherPath))
            {
                Exit("The launcher " + settings.launcherPath + "was not found");
            }

            LaunchGame(settings.launcherPath, settings.args);
        }

        static void CreateBackups(Settings settings)
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
                    try
                    {
                        Directory.CreateDirectory(backupPath);
                    }
                    catch (Exception ex)
                    {
                        Exit("Failed to create " + backupPath, ex);
                    }
                }

                //check if a path is a valid backup
                bool isBackup(string path)
                {
                    DateTime dateValue;
                    string fileName = Path.GetFileName(path);
                    string dateStr = fileName.Substring(0, 10);
                    string fmt = "yyyy-MM-dd";
                    CultureInfo culture = CultureInfo.InvariantCulture;
                    DateTimeStyles style = DateTimeStyles.None;
                    
                    bool isDate = DateTime.TryParseExact(dateStr, fmt, culture, style, out dateValue);
                    bool isZip = fileName.EndsWith(".zip");

                    return isDate && isZip;
                }

                //get and sort backups
                List<string> backups = Directory.GetFiles(backupPath).ToList();
                backups = backups.Where(backup => isBackup(backup)).ToList();
                backups.Sort();

                //create backup
                DateTime lastBackupDate = DateTime.MinValue;

                if (backups.Count > 0)
                {
                    lastBackupDate = DateTime.ParseExact(Path.GetFileName(backups.Last()).Substring(0,10), "yyyy-MM-dd", CultureInfo.InvariantCulture);
                }
                    
                if ((DateTime.Now - lastBackupDate).Days >= settings.backupFreq)
                {
                    Console.WriteLine("Creating backup for " + Path.GetFileName(hoodPath) + "...");
                    string newBackupPath = Path.Join(backupPath, DateTime.Now.ToString("yyyy-MM-dd")) + ".zip";
                    CreateBackup(hoodPath, newBackupPath);
                    backups.Add(newBackupPath);
                }

                //delete old backups
                for (int i = 0; i <  backups.Count - settings.nBackups; i++)
                {
                    try
                    {
                        File.Delete(backups[i]);
                    }
                    catch (Exception ex) {
                        Exit("Failed to delete " + backups[i], ex);
                    }
                }

            }

            Console.WriteLine("");
        }

        static void CreateBackup(string folderPath, string backupPath)
        {
            try
            {
                ZipFile.CreateFromDirectory(folderPath, backupPath, CompressionLevel.Fastest, true);
            }
            catch (Exception ex)
            {
                Exit("Failed to create backup for " + Path.GetFileName(folderPath), ex);
            }
        }

        static void LaunchGame(string launcherPath, string args)
        {
            Console.WriteLine("Starting the game...");
            Thread.Sleep(1000);

            try
            {
                ProcessStartInfo process = new ProcessStartInfo(launcherPath, args);
                process.UseShellExecute = true;
                Process.Start(process);
            }
            catch (Exception ex)
            {
                Exit("Failed to launch the game", ex);
            }
        }
        static void Exit(string error, Exception? exception = null)
        {
            Console.WriteLine("\n" + error);

            if (exception != null)
            {
                Console.WriteLine("\n" + exception.ToString());
            }

            Console.ReadKey();
            Environment.Exit(1);
        }
    }
}
