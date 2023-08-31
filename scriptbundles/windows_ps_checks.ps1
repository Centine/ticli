# Define a list to hold the results
$results = @()

# Check if script is running with elevated permissions
$IsElevated = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $IsElevated) {
    Write-Host "ElevatedPermissions;FAIL;Please run the script with elevated permissions." -ForegroundColor Red
    Exit
} else {
    $results += "ElevatedPermissions;SUCCESS;"
}


# Check 1: Execution policy is not 'Restricted'
$executionPolicy = Get-ExecutionPolicy
if ($executionPolicy -ne "Restricted") {
    $results += "ExecutionPolicy Check;SUCCESS;Policy is $executionPolicy"
} else {
    $results += "ExecutionPolicy Check;FAIL;Policy is $executionPolicy"
}

# Check 2: Hyper-V Management Tools are installed
if (Get-WindowsFeature -Name RSAT-Hyper-V-Tools) {
    $results += "HyperV Management Tools Check;SUCCESS;Installed"
} else {
    $results += "HyperV Management Tools Check;WARNING;Not Installed"
}

# Check 3: Chocolatey package manager is installed
if (Test-Path 'C:\ProgramData\chocolatey\choco.exe') {
    $results += "Chocolatey Check;SUCCESS;Installed"
} else {
    $results += "Chocolatey Check;WARNING;Not Installed"
}

# Check 4: Microsoft Windows Terminal is installed
$terminalCheck = Get-AppxPackage | Where-Object { $_.Name -like "Microsoft.WindowsTerminal" }
if ($terminalCheck) {
    $results += "Microsoft Windows Terminal Check;SUCCESS;Installed"
} else {
    $results += "Microsoft Windows Terminal Check;WARNING;Not Installed"
}

# Check 5: Docker Desktop is installed
if (Test-Path 'C:\Program Files\Docker\Docker\Docker Desktop.exe') {
    $results += "Docker Desktop Installation Check;SUCCESS;Installed"
} else {
    $results += "Docker Desktop Installation Check;FAIL;Not Installed"
}

# Check 6: Docker Desktop is running
try {
    $dockerService = Get-Service -Name com.docker.service
    if ($dockerService.Status -eq "Running") {
        $results += "Docker Desktop Running Check;SUCCESS;Running"
    } else {
        $results += "Docker Desktop Running Check;WARNING;Not Running"
    }
} catch {
    $results += "Docker Desktop Running Check;WARNING;Service not found"
}

# Check 7: WSL version 2 is installed
try {
    $wslVersion = wsl --list --verbose | Where-Object { $_ -match '^\s*\*' } | ForEach-Object { ($_ -split '\s+')[2] }
    if ($wslVersion -eq "2") {
        $results += "WSL Version Check;SUCCESS;WSL version is 2"
    } else {
        $results += "WSL Version Check;FAIL;WSL version is $wslVersion"
    }
} catch {
    $results += "WSL Version Check;FAIL;WSL not installed or encountered an error"
}

# Check 8: Ubuntu 20.04 is the default WSL distribution
try {
    $defaultDist = wsl --list --verbose | Where-Object { $_ -match '^\s*\*' } | ForEach-Object { ($_ -split '\s+')[1] }
    if ($defaultDist -like "Ubuntu-20.04") {
        $results += "Ubuntu 20.04 Default Check;SUCCESS;Ubuntu 20.04 is the default"
    } else {
        $results += "Ubuntu 20.04 Default Check;WARNING;Default is $defaultDist"
    }
} catch {
    $results += "Ubuntu 20.04 Default Check;WARNING;Could not determine default distribution"
}


# Export the results to a CSV
$results | Out-File -FilePath 'C:\path\to\your\output.csv' -Encoding utf8 -Force

