import { isIPv4, isIPv6 } from 'is-ip';

/**
 * @method isValidCidr decide if arg is a valid cidr
 * @return boolean
 */
export const isValidCidr = (target: string | undefined): boolean => {
    if (!target) {
        return false;
    }

    const argArr = target.split('/');
    if (argArr.length !== 2) {
        return false;
    }

    const ip = argArr[0];
    if (!isValidIp(ip)) {
        return false;
    }

    const rangeAsString = argArr[1];
    const ipRange = rangeAsString?.length > 0 ? +rangeAsString : -1;
    return ipRange >= 0 && ipRange < 32;
};

// Returns true if argument is a comma-separated list of valid IP or FQDN values
export const isValidCommaSeparatedIpOrFqdn = (arg: string | undefined): boolean => {
    if (!arg) {
        return false;
    }
    const ips = arg.split(',');
    return ips.map((ip) => isValidIp(ip) || isValidFqdn(ip)).reduce((a, b) => a && b, true);
};

export const isValidProxyServer = (arg: string | undefined): boolean => {
    if (!arg) {
        return false;
    }
    const stringParts = arg.split('://');
    const rightNumberOfStringParts = stringParts.length === 2;
    const protocol = stringParts[0];
    const validProtocol = protocol === 'http' || protocol === 'https';
    const server = rightNumberOfStringParts ? stringParts[1] : undefined;
    const validServer = server ? isValidIp(server) || isValidFqdn(server) : false;
    return validProtocol && validServer;
};

/**
 * @method isValidFqdn decide if arg is a valid FQDN
 * @return boolean
 */
export const isValidFqdn = (arg: string | undefined): boolean => {
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
export const isValidIp = (arg: string | undefined): boolean => {
    if (!arg) {
        return false;
    }
    const regexPattern =
        // eslint-disable-next-line max-len
        /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
    return regexPattern.test(arg.trim());
};

export const isValidIp4 = (ip: string | undefined): boolean => {
    if (!ip) {
        return false;
    }
    return isIPv4(ip);
};

export const isValidIp6 = (ip: string | undefined): boolean => {
    if (!ip) {
        return false;
    }
    return isIPv6(ip);
};

/**
 * @method isK8sCompliantString decide if arg is a valid k8s string, eg for cluster name
 * @return boolean
 */
export const isK8sCompliantString = (arg: string | undefined) => {
    if (!arg) {
        return false;
    }
    const regexPattern = /^[a-z0-9][a-z0-9-]*[a-z0-9]$/;
    return regexPattern.test(arg);
};
