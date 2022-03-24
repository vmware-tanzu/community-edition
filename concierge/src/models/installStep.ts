export interface InstallStep<PARAM> {
    name: string,
    execute: (arg: PARAM) => boolean,
}
