export function first<ARRAY_TYPE>(source: ARRAY_TYPE[]): ARRAY_TYPE | undefined {
    return elementAt(source, 0);
}

export function last<ARRAY_TYPE>(source: ARRAY_TYPE[]): ARRAY_TYPE | undefined {
    return elementAt(source, source?.length - 1);
}

export function middle<ARRAY_TYPE>(source: ARRAY_TYPE[]): ARRAY_TYPE | undefined {
    // NOTE: for an even array, we've decided to take the FIRST of the two middle elements, so we use floor() instead of round()
    const index = Math.floor((source?.length - 1) / 2);
    return elementAt(source, index);
}

/**
 * Find the first matching element which is also in options.
 * @param options is a list of string.
 * @param validMatch contains all matching options which have 0 or 1 element in the options.
 * @returns string is the element in both options and validMatch.
 */
export function findFirstMatchingOption(options: string[], validMatch: string[]): string {
    const set = new Set(options);
    for (let i = 0; i < validMatch.length; i++) {
        if (set.has(validMatch[i])) {
            return validMatch[i];
        }
    }
    return '';
}

function elementAt<ARRAY_TYPE>(source: ARRAY_TYPE[], index: number): ARRAY_TYPE | undefined {
    return source && source.length > index && index >= 0 ? source[index] : undefined;
}
