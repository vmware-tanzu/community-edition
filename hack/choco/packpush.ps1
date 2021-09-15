# This script is for LOCAL DEVELOPMENT ONLY.
# Use it after you've made changes to the powershell scripts to clean up artifacts and get ready for re-installing.
# It assumes you are using a local directory in ${HOME}\tcerepo as a --source for `choco install`.
# Invoke from hack/choco directory

# Remove the nupkg file made by `choco pack` in the working dir
Remove-Item *.nupkg -ErrorAction Continue

# Remove the nupkgs from our local repo
Remove-Item ${HOME}\tce-pkg\*.nupkg -ErrorAction Continue

# Delete any versions that Chocolatey may have cached for us.
Remove-Item ${env:LOCALAPPDATA}\Temp\chocolatey\tanzu-community-edition\ -Recurse -ErrorAction Continue

# Bundle the powershell scripts and nuspec into a nupkg file
choco pack
# Upload the nupkg file to the directory we wanna use.
choco push --source ${HOME}\tce-pkg

# Now you can do `choco install tanzu-community-edition --source ${HOME}\tce-pkg`