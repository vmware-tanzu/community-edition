// Library imports
import * as yup from 'yup';
import { IPFAMILIES, VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';
import { isValidFqdn, isValidIp4, isValidIp6 } from '../../../../shared/validations/Validation.service';

export function createSchema(ipFamily: IPFAMILIES) {
    const { prompt, validator } = serverValidation(ipFamily);
    return yup
        .object({
            [VSPHERE_FIELDS.SERVERNAME]: yup.string().test('', prompt, validator).required('vSphere server name is required'),
            [VSPHERE_FIELDS.USERNAME]: yup.string().required('username is required'),
            [VSPHERE_FIELDS.PASSWORD]: yup.string().required('password is required'),
            [VSPHERE_FIELDS.DATACENTER]: yup.string().required('Please select a data center'),
        })
        .required();
}

// given an IP family (v4 or v6) returns an object that holds the appropriate prompt and validation method
function serverValidation(ipFamily: IPFAMILIES): { prompt: string; validator: (value: string | undefined) => boolean } {
    if (ipFamily === IPFAMILIES.IPv6) {
        return {
            prompt: 'Please enter a valid IP (v6) or FQDN',
            validator: (value: string | undefined) => isValidFqdn(value) || isValidIp6(value),
        };
    }
    return {
        prompt: 'Please enter a valid IP or FQDN',
        validator: (value: string | undefined) => isValidFqdn(value) || isValidIp4(value),
    };
}
