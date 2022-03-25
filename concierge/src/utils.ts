
export function stringContains(src, target: string): boolean {
    return src && target && src.indexOf(target) >= 0
}

export function stringContainsAny(src: string, targets: string[]): boolean {
    return targets && targets.length && (targets.find(target => stringContains(src, target)) != undefined)
}
