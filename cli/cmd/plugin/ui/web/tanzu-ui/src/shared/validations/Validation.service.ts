/**
 * @method isValidFqdn decide if arg is a valid FQDN
 * @return boolean
 */
export const isValidFqdn = (arg: string | undefined) => {
    if (!arg) {
        return false;
    }
    const regexPattern = /^[a-z0-9]+([-.][a-z0-9]+)*\.[a-z]{2,}$/i;
    return regexPattern.test(arg.trim());
};

/**
 * @method isValidIp decide if arg is a valid IP after trimming whitespaces
 * @return boolean
 */
export const isValidIp = (arg: string | undefined) => {
    if (!arg) {
        return false;
    }
    const regexPattern =
        // eslint-disable-next-line max-len
        /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
    return regexPattern.test(arg.trim());
};

export const isValidClusterName = (arg: string | undefined) => {
    if (!arg) {
        return false;
    }
    const regexPattern = /^([a-z-]*)$/;
    return regexPattern.test(arg);
};
