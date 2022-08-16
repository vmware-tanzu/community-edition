// Library imports
import * as yup from 'yup';

// App imports
import { AWS_FIELDS } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';

export const managementCredentialFormSchema = yup
    .object({
        [AWS_FIELDS.PROFILE]: yup.string(),
        [AWS_FIELDS.REGION]: yup.string().required('Please select an AWS region'),
        [AWS_FIELDS.SECRET_ACCESS_KEY]: yup.string(),
        [AWS_FIELDS.SESSION_TOKEN]: yup.string(),
        [AWS_FIELDS.ACCESS_KEY_ID]: yup.string(),
        [AWS_FIELDS.EC2_KEY_PAIR]: yup.string().required('Please select an EC2 key pair'),
    })
    .required();
