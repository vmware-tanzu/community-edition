// Library imports
import * as yup from 'yup';
import { IPFAMILIES, VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';
import { isValidFqdn, isValidIp4, isValidIp6 } from '../../../../shared/validations/Validation.service';
import { AnyObject } from 'yup/es/types';

export function createSchema(ipFamily: IPFAMILIES) {
    return yup
        .object({
            [VSPHERE_FIELDS.SERVERNAME]: yupServerTest(ipFamily).required('vSphere server name is required'),
            [VSPHERE_FIELDS.USERNAME]: yup.string().required('username is required'),
            [VSPHERE_FIELDS.PASSWORD]: yup.string().required('password is required'),
            [VSPHERE_FIELDS.DATACENTER]: yup.string().required('Please select a data center'),
        })
        .required();
}

function yupServerTest(ipFamily: IPFAMILIES): yup.StringSchema<string | undefined, AnyObject, string | undefined> {
    if (ipFamily === IPFAMILIES.IPv6) {
        return yup
            .string()
            .test('', 'To use an IP (v4), toggle the IP family', (value: string | undefined) => !value || !isValidIp4(value))
            .test('', 'Please enter a valid IP (v6) or FQDN', (value: string | undefined) => isValidFqdn(value) || isValidIp6(value));
    }
    return yup
        .string()
        .test('', 'To use an IP (v6), toggle the IP family', (value: string | undefined) => !value || !isValidIp6(value))
        .test('', 'Please enter a valid IP (v4) or FQDN', (value: string | undefined) => isValidFqdn(value) || isValidIp4(value));
}
