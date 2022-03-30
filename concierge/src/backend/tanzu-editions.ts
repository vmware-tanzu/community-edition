'use strict'
const platform = require('./os-platform.ts')
const TceEdition = 'tce'
const TceVersion0110 = { major: 0, minor: 11, patch: 0, released: 'March 29, 2022'} as TanzuVersion
const TceVersion0100 = { major: 0, minor: 10, patch: 0, released: 'Feb 10, 2022'} as TanzuVersion
const TceVersion0091 =  { major: 0, minor: 9, patch: 1, released: 'Sept 29, 2021'} as TanzuVersion

const TceMacRelease0110 = { edition: TceEdition, os: platform.osMac, version: TceVersion0110,
    sha256: 'b2e986117b1e4c1a4e3f256a688893cbb85e453d108fe50f17a8e969f38672a2' }
const TceMacRelease0100 = { edition: TceEdition, os: platform.osMac, version: TceVersion0100,
    sha256: '2ce93bb148383a754d65150105d345ccbdcad60054d95446b1e1c1f7a5765a26' }
const TceMacRelease0091 = { edition: TceEdition, os: platform.osMac, version: TceVersion0091,
    sha256: 'c07f55c6c12e658a33cd7049569b166bb26c70d89a8d5acf49a62295e351bb61' }
const TceReleasesMac = [TceMacRelease0110, TceMacRelease0100, TceMacRelease0091] as TanzuRelease[]

const TceWinRelease0110 = { edition: TceEdition, os: platform.osWin, version: TceVersion0110,
    sha256: '09793b35f3850933d67b5e9103a191c93c1b7570ede9040065138ef78c3127fd' }
const TceWinRelease0100 = { edition: TceEdition, os: platform.osWin, version: TceVersion0100,
    sha256: '22e6046a0f8d74c42b47d6b699ae242474a50b1005d65105d0e8209cce4440db' }
const TceWinRelease0091 = { edition: TceEdition, os: platform.osWin, version: TceVersion0091,
    sha256: 'd301a6d392bab76f0aa31e54724a9603b52a7a0957831b1b5b9516018f8a5ccd' }
const TceReleasesWin = [TceWinRelease0110, TceWinRelease0100, TceWinRelease0091] as TanzuRelease[]

const TceLinuxRelease0110 = { edition: TceEdition, os: platform.osLinux, version: TceVersion0110,
    sha256: '746b70015bbb90701cb20a9a2b61212250c7b0c12fb184694adfbdf49b9a6747' }
const TceLinuxRelease0100 = { edition: TceEdition, os: platform.osLinux, version: TceVersion0100,
    sha256: '0ea2b3f1c78028e56fdae0d12d1a6328d83b1776e15bb4738bf44aa34718e8f3' }
const TceLinuxRelease0091 = { edition: TceEdition, os: platform.osLinux, version: TceVersion0091,
    sha256: 'ad7cd54a77def27708f64793c76b9115e32b7e3d5893599a3432e39da36b3fc2' }
const TceReleasesLinux = [TceLinuxRelease0110, TceLinuxRelease0100, TceLinuxRelease0091] as TanzuRelease[]

const ReleaseMap = new Map([
    [ platform.osMac, TceReleasesMac],
    [ platform.osWin, TceReleasesWin],
    [ platform.osLinux, TceReleasesLinux],
]) as Map<string, TanzuRelease[]>

export interface TanzuRelease {
    edition: string,
    os: string,
    sha256: string,
    version: TanzuVersion,
}

export interface TanzuVersion {
    major: number,
    minor: number,
    patch: number,
    released: string,
}

// returns an array of TanzuReleases on the current platform
export function knownReleases() : TanzuRelease[] {
    if (ReleaseMap[process.platform]) {
        return ReleaseMap[process.platform]
    }
    console.log('WARNING: Unable to find releases for a process.platform of: ' + process.platform)
    return []
}

// returns a TanzuRelease that matches the given sha, or undefined if none found
export function findRelease(sha: string): TanzuRelease {
    return knownReleases().find(release => release.sha256 === sha)
}
