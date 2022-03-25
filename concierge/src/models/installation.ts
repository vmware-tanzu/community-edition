import { ProgressMessenger } from './progressMessage';

export interface InstallationState {
    currentStep: string,                         // name of the current step (for logging and progress messages)
    currentStepIndex: number                     // index of the current step
    totalSteps: number                           // number of total steps
    stop?: boolean,                              // TRUE if installation should halt (error encountered)
    success?: boolean,                           // TRUE if installation completed successfully
    existingInstallation?: ExistingInstallation, // if there's an existing installation, this object describes it
    pathToInstallFiles?: string,                 // path to the unzipped tanzu files we're installing
    pathToTanzuInstallation?: string,            // path of the tanzu destination for installation
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
