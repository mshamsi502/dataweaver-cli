package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config struct to hold all configuration parameters
type Config struct {
	MongoDB struct {
		RemoteURI string `mapstructure:"remote_uri"`
		LocalURI  string `mapstructure:"local_uri"`
	} `mapstructure:"mongodb"`
	Paths struct {
		Backup     string `mapstructure:"backup"`
		MongoTools string `mapstructure:"mongo_tools"`
	} `mapstructure:"paths"`
}

// Global variable to hold the loaded configuration
var AppConfig Config

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() error {
	viper.SetConfigName("config") // نام فایل کانفیگ (بدون پسوند)
	viper.SetConfigType("yaml")   // نوع فایل کانفیگ

	// مسیر اول: پوشه فعلی (برای زمانی که شاید فایل کانفیگ کنار اجرایی باشد)
	viper.AddConfigPath(".")

	// مسیر دوم: پوشه .dataweaver-cli در دایرکتوری خانگی کاربر
	home, err := os.UserHomeDir()
	if err == nil {
		configDirPath := filepath.Join(home, ".dataweaver-cli")
		viper.AddConfigPath(configDirPath)
	} else {
		fmt.Fprintln(os.Stderr, "Warning: Could not get user home directory. Will only check current directory for config.", err)
	}

	// خواندن فایل کانفیگ
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// فایل کانفیگ پیدا نشد. این می‌تواند خطا نباشد اگر تنظیمات پیش‌فرض کافی هستند
			// یا اگر برنامه طوری طراحی شده که بدون فایل کانفیگ هم کار کند.
			// برای ابزار ما، احتمالاً نیاز داریم که کاربر ابتدا configure را اجرا کند.
			return fmt.Errorf("config file not found. Please run 'dataweaver-cli configure' first: %w", err)
		} else {
			// خطای دیگری هنگام خواندن فایل کانفیگ رخ داده
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal کردن (تبدیل) تنظیمات خوانده شده به struct AppConfig
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("unable to decode config into struct: %w", err)
	}

	fmt.Println("Configuration loaded successfully from:", viper.ConfigFileUsed())
	return nil
}
