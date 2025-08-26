using IniParser;
using IniParser.Model;

namespace BackupLauncherSettings
{
    public partial class backupLauncherSettings : Form
    {

        const string settingsIni = @"[BackupSettings]
BackupFrequency = 5

NumberOfBackups = 3

Exceptions = Tutorial

[Paths]
LauncherPath =

Arguments =

SavePath =

BackupPath =";

        Settings settings;

        public backupLauncherSettings()
        {
            InitializeComponent();

            string AppDataPath = Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData);
            string settingsPath = Path.Join(AppDataPath, "BackupLauncher", "settings.ini");

            if (!Directory.Exists(Path.GetDirectoryName(settingsPath)))
            {
                Directory.CreateDirectory(Path.GetDirectoryName(settingsPath));
            }

            if (!File.Exists(settingsPath))
            {
                File.WriteAllText(settingsPath, settingsIni);
            }

            settings = new Settings(settingsPath);

            backupFreqNumberBox.Value = settings.backupFreq;
            nBackupsNumberBox.Value = settings.nBackups;
            launcherTextBox.Text = settings.launcherPath;
            saveTextBox.Text = settings.savePath;
            backupTextBox.Text = settings.backupPath;
        }

        private void saveTextBox_TextChanged(object sender, EventArgs e)
        {
            hoodsBox.Items.Clear();
            string hoodsPath = saveTextBox.Text;

            if (Path.GetFileName(hoodsPath) != "Neighborhoods")
            {
                hoodsPath = Path.Join(hoodsPath, "Neighborhoods");
            }

            if (Directory.Exists(hoodsPath))
            {
                foreach (string hoodPath in Directory.GetDirectories(hoodsPath))
                {
                    string hood = Path.GetFileName(hoodPath);
                    hoodsBox.Items.Add(hood, !settings.exceptions.Contains(hood));
                }
            }
        }

        private void launcherBrowseButton_Click(object sender, EventArgs e)
        {
            OpenFileDialog openFileDialog = new OpenFileDialog();
            if (openFileDialog.ShowDialog() == DialogResult.OK)
            {
                launcherTextBox.Text = openFileDialog.FileName;
            }
        }

        private void saveBrowseButton_Click(object sender, EventArgs e)
        {
            FolderBrowserDialog folderBrowserDialog = new FolderBrowserDialog();
            if (folderBrowserDialog.ShowDialog() == DialogResult.OK)
            {
                saveTextBox.Text = folderBrowserDialog.SelectedPath;
            }
        }

        private void backupBrowseButton_Click(object sender, EventArgs e)
        {
            FolderBrowserDialog folderBrowserDialog = new FolderBrowserDialog();
            if (folderBrowserDialog.ShowDialog() == DialogResult.OK)
            {
                backupTextBox.Text = folderBrowserDialog.SelectedPath;
            }
        }

        private void saveButton_Click(object sender, EventArgs e)
        {
            settings.backupFreq = (int)backupFreqNumberBox.Value;
            settings.nBackups = (int)nBackupsNumberBox.Value;

            List<string> exceptions = [];
            for (int i = 0; i < hoodsBox.Items.Count; i++)
            {
                if (!hoodsBox.GetItemChecked(i))
                {
                    exceptions.Add(hoodsBox.Items[i].ToString());
                }
            }
            
            settings.exceptions = exceptions.ToArray();

            if (launcherTextBox.Text == "")
            {
                MessageBox.Show("The launcher was not specified", "Error", MessageBoxButtons.OK, MessageBoxIcon.Information);
                return;
            }

            if (saveTextBox.Text == "")
            {
                MessageBox.Show("The save folder was not specified", "Error", MessageBoxButtons.OK, MessageBoxIcon.Information);
                return;
            }

            settings.launcherPath = launcherTextBox.Text;
            settings.savePath = saveTextBox.Text;
            settings.backupPath = backupTextBox.Text;

            settings.Write();

            MessageBox.Show("The changes were saved!", "Info", MessageBoxButtons.OK);
        }
    }

    internal class Settings
    {
        public string settingsPath;
        public IniData data;
        public int backupFreq;
        public int nBackups;
        public string[] exceptions;
        public string launcherPath;
        public string args;
        public string savePath;
        public string backupPath;

        public Settings(string path)
        {
            this.settingsPath = path;
            FileIniDataParser parser = new FileIniDataParser();
            data = parser.ReadFile(path);

            this.backupFreq = int.Parse(data["BackupSettings"]["BackupFrequency"]);
            this.nBackups = int.Parse(data["BackupSettings"]["NumberOfBackups"]);
            this.exceptions = data["BackupSettings"]["Exceptions"].Split(",");
            this.launcherPath = data["Paths"]["LauncherPath"];
            this.args = data["Paths"]["Arguments"];
            this.savePath = Path.Join(data["Paths"]["SavePath"]);
            this.backupPath = data["Paths"]["BackupPath"];
        }

        public void Write()
        {
            this.data["BackupSettings"]["BackupFrequency"] = this.backupFreq.ToString();
            this.data["BackupSettings"]["NumberOfBackups"] = this.nBackups.ToString();
            this.data["BackupSettings"]["Exceptions"] = String.Join(",", this.exceptions);
            this.data["Paths"]["LauncherPath"] = this.launcherPath;

            this.data["Paths"]["SavePath"] = this.savePath;

            if (Path.GetFileName(this.savePath) == "Neighborhoods")
            {
                this.savePath = Path.GetDirectoryName(this.savePath);
            }

            this.data["Paths"]["BackupPath"] = this.backupPath;

            FileIniDataParser parser = new FileIniDataParser();
            parser.WriteFile(this.settingsPath, this.data);
        }
    }
}
