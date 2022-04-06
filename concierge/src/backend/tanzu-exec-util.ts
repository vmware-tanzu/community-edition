import { ProgressMessage } from '_models/progressMessage';
import { spawn, spawnSync } from 'child_process';

export interface ExecOptions {
    stderrImpliesError?: boolean
}

export function exec(options: ExecOptions = {}, ...args: string[]) : ProgressMessage {
    if (!args || args.length === 0) {
        return {message: 'no arguments passed to doExec()?!', warning: true}
    }

    console.log(`doExec(): ${args.join(' ')}`)
    const result = {message: '', details: '', error: false }
    try {
        // pull off first arg to be the command
        const [command, ...commandArgs] = [...args]
        const execResult = spawnSync(command, commandArgs, {stdio: 'pipe', encoding: 'utf8'})
        console.log(`doExec(${command} ${commandArgs.join(' ')}) yields: ${JSON.stringify(execResult)}`)
        result.message = execResult.stdout?.toString()
        result.details = execResult.stderr?.toString()
        // TODO: LEFT OFF: the following line does not do what we're hoping, which is to set the error only if there was
        // output to stderr
        result.error = (result.details) && options.stderrImpliesError
    } catch (e) {
        console.log(e)
        result.message = 'ERROR: ' + e.toString()
        result.error = true
    }
    console.log(`doExec() returning: ${JSON.stringify(result)}`)
    return result
}

export function execAsync(options: ExecOptions = {}, ...args: string[]) : ProgressMessage {
    if (!args || args.length === 0) {
        return {message: 'no arguments passed to doExec()?!', warning: true}
    }

    console.log(`execAsync(): ${args.join(' ')}`)
    const result = {message: '', details: '', error: false }
    try {
        // pull off first arg to be the command
        const [command, ...commandArgs] = [...args]
        const execResult = spawn(command, commandArgs, {stdio: 'pipe'})
        console.log(`execAsync(${command}) yields: ${JSON.stringify(execResult)}`)
        result.message = execResult.stdout?.toString()
        result.details = execResult.stderr?.toString()
        // TODO: LEFT OFF: the following line does not do what we're hoping, which is to set the error only if there was
        // output to stderr
        result.error = (result.details) && options.stderrImpliesError
    } catch (e) {
        console.log(e)
        result.message = 'ERROR: ' + e.toString()
        result.error = true
    }
    console.log(`doExec() returning: ${JSON.stringify(result)}`)
    return result
}

function doExec(fxn: any, options: ExecOptions, ...args: string[]) : ProgressMessage {
    if (!args || args.length === 0) {
        return {message: 'no arguments passed to doExec()?!', warning: true}
    }

    console.log(`doExec(): ${args.join(' ')}`)
    const result = {message: '', details: '', error: false }
    try {
        // pull off first arg to be the command
        const [command, ...commandArgs] = [...args]
        const execResult = fxn(command, commandArgs, {stdio: 'pipe', encoding: 'utf8'})
        console.log(`doExec(${command}) yields: ${JSON.stringify(execResult)}`)
        result.message = execResult.stdout?.toString()
        result.details = execResult.stderr?.toString()
        // TODO: LEFT OFF: the following line does not do what we're hoping, which is to set the error only if there was
        // output to stderr
        result.error = (result.details) && options.stderrImpliesError
    } catch (e) {
        console.log(e)
        result.message = 'ERROR: ' + e.toString()
        result.error = true
    }
    console.log(`doExec() returning: ${JSON.stringify(result)}`)
    return result
}
