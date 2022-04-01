// Library imports
import * as yup from 'yup';

// App imports
import { isValidFqdn, isValidIp } from '../shared/validations/Validation.service';

export const authFormSchema = yup.object({
    VCENTER_SERVER: yup.string().required().test('', 'It is an invalid ip or fqdn', value => isValidFqdn(value) || isValidIp(value)),
    VCENTER_USERNAME:yup.string().min(2).max(65).required(),
    VCENTER_PASSWORD:yup.string().min(5).required()
}).required();