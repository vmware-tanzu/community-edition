export function removeErrorInfo(errorObject: { [key: string]: any }, field: string) {
    const copy = { ...errorObject };
    delete copy[field];
    return copy;
}

export function addErrorInfo(errorObject: { [key: string]: any }, error: any, field: string) {
    return {
        ...errorObject,
        [field]: error,
    };
}
