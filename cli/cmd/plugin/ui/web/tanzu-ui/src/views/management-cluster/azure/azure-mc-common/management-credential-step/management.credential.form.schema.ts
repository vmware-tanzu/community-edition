// Library imports
import * as yup from 'yup';

// App imports
import { AZURE_FIELDS } from '../../azure-mc-basic/AzureManagementClusterBasic.constants';

export const managementCredentialFormSchema = yup
    .object({
        [AZURE_FIELDS.TENANT_ID]: yup.string().required('Please enter a tenant ID'),
        [AZURE_FIELDS.CLIENT_ID]: yup.string().required('Please enter a client ID'),
        [AZURE_FIELDS.CLIENT_SECRET]: yup.string().required('Please enter a client secret'),
        [AZURE_FIELDS.SUBSCRIPTION_ID]: yup.string().required('Please enter a subscription ID'),
        [AZURE_FIELDS.AZURE_ENVIRONMENT]: yup.string().required('Please select an Azure environment'),
        [AZURE_FIELDS.REGION]: yup.string().required('Please select a region'),
        [AZURE_FIELDS.SSH_PUBLIC_KEY]: yup.string().required('Please enter an SSH public key'),
    })
    .required();
