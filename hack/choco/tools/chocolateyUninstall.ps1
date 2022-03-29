# Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# Clean out old runtime directories/files, but *keep configs*
$deleteFiles = (
    "${HOME}\.cache\tanzu\catalog.yaml",
    "${HOME}\.config\tanzu\config.yaml",
    "${HOME}\.config\tanzu-plugins",
    "${HOME}\.config\tanzu\tkg\bom",
    "${HOME}\.config\tanzu\tkg\providers",
    "${HOME}\.config\tanzu\tkg\.tanzu.lock",
    "${HOME}\.config\tanzu\tkg\compatibility\tkg-compatibility.yaml",
    "${HOME}\tce",
    "${env:LOCALAPPDATA}\tanzu-cli",
    "${env:ProgramFiles}\tanzu"
)

Write-Host "Removing catalog files" -ForegroundColor Green
foreach ($file in $deleteFiles) {
    # If the file doesn't exist, that's fine, just keep going.
    Write-Host " - Removing file $file" -ForegroundColor Cyan
    Remove-Item -Force -Recurse -Path $file -ErrorAction SilentlyContinue
}
