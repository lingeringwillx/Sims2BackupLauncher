namespace BackupLauncherSettings
{
    partial class backupLauncherSettings
    {
        /// <summary>
        ///  Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        ///  Clean up any resources being used.
        /// </summary>
        /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows Form Designer generated code

        /// <summary>
        ///  Required method for Designer support - do not modify
        ///  the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            backupEveryText = new Label();
            backupFreqNumberBox = new NumericUpDown();
            numberBackupsText = new Label();
            nBackupsNumberBox = new NumericUpDown();
            launcherText = new Label();
            launcherTextBox = new TextBox();
            saveTextBox = new TextBox();
            argsText = new Label();
            backupTextBox = new TextBox();
            saveText = new Label();
            hoodsBox = new CheckedListBox();
            backupText = new Label();
            hoodsText = new Label();
            argsTextBox = new TextBox();
            daysText = new Label();
            launcherBrowseButton = new Button();
            saveBrowseButton = new Button();
            backupBrowseButton = new Button();
            saveButton = new Button();
            ((System.ComponentModel.ISupportInitialize)backupFreqNumberBox).BeginInit();
            ((System.ComponentModel.ISupportInitialize)nBackupsNumberBox).BeginInit();
            SuspendLayout();
            // 
            // backupEveryText
            // 
            backupEveryText.AutoSize = true;
            backupEveryText.Location = new Point(12, 9);
            backupEveryText.Name = "backupEveryText";
            backupEveryText.Size = new Size(136, 21);
            backupEveryText.TabIndex = 0;
            backupEveryText.Text = "Backup Frequency";
            // 
            // backupFreqNumberBox
            // 
            backupFreqNumberBox.Location = new Point(165, 7);
            backupFreqNumberBox.Minimum = new decimal(new int[] { 1, 0, 0, 0 });
            backupFreqNumberBox.Name = "backupFreqNumberBox";
            backupFreqNumberBox.Size = new Size(50, 29);
            backupFreqNumberBox.TabIndex = 1;
            backupFreqNumberBox.Value = new decimal(new int[] { 1, 0, 0, 0 });
            // 
            // numberBackupsText
            // 
            numberBackupsText.AutoSize = true;
            numberBackupsText.Location = new Point(12, 44);
            numberBackupsText.Name = "numberBackupsText";
            numberBackupsText.Size = new Size(147, 21);
            numberBackupsText.TabIndex = 2;
            numberBackupsText.Text = "Number of Backups";
            // 
            // nBackupsNumberBox
            // 
            nBackupsNumberBox.Location = new Point(165, 42);
            nBackupsNumberBox.Minimum = new decimal(new int[] { 1, 0, 0, 0 });
            nBackupsNumberBox.Name = "nBackupsNumberBox";
            nBackupsNumberBox.Size = new Size(50, 29);
            nBackupsNumberBox.TabIndex = 3;
            nBackupsNumberBox.Value = new decimal(new int[] { 1, 0, 0, 0 });
            // 
            // launcherText
            // 
            launcherText.AutoSize = true;
            launcherText.Location = new Point(12, 80);
            launcherText.Name = "launcherText";
            launcherText.Size = new Size(119, 21);
            launcherText.TabIndex = 4;
            launcherText.Text = "Game Launcher";
            // 
            // launcherTextBox
            // 
            launcherTextBox.BackColor = SystemColors.Window;
            launcherTextBox.Location = new Point(137, 77);
            launcherTextBox.Name = "launcherTextBox";
            launcherTextBox.ReadOnly = true;
            launcherTextBox.Size = new Size(500, 29);
            launcherTextBox.TabIndex = 5;
            // 
            // saveTextBox
            // 
            saveTextBox.BackColor = SystemColors.Window;
            saveTextBox.Location = new Point(137, 147);
            saveTextBox.Name = "saveTextBox";
            saveTextBox.ReadOnly = true;
            saveTextBox.Size = new Size(500, 29);
            saveTextBox.TabIndex = 6;
            saveTextBox.TextChanged += saveTextBox_TextChanged;
            // 
            // argsText
            // 
            argsText.AutoSize = true;
            argsText.Location = new Point(12, 115);
            argsText.Name = "argsText";
            argsText.Size = new Size(87, 21);
            argsText.TabIndex = 7;
            argsText.Text = "Arguments";
            // 
            // backupTextBox
            // 
            backupTextBox.BackColor = SystemColors.Window;
            backupTextBox.Location = new Point(137, 182);
            backupTextBox.Name = "backupTextBox";
            backupTextBox.ReadOnly = true;
            backupTextBox.Size = new Size(500, 29);
            backupTextBox.TabIndex = 8;
            // 
            // saveText
            // 
            saveText.AutoSize = true;
            saveText.Location = new Point(12, 150);
            saveText.Name = "saveText";
            saveText.Size = new Size(91, 21);
            saveText.TabIndex = 9;
            saveText.Text = "Save Folder";
            // 
            // hoodsBox
            // 
            hoodsBox.CheckOnClick = true;
            hoodsBox.FormattingEnabled = true;
            hoodsBox.Location = new Point(137, 217);
            hoodsBox.Name = "hoodsBox";
            hoodsBox.Size = new Size(160, 124);
            hoodsBox.TabIndex = 10;
            // 
            // backupText
            // 
            backupText.AutoSize = true;
            backupText.Location = new Point(12, 185);
            backupText.Name = "backupText";
            backupText.Size = new Size(108, 21);
            backupText.TabIndex = 11;
            backupText.Text = "Backup Folder";
            // 
            // hoodsText
            // 
            hoodsText.AutoSize = true;
            hoodsText.Location = new Point(12, 217);
            hoodsText.Name = "hoodsText";
            hoodsText.Size = new Size(119, 21);
            hoodsText.TabIndex = 12;
            hoodsText.Text = "Neighborhoods";
            // 
            // argsTextBox
            // 
            argsTextBox.Location = new Point(137, 112);
            argsTextBox.Name = "argsTextBox";
            argsTextBox.Size = new Size(500, 29);
            argsTextBox.TabIndex = 13;
            // 
            // daysText
            // 
            daysText.AutoSize = true;
            daysText.Location = new Point(221, 9);
            daysText.Name = "daysText";
            daysText.Size = new Size(42, 21);
            daysText.TabIndex = 14;
            daysText.Text = "days";
            // 
            // launcherBrowseButton
            // 
            launcherBrowseButton.Location = new Point(643, 75);
            launcherBrowseButton.Name = "launcherBrowseButton";
            launcherBrowseButton.Size = new Size(100, 31);
            launcherBrowseButton.TabIndex = 15;
            launcherBrowseButton.Text = "Browse";
            launcherBrowseButton.UseVisualStyleBackColor = true;
            launcherBrowseButton.Click += launcherBrowseButton_Click;
            // 
            // saveBrowseButton
            // 
            saveBrowseButton.Location = new Point(643, 145);
            saveBrowseButton.Name = "saveBrowseButton";
            saveBrowseButton.Size = new Size(100, 31);
            saveBrowseButton.TabIndex = 16;
            saveBrowseButton.Text = "Browse";
            saveBrowseButton.UseVisualStyleBackColor = true;
            saveBrowseButton.Click += saveBrowseButton_Click;
            // 
            // backupBrowseButton
            // 
            backupBrowseButton.Location = new Point(643, 182);
            backupBrowseButton.Name = "backupBrowseButton";
            backupBrowseButton.Size = new Size(100, 31);
            backupBrowseButton.TabIndex = 17;
            backupBrowseButton.Text = "Browse";
            backupBrowseButton.UseVisualStyleBackColor = true;
            backupBrowseButton.Click += backupBrowseButton_Click;
            // 
            // saveButton
            // 
            saveButton.Location = new Point(643, 407);
            saveButton.Name = "saveButton";
            saveButton.Size = new Size(100, 31);
            saveButton.TabIndex = 18;
            saveButton.Text = "Save";
            saveButton.UseVisualStyleBackColor = true;
            saveButton.Click += saveButton_Click;
            // 
            // backupLauncherSettings
            // 
            AutoScaleDimensions = new SizeF(9F, 21F);
            AutoScaleMode = AutoScaleMode.Font;
            ClientSize = new Size(757, 452);
            Controls.Add(saveButton);
            Controls.Add(backupBrowseButton);
            Controls.Add(saveBrowseButton);
            Controls.Add(launcherBrowseButton);
            Controls.Add(daysText);
            Controls.Add(argsTextBox);
            Controls.Add(hoodsText);
            Controls.Add(backupText);
            Controls.Add(hoodsBox);
            Controls.Add(saveText);
            Controls.Add(backupTextBox);
            Controls.Add(argsText);
            Controls.Add(saveTextBox);
            Controls.Add(launcherTextBox);
            Controls.Add(launcherText);
            Controls.Add(nBackupsNumberBox);
            Controls.Add(numberBackupsText);
            Controls.Add(backupFreqNumberBox);
            Controls.Add(backupEveryText);
            FormBorderStyle = FormBorderStyle.FixedSingle;
            MaximizeBox = false;
            Name = "backupLauncherSettings";
            StartPosition = FormStartPosition.CenterScreen;
            Text = "Launcher Settings";
            ((System.ComponentModel.ISupportInitialize)backupFreqNumberBox).EndInit();
            ((System.ComponentModel.ISupportInitialize)nBackupsNumberBox).EndInit();
            ResumeLayout(false);
            PerformLayout();
        }

        #endregion

        private Label backupEveryText;
        private NumericUpDown backupFreqNumberBox;
        private Label numberBackupsText;
        private NumericUpDown nBackupsNumberBox;
        private Label launcherText;
        private TextBox launcherTextBox;
        private TextBox saveTextBox;
        private Label argsText;
        private TextBox backupTextBox;
        private Label saveText;
        private CheckedListBox hoodsBox;
        private Label backupText;
        private Label hoodsText;
        private TextBox argsTextBox;
        private Label daysText;
        private Button launcherBrowseButton;
        private Button saveBrowseButton;
        private Button backupBrowseButton;
        private Button saveButton;
    }
}
