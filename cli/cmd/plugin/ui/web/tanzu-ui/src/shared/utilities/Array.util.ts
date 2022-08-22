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

function elementAt<ARRAY_TYPE>(source: ARRAY_TYPE[], index: number): ARRAY_TYPE | undefined {
    return source && source.length > index && index >= 0 ? source[index] : undefined;
}
