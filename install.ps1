$ErrorActionPreference = "Stop"

$RELEASES = "https://github.com/Dadarzz2405/ai-cli/releases/latest/download"
$BINARY   = "ai-windows-amd64.exe"
$INSTALL  = "$env:USERPROFILE\AppData\Local\Microsoft\WindowsApps"

Write-Host "-> downloading $BINARY..."
Invoke-WebRequest -Uri "$RELEASES/$BINARY" -OutFile "$env:TEMP\ai.exe"

Write-Host "-> installing to $INSTALL..."
Move-Item -Force "$env:TEMP\ai.exe" "$INSTALL\ai.exe"

Write-Host ""
Write-Host "done! run: ai"
Write-Host "first launch will walk you through setup."