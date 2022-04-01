
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
