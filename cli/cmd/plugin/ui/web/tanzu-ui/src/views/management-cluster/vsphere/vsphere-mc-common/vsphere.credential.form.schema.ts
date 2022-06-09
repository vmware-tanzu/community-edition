// Library imports
import * as yup from 'yup';
import { VSPHERE_FIELDS } from './VsphereManagementClusterCommon.constants';
import { isValidFqdn, isValidIp } from '../../../../shared/validations/Validation.service';

export const vsphereCredentialFormSchema = yup
    .object({
        // TODO: make these fields required
        [VSPHERE_FIELDS.SERVERNAME]: yup
            .string()
            .test('', 'Please enter a valid ip or fqdn', (value) => isValidFqdn(value) || isValidIp(value))
            .required('vSphere server name is required'),
        [VSPHERE_FIELDS.USERNAME]: yup.string().required('username is required'),
        [VSPHERE_FIELDS.PASSWORD]: yup.string().required('password is required'),
    })
    .required();
