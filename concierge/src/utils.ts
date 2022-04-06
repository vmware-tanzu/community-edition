const fs = require('fs')
const { shellPathSync } = require('shell-path');

export interface FunctionResult {
    error?: boolean,
    displayMessage?: string,
    techMessage?: string,
    data?: any
}

export function stringContains(src, target: string): boolean {
    return src && target && src.indexOf(target) >= 0
}

export function percentage(numerator: number, denominator: number): number {
    if (denominator === 0) {
        return 0
    }
    const rawPercentage = (100 * numerator) / denominator
    return Math.round(rawPercentage)
}

export function fixPath(): FunctionResult {
    if (isWin32()) {
        return {techMessage: 'Win32 detected; no need to adjust path'}
    }
    const newPath = shellPathSync()
    if (newPath) {
        const oldPath = process.env.PATH
        process.env.PATH = newPath
        const techMessage = `Set path from ${oldPath} to ${newPath}`
        return {techMessage}
    }
    const techMessage = 'ERROR: unable to get shell path, so unable to set env path on non-Win32 OS!'
    const displayMessage = 'Due to an internal problem setting the shell path, we may be unable to execute shell commands'
    return {error: true, techMessage, displayMessage }
}

export function isWin32(): boolean {
    return process.platform === 'win32'
}

export function pathExists(path): boolean {
    try {
        fs.accessSync(path)
        return true
    } catch (e) {
        console.log(`pathExists(${path}) encounters error: ${JSON.stringify(e)}`)
        return false
    }
}
