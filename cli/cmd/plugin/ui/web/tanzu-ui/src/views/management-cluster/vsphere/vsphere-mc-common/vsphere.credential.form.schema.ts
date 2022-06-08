// Library imports
import * as yup from 'yup';

export const vsphereCredentialFormSchema = yup
    .object({
        // TODO: make these fields required
        SERVER: yup.string(), //.required(),
        USERNAME: yup.string(), //.required(),
        PASSWORD: yup.string(), //.required(),
    })
    .required();
