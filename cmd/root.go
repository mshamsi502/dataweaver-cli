package cmd

import (
	"fmt"
	"os"

	"github.com/mshamsi502/dataweaver-cli/internal/config" // مسیر پکیج کانفیگ

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dataweaver-cli",
	Short: "A CLI tool for managing database operations",
	Long:  `DataWeaver CLI is a powerful, self-contained tool to handle backup, restore, and other database operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		// به جای ارجاع به متغیر سراسری rootCmd، خود cmd را به تابع پاس می‌دهیم
		runInteractiveMenu(cmd)
	},
}

// تابع اصلی برای اجرای منوی تعاملی
func runInteractiveMenu(cmd *cobra.Command) {
	fmt.Println("Welcome to DataWeaver CLI! (Press Ctrl+C to exit at any time)")

	// بارگذاری اولیه کانفیگ
	if err := config.LoadConfig(); err != nil {
		fmt.Printf("Notice: Could not load config file. Use 'Configure' to create one. (%v)\n", err)
	}

	// **تغییر اصلی اینجاست:**
	// ما از متغیر 'cmd' که به تابع پاس داده شده استفاده می‌کنیم، نه از 'rootCmd' سراسری.
	backupMongo, _, _ := cmd.Find([]string{"backup", "mongo"})
	restoreMongo, _, _ := cmd.Find([]string{"restore", "mongo"})
	configure, _, _ := cmd.Find([]string{"configure"})
	downloadTools, _, _ := cmd.Find([]string{"download-tools"})
	configurePath, _, _ := cmd.Find([]string{"configure", "path"})
	configureEdit, _, _ := cmd.Find([]string{"configure", "edit"})

	for {
		options := []string{
			"Backup MongoDB",
			"Restore MongoDB",
			"---", // جداکننده برای زیبایی
			"Configure Settings (Interactive)",
			"Download/Setup Tools",
			"---", // جداکننده
			"Edit Config File",
			"Show Config Path",
			"Exit",
		}

		var selectedOption string
		prompt := &survey.Select{
			Message:  "Choose an option:",
			Options:  options,
			PageSize: 10,
		}
		err := survey.AskOne(prompt, &selectedOption, survey.WithStdio(os.Stdin, os.Stderr, os.Stdout))
		if err != nil {
			return // خروج در صورت فشردن Ctrl+C
		}

		fmt.Println() // یک خط خالی برای فاصله

		switch selectedOption {
		case "Backup MongoDB":
			fmt.Println("--- Running Backup ---")
			// اجرای دستور پیدا شده
			if backupMongo != nil {
				backupMongo.Run(backupMongo, []string{})
			}
			fmt.Println("--- Backup Finished ---\n")

		case "Restore MongoDB":
			fmt.Println("\n--- Running Restore ---")
			if restoreMongo != nil {
				restoreMongo.Run(restoreMongo, []string{})
			}
			fmt.Println("--- Restore Finished ---\n")

		case "Configure Settings (Interactive)":
			fmt.Println("\n--- Running Interactive Configuration ---")
			if configure != nil {
				// تابع Run دستور configure را برای حالت تعاملی فراخوانی می‌کنیم
				runConfiguration(configure)
			}
			fmt.Println("--- Configuration Finished ---\n")

		case "Edit Config File":
			fmt.Println("\n--- Opening Config File ---")
			if configureEdit != nil {
				// تابع Run دستور 'configure edit' را فراخوانی می‌کنیم
				configureEdit.Run(configureEdit, []string{})
			}
			fmt.Println("--- Action Finished ---\n")

		case "Show Config Path":
			fmt.Println("\n--- Config File Path ---")
			if configurePath != nil {
				configurePath.Run(configurePath, []string{})
			}
			fmt.Println("--- Done ---\n")

		case "Download/Setup Tools":
			fmt.Println("\n--- Running Download/Setup Tools ---")
			if downloadTools != nil {
				downloadTools.Run(downloadTools, []string{})
			}
			fmt.Println("--- Download/Setup Finished ---\n")

		case "Exit":
			fmt.Println("Exiting. Goodbye!")
			return

		case "---":
			// این گزینه کاری انجام نمی‌دهد و فقط جداکننده است
			continue
		}

		fmt.Println("\nPress Enter to return to the main menu...")
		fmt.Scanln()
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// این تابع خالی است. ثبت دستورات در فایل‌های خودشان انجام می‌شود.
}
