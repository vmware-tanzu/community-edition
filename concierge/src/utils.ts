
export function stringContains(src, target: string): boolean {
    return src && target && src.indexOf(target) >= 0
}

export function stringContainsAny(src: string, targets: string[]): boolean {
    return targets && targets.length && (targets.find(target => stringContains(src, target)) != undefined)
}

export function percentage(numerator, denominator: number): number {
    if (denominator === 0) {
        return 0
    }
    const rawPercentage = (100 * numerator) / denominator
    return Math.round(rawPercentage)
}
