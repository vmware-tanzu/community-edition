import { ProgressMessenger } from './progressMessage';

// TODO: rename ExistingInstallation to Installation and use for both existing and new
export interface InstallationState {
    currentStep: string,                         // name of the current step (for logging and progress messages)
    currentStepIndex: number                     // index of the current step
    totalSteps: number                           // number of total steps
    stop?: boolean,                              // TRUE if installation should halt (error encountered)
    success?: boolean,                           // TRUE if installation completed successfully
    existingInstallation?: ExistingInstallation, // if there's an existing installation, this object describes it
    pathInstallFiles?: string,                   // path to the unzipped tanzu files we're installing
    pathTanzuInstallation?: string,              // path of the tanzu destination for installation, aka TANZU_BIN_PATH in install.sh
    pathTanzuConfigHome?: string,                // path to the tanzu configuration file, aka XDG_CONFIG_HOME in install.sh
    pathTanzuDataHome?: string,                  // path to the data dir where a tce subdir lives with uninstall script, aka XDG_DATA_HOME
    tarballInfo?: InstallationTarball,           // information about the tarball file associated with this installation
    version?: string,                            // version of the software we're installing
    edition?: string,                            // edition of the software we're installing
}

export interface InstallStep {
    name: string,
    execute: (state: InstallationState, progressMessenger: ProgressMessenger) => InstallationState,
}

export interface ExistingInstallation {
    edition: string,
    path: string,
    version: string,
}

export interface InstallationTarball {
    dir: string,
    file: string,
    fullPath: string,
}
