// Library imports
import * as yup from 'yup';

export const managementCredentialFormSchema = yup
    .object({
        TENANT_ID: yup.string().required('Please enter a tenant ID'),
        CLIENT_ID: yup.string().required('Please enter a client ID'),
        CLIENT_SECRET: yup.string().required('Please enter a client secret'),
        SUBSCRIPTION_ID: yup.string().required('Please enter a subscription ID'),
        AZURE_ENVIRONMENT: yup.string().required('Please select an Azure environment'),
        REGION: yup.string().required('Please select a region'),
        SSH_PUBLIC_KEY: yup.string().required('Please enter an SSH public key'),
    })
    .required();
