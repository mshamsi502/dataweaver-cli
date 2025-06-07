# DataWeaver CLI
![DataWeaver CLI Cover](./.guthub/assets/dataveaver-cli-cover.jpg)


A simple and powerful CLI tool for database operations.

- [![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org/dl/)
- [![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
- [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](#)

DataWeaver CLI is a self-contained, cross-platform command-line tool designed to simplify and automate common database tasks, starting with backup and restore operations for MongoDB. It features an easy-to-use interactive menu and a self-updating mechanism for its dependencies.

## ✨ Key Features

- **MongoDB Backup & Restore**: Easily create compressed backups of your remote MongoDB and restore them to a local instance.
- **Automated Tooling**: The `download-tools` command automatically fetches and installs the necessary MongoDB Database Tools for your platform (currently supports Windows).
- **Interactive Menu**: A user-friendly menu for guided operations, making it easy for anyone to use without memorizing commands.
- **Configuration-Driven**: Manages all settings like connection strings and paths in a simple `config.yaml` file.
- **Portable**: Designed to keep its dependencies in a local project folder, avoiding the need for system-wide installations.

## 📁 Project Structure
```
dataweaver-cli/
│
├── .github/
│   └── assets/              # For repository-specific assets like images.
│
├── .vscode/                 # VS Code editor-specific settings (usually gitignored).
│
├── backups/                 # Default directory for storing backup archives.
│   └── backup-*.gz
│
├── cmd/                     # Source code for all CLI commands.
│   ├── root.go              # Defines the root command and the interactive menu.
│   ├── backup.go            # Defines the parent 'backup' command.
│   ├── backup_mongo.go      # Defines the 'backup mongo' subcommand.
│   ├── configure.go         # Defines the 'configure' command and its subcommands.
│   ├── download-tools.go    # Defines the 'download-tools' command.
│   ├── restore.go           # Defines the parent 'restore' command.
│   └── restore_mongo.go     # Defines the 'restore mongo' subcommand.
│
├── downloads/               # Stores temporary downloaded files (e.g., .zip archives).
│   └── mongodb-database-tools-windows-x86_64-100.12.2.zip
│
├── internal/                # Private application packages (not for external use).
│   ├── config/
│   │   └── config.go
│   ├── downloader/
│   └── mongodb/
│
├── scripts/
│   └── install_tools.ps1    # PowerShell script for automated installation on Windows.
│
├── tools/                   # Location for extracted/installed tools for portable use.
│   └── mongodb-database-tools-windows-x86_64-100.12.2/
│       ├── bin/
│       │   ├── mongodump.exe
│       │   └── mongorestore.exe
│       │   └── (and other tools)...
│       ├── THIRD-PARTY-NOTICES
│       └── LICENSE.md
│
├── .gitignore               # Specifies intentionally untracked files to ignore.
├── dataweaver-cli.exe       # Compiled application binary.
├── go.mod                   # Go module definition file.
├── go.sum                   # Dependency checksums.
├── LICENSE                  # Your project's software license.
├── main.go                  # Main application entry point.
└── README.md                # The project's documentation file.
```

### Directory Explanations:
- ```.github/```: Contains GitHub-specific files, such as workflow definitions for CI/CD or issue templates. The ```assets``` subfolder is a good place to store images for the ```README```.
- ```cmd/```: Contains the main application code for all your CLI commands. Following the standard Go project layout, each command can have its own file for better organization.
- ```internal/```: This directory holds the core logic of your application that is not meant to be imported by other projects. This includes packages for configuration management, database-specific logic, etc.
- ```scripts/```: Includes helper scripts used by the project. In this case, it holds the PowerShell script responsible for automating the installation of MongoDB Tools on Windows.
- ```backups/```, ```downloads/```, ```tools/```: These directories are created by the application at runtime to store backups, downloaded files, and the extracted tools, respectively. They should typically be added to your ```.gitignore``` file.


## 🚀 Getting Started

### Prerequisites

- **Go**: Version 1.20 or higher is required to build from source.
- **PowerShell**: Required on Windows for the automated installer script.

### Installation

You can get DataWeaver CLI in one of two ways:

#### 1. From a Release (Recommended)

This is the easiest way to get started.

1.  Go to the [**Releases**](https://github.com/YOUR_USERNAME/dataweaver-cli/releases) page for this project.
2.  Download the appropriate binary for your operating system (e.g., `dataweaver-cli.exe` for Windows).
3.  For Windows, also download the `install_tools.ps1` script and place it in the same directory as the executable.
4.  Run the application from your terminal.

#### 2. From Source (For Developers)

If you have Go installed, you can build the CLI from source.

1.  Clone the repository:
    ```bash
    git clone [https://github.com/YOUR_USERNAME/dataweaver-cli.git](https://github.com/YOUR_USERNAME/dataweaver-cli.git)
    cd dataweaver-cli
    ```
2.  Fetch the dependencies:
    ```bash
    go mod tidy
    ```
3.  Run the application:
    ```bash
    go run main.go
    ```

## 📖 Usage

Using the tool is a simple, three-step process for the first run.

### Step 1: Download & Install Dependencies

The very first time you run the tool, you need to download the required MongoDB command-line tools. Our CLI automates this for you.

```bash
dataweaver-cli download-tools
```

This will download the tools, and on Windows, it will launch an installer script which will prompt for Administrator privileges to install the tools system-wide and find the installation path automatically.

### Step 2: Configure the CLI
Next, set up your database connection strings and paths. Run the configure command:

```bash
dataweaver-cli configure
```

This will launch an interactive wizard that asks for:

- Remote MongoDB URI: The connection string for the database you want to back up.
- Local MongoDB URI: The connection string for the database where you want to restore backups.
- Backup Path: A local directory where backup files will be stored (e.g., ```./backups```).
- Mongo Tools Path: This is set automatically by the ```download-tools``` command.

### Step 3: Use the Interactive Menu
Now you are all set! Just run the tool without any commands to open the main menu.

```bash
dataweaver-cli
```
You will see a list of options:

- Backup MongoDB
- Restore MongoDB
- Configure Settings
- Download/Setup Tools
- ...and more.

Simply select an option and follow the prompts.

## 🕹️ Command Reference
You can also use commands directly without the interactive menu.
```
dataweaver-cli
├── configure              # Configure application settings interactively or with flags.
│   ├── path               # Show the path to the active configuration file.
│   └── edit               # Open the configuration file in the default editor.
│
├── download-tools         # Download and set up required dependencies (e.g., MongoDB Tools).
│
├── backup
│   └── mongo              # Create a new backup of the remote MongoDB database.
│
└── restore
    └── mongo              # Restore a MongoDB database from an existing backup.
```

## ⚙️ Configuration
The CLI uses a ```config.yaml``` file to store settings. This file is typically located at:

- Windows: ```C:\Users\<YourUser>\.dataweaver-cli\config.yaml```
- Linux/macOS: ```~/.dataweaver-cli/config.yaml```

A typical configuration file looks like this:
```YAML
mongodb:
  local_uri: mongodb://localhost:27017
  remote_uri: mongodb://user:password@remote-host:27017/mydatabase
paths:
  backup: ./backups
  mongo_tools: C:\Program Files\MongoDB\Tools\100.12.1\bin
```
## 🛣️ Roadmap
This project is actively being developed. Future enhancements include:

- [ ] Support for other databases (e.g., PostgreSQL, QuestDB).
- [ ] Add support for Linux and macOS to the ```download-tools``` command.
- [ ] Add progress bars for long-running operations like downloads and backups.
- [ ] Add more backup management commands (e.g., list, clean old backups).

## 🤝 Contributing

Contributions are welcome! If you have a feature request, bug report, or want to contribute code, please feel free to open an issue or submit a pull request.

- Fork the repository.
- Create your feature branch (```git checkout -b feature/AmazingFeature```).
- Commit your changes (```git commit -m 'Add some AmazingFeature'```).
- Push to the branch (```git push origin feature/AmazingFeature```).
- Open a Pull Request.

## 📄 License
This project is licensed under the MIT License. See the LICENSE file for details.