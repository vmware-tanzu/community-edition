// Library imports
import * as yup from 'yup';

export const managementCredentialFormSchema = yup
    .object({
        TENANT_ID: yup.string().required(),
        CLIENT_ID: yup.string().required(),
        CLIENT_SECRET: yup.string().required(),
        SUBSCRIPTION_ID: yup.string().required(),
        AZURE_ENVIRONMENT: yup.string().required(),
        REGION: yup.string().required(),
        SSH_PUBLIC_KEY: yup.string().required(),
    })
    .required();
