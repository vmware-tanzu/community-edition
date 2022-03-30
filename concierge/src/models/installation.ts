import { ProgressMessenger } from './progressMessage';

// TODO: rename ExistingInstallation to Installation and use for both existing and new
export interface InstallationState {
    currentStep: string,                         // name of the current step (for logging and progress messages)
    currentStepIndex: number                     // index of the current step
    totalSteps: number                           // number of total steps
    stop?: boolean,                              // TRUE if installation should halt (error encountered)
    success?: boolean,                           // TRUE if installation completed successfully
    existingInstallation?: ExistingInstallation, // if there's an existing installation, this object describes it
    dirInstallFiles?: string,                    // path to the unzipped tanzu files we're installing
    dirTanzuInstallation?: string,               // path of the tanzu destination for installation, aka TANZU_BIN_PATH in install.sh
    dirTanzuConfig?: string,                     // path to config dir , aka XDG_CONFIG_HOME in install.sh
    pathTanzuConfig?: string,                    // path to the tanzu configuration file
    dirTanzuData?: string,                       // path to the data dir where a tce subdir lives with uninstall script, aka XDG_DATA_HOME
    version?: string,                            // version of the software we're installing
    edition?: string,                            // edition of the software we're installing
    plugins?: string[],                          // list of plugins discovered in the installation tarball
    chosenInstallation?: AvailableInstallation,  // the installation the user has chosen
}

export interface InstallStep {
    name: string,
    execute: (state: InstallationState, progressMessenger: ProgressMessenger) => InstallationState,
}

export interface ExistingInstallation {
    edition: string,
    editionVersion: string,
    path: string,
    tanzuBinaryVersion: string,
}

export interface InstallationTarball {
    dir: string,
    file: string,
    fullPath: string,
}

export interface AvailableInstallation {
    version: string,
    edition: string,
    tarball: InstallationTarball,
}

export interface PreInstallation {
    existingInstallation: ExistingInstallation,
    availableInstallations: AvailableInstallation[],
    chosenInstallation?: AvailableInstallation,
}
