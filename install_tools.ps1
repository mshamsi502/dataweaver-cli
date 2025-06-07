# ======================================================================================
#  PowerShell Script to Install MongoDB Tools and Return Install Path (V5 - Final Logic)
# ======================================================================================

# 1. --- Admin Privileges Check ---
if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    try {
        Start-Process powershell -Verb RunAs -ArgumentList "& `"$($myinvocation.mycommand.definition)`""
    } catch {
        Write-Error "Failed to re-launch as Administrator."
    }
    exit
}

# --- Script Body (Running as Admin) ---
$Host.UI.RawUI.WindowTitle = "MongoDB Tools Installer - DataWeaver CLI"

# 2. --- Define Search Paths and MSI Info ---
$MsiFileName = "mongodb-database-tools-windows-x86_64-100.12.1.msi"
$MsiFilePath = Join-Path $PSScriptRoot "files\$MsiFileName"
$uninstallPaths = @(
    "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\*",
    "HKLM:\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall\*",
    "HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall\*"
)
# از یک الگوی جستجوی دقیق‌تر بر اساس خروجی شما استفاده می‌کنیم
$searchPattern = "MongoDB Tools*" 

# 3. --- Search, Install, and Find Path Logic ---
Write-Host "Searching for existing MongoDB Tools installation..." -ForegroundColor Cyan

$mongoToolsEntry = $null
# --- STEP 1: Perform the search ONCE before installation ---
foreach ($path in $uninstallPaths) {
    $foundEntry = Get-ItemProperty $path -ErrorAction SilentlyContinue | Where-Object { ($_.DisplayName -like $searchPattern) -and ($_.DisplayName) }
    if ($foundEntry) {
        # اولین مورد پیدا شده رو انتخاب می‌کنیم
        $mongoToolsEntry = $foundEntry | Select-Object -First 1 
        break
    }
}

# --- STEP 2: Decide whether to install based on search result ---
if ($mongoToolsEntry) {
    Write-Host "MongoDB Tools already installed. Skipping installation." -ForegroundColor Yellow
} else {
    # If not found, try to install it
    Write-Host "MongoDB Tools not found. Starting silent installation..." -ForegroundColor Cyan
    if (-NOT (Test-Path $MsiFilePath)) {
        Write-Error "MSI file not found at '$MsiFilePath'. Please run 'download-tools' first."
        exit 1
    }
    $msiArgs = "/i `"$MsiFilePath`" /qn"
    $installProcess = Start-Process msiexec.exe -ArgumentList $msiArgs -Wait -PassThru

    if ($installProcess.ExitCode -ne 0) {
        Write-Error "Installation failed with exit code: $($installProcess.ExitCode)"
        exit $installProcess.ExitCode
    }
    
    Write-Host "Installation completed. Re-searching for registry entry..." -ForegroundColor Green
    # --- STEP 3: Re-search for the entry AFTER installation ---
    foreach ($path in $uninstallPaths) {
        $foundEntry = Get-ItemProperty $path -ErrorAction SilentlyContinue | Where-Object { $_.DisplayName -like $searchPattern }
        if ($foundEntry) {
            $mongoToolsEntry = $foundEntry | Select-Object -First 1
            break
        }
    }
}

# --- STEP 4: Extract the path from the (now hopefully found) entry ---
if ($mongoToolsEntry) {
    $installLocation = $mongoToolsEntry.InstallLocation
    if ($installLocation -and (Test-Path $installLocation)) {
        $binPath = Join-Path $installLocation "bin"
        
        # On success, write ONLY the path to the standard output. The Go program will capture this.
        Write-Output $binPath
        exit 0
    }
}

# If we reach here, we failed to find the path.
Write-Error "Could not automatically determine the installation 'bin' path from the registry."
exit 1