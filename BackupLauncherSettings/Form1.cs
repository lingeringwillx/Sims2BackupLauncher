using IniParser;
using IniParser.Model;
using IniParser.Parser;

namespace BackupLauncherSettings
{
    public partial class backupLauncherSettings : Form
    {

        Settings settings;

        public backupLauncherSettings()
        {
            InitializeComponent();

            try
            {
                settings = new Settings();

                backupFreqNumberBox.Value = settings.backupFreq;
                nBackupsNumberBox.Value = settings.nBackups;
                launcherTextBox.Text = settings.launcherPath;
                saveTextBox.Text = settings.savePath;
                backupTextBox.Text = settings.backupPath;
            }
            catch (Exception ex)
            {
                this.saveButton.Enabled = false;
                displayException(ex);
            }
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

            settings.exceptions = [];
            for (int i = 0; i < hoodsBox.Items.Count; i++)
            {
                if (!hoodsBox.GetItemChecked(i))
                {
                    settings.exceptions.Add(hoodsBox.Items[i].ToString());
                }
            }

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

            try
            {
                settings.Write();
            }
            catch (Exception ex)
            {
                displayException(ex);
            }

            MessageBox.Show("The changes were saved!", "Info", MessageBoxButtons.OK);
        }

        public void displayException(Exception exception)
        {
            MessageBox.Show(exception.ToString(), "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
        }
    }

    internal class Settings
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

        public string settingsPath;
        public IniData data;
        public int backupFreq;
        public int nBackups;
        public List<string> exceptions;
        public string launcherPath;
        public string args;
        public string savePath;
        public string backupPath;

        public Settings()
        {
            string AppDataPath = Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData);
            this.settingsPath = Path.Join(AppDataPath, "BackupLauncher", "settings.ini");

            if (File.Exists(this.settingsPath))
            {
                FileIniDataParser parser = new FileIniDataParser();
                data = parser.ReadFile(this.settingsPath);
            }
            else
            {
                IniDataParser parser = new IniDataParser();
                data = parser.Parse(settingsIni);
            }

            this.backupFreq = int.Parse(data["BackupSettings"]["BackupFrequency"]);
            this.nBackups = int.Parse(data["BackupSettings"]["NumberOfBackups"]);
            this.exceptions = data["BackupSettings"]["Exceptions"].Split(",").ToList();
            this.launcherPath = data["Paths"]["LauncherPath"];
            this.args = data["Paths"]["Arguments"];
            this.savePath = data["Paths"]["SavePath"];
            this.backupPath = data["Paths"]["BackupPath"];
        }

        public void Write()
        {
            this.data["BackupSettings"]["BackupFrequency"] = this.backupFreq.ToString();
            this.data["BackupSettings"]["NumberOfBackups"] = this.nBackups.ToString();
            this.data["BackupSettings"]["Exceptions"] = String.Join(",", this.exceptions);
            this.data["Paths"]["LauncherPath"] = this.launcherPath;

            if (Path.GetFileName(this.savePath) == "Neighborhoods")
            {
                this.savePath = Path.GetDirectoryName(this.savePath);
            }

            this.data["Paths"]["SavePath"] = this.savePath;
            this.data["Paths"]["BackupPath"] = this.backupPath;

            FileIniDataParser parser = new FileIniDataParser();

            if (!Directory.Exists(Path.GetDirectoryName(this.settingsPath)))
            {
                Directory.CreateDirectory(Path.GetDirectoryName(this.settingsPath));
            }

            parser.WriteFile(this.settingsPath, this.data);
        }
    }
}
