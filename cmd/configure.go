// فایل: cmd/configure.go
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// --- متغیرهای فلگ‌ها برای دستور اصلی configure ---
var (
	mongoRemoteURI string
	mongoLocalURI  string
	backupPath     string
	mongoToolsPath string
)

// --- دستور اصلی configure ---
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure DataWeaver CLI settings.",
	Long: `This command allows you to configure various settings for the DataWeaver CLI,
such as database connection details, backup paths, etc.
Run without flags for an interactive setup.`,
	Run: func(cmd *cobra.Command, args []string) {
		// این تابع مقادیر را از فلگ‌ها یا به صورت تعاملی می‌گیرد و ذخیره می‌کند
		// (منطق این تابع را برای خوانایی خلاصه می‌کنیم، چون قبلاً پیاده‌سازی شده)
		runConfiguration(cmd)
	},
}

// --- زیردستور 'configure path' ---
var configurePathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show the path to the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		showConfigPath()
	},
}

// --- زیردستور 'configure edit' ---
var configureEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open the configuration file in the default editor",
	Run: func(cmd *cobra.Command, args []string) {
		editConfigFile()
	},
}

// توابع کمکی که منطق اصلی را انجام می‌دهند
func runConfiguration(cmd *cobra.Command) {
	fmt.Println("Starting interactive configuration...")
	// ... (منطق کامل پرسیدن سوالات با survey و گرفتن مقادیر از فلگ‌ها که قبلاً نوشتیم) ...

	// گرفتن مقادیر از viper به عنوان پیش‌فرض
	survey.AskOne(&survey.Input{Message: "Enter MongoDB Remote Server URI:", Default: viper.GetString("mongodb.remote_uri")}, &mongoRemoteURI, survey.WithValidator(survey.Required))
	survey.AskOne(&survey.Input{Message: "Enter MongoDB Local Server URI:", Default: viper.GetString("mongodb.local_uri")}, &mongoLocalURI, survey.WithValidator(survey.Required))
	survey.AskOne(&survey.Input{Message: "Enter path to store backups:", Default: viper.GetString("paths.backup")}, &backupPath, survey.WithValidator(survey.Required))
	survey.AskOne(&survey.Input{Message: "Enter path to MongoDB Database Tools 'bin' directory:", Default: viper.GetString("paths.mongo_tools")}, &mongoToolsPath)

	viper.Set("mongodb.remote_uri", mongoRemoteURI)
	viper.Set("mongodb.local_uri", mongoLocalURI)
	viper.Set("paths.backup", backupPath)
	viper.Set("paths.mongo_tools", mongoToolsPath)

	saveConfiguration()
}

func showConfigPath() {
	setupViperConfigPaths()
	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		fmt.Printf("Configuration file is located at:\n%s\n", configFile)
	} else {
		defaultPath := getConfigFilePath()
		fmt.Printf("Configuration file not yet created. It will be created at:\n%s\n", defaultPath)
	}
}

func editConfigFile() {
	setupViperConfigPaths()
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		configFile = getConfigFilePath()
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			log.Printf("Configuration file not found. Creating a new one at: %s\n", configFile)
			saveConfiguration()
		}
	}

	fmt.Printf("Opening config file: %s\n", configFile)

	var editorCmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		editorCmd = exec.Command("cmd", "/C", "start", "", configFile)
	case "darwin":
		editorCmd = exec.Command("open", configFile)
	case "linux":
		editorCmd = exec.Command("xdg-open", configFile)
	default:
		fmt.Printf("Unsupported OS: %s. Please open the file manually at:\n%s\n", runtime.GOOS, configFile)
		return
	}

	err := editorCmd.Run()
	if err != nil {
		log.Printf("Failed to open editor: %v. Please open the file manually at:\n%s\n", err, configFile)
	}
}

func saveConfiguration() {
	configFile := getConfigFilePath()
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatalf("Error creating config directory: %v", err)
	}
	if err := viper.WriteConfigAs(configFile); err != nil {
		log.Fatalf("Error writing config file: %v", err)
	}
	viper.SetConfigFile(configFile)
	fmt.Printf("Configuration saved successfully to: %s\n", configFile)
}

func getConfigFilePath() string {
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		return filepath.Join(home, ".dataweaver-cli", "config.yaml")
	}
	currentDir, _ := os.Getwd()
	return filepath.Join(currentDir, "config.yaml")
}

func setupViperConfigPaths() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	home, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(filepath.Join(home, ".dataweaver-cli"))
	}
	viper.ReadInConfig()
}

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.AddCommand(configurePathCmd)
	configureCmd.AddCommand(configureEditCmd)

	// تعریف فلگ‌ها برای دستور configure
	configureCmd.Flags().StringVarP(&mongoRemoteURI, "mongo-remote-uri", "r", "", "MongoDB remote server URI")
	configureCmd.Flags().StringVarP(&mongoLocalURI, "mongo-local-uri", "l", "mongodb://localhost:27017", "MongoDB local server URI")
	configureCmd.Flags().StringVarP(&backupPath, "backup-path", "b", "./backups", "Path to store backups")
	configureCmd.Flags().StringVarP(&mongoToolsPath, "mongo-tools-path", "t", "", "Path to MongoDB Tools 'bin' directory")
}
