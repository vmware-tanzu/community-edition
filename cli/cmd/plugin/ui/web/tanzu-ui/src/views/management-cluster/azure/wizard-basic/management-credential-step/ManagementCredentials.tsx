// React imports
import React, { ChangeEvent, useContext, useState } from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsSelect } from '@cds/react/select';
import { SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsTextarea } from '@cds/react/textarea';

// App import
import { AzureAccountParams, AzureLocation, AzureService } from '../../../../../swagger-api';
import { AzureStore } from '../../../../../state-management/stores/Azure.store';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { StepProps } from '../../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import { managementCredentialFormSchema } from './management.credential.form.schema';
import Login from './login/Login';
import './ManagementCredentials.scss';
import { AZURE_FIELDS } from '../../AzureManagementCluster.constants';

export interface FormInputs {
    TENANT_ID: string;
    CLIENT_ID: string;
    CLIENT_SECRET: string;
    SUBSCRIPTION_ID: string;
    AZURE_ENVIRONMENT: string;
    REGION: string;
    SSH_PUBLIC_KEY: string;
}
export type FormField = 'TENANT_ID' | 'CLIENT_ID' | 'CLIENT_SECRET' | 'SUBSCRIPTION_ID' | 'AZURE_ENVIRONMENT' | 'REGION' | 'SSH_PUBLIC_KEY';
const placeholderText = `Begins with 'ssh-rsa', 'ecdsa-sha2-nistp256', 'ecdsa-sha2-
nistp384', 'ecdsa-sha2-nistp521', 'ssh-ed25519', 'sk-ecdsa-sha2-nistp256@openssh.com', or 'sk-ssh-ed25519@openssh.com'`;

function ManagementCredentials(props: Partial<StepProps>) {
    const { handleValueChange, currentStep, goToStep, submitForm } = props;
    const { azureState } = useContext(AzureStore);
    const methods = useForm<FormInputs>({
        resolver: yupResolver(managementCredentialFormSchema),
    });
    const {
        register,
        handleSubmit,
        setValue,
        formState: { errors },
    } = methods;
    const [connected, setConnection] = useState(false);
    const [regions, setRegions] = useState<AzureLocation[]>([]);

    const handleInputChange = (field: FormField, value: string) => {
        if (field !== AZURE_FIELDS.SSH_PUBLIC_KEY && field !== AZURE_FIELDS.REGION) {
            setConnection(false);
        }
        setValue(field, value, { shouldValidate: true });
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, field, value, currentStep, errors);
        }
    };

    const retrieveRegions = async () => {
        try {
            const azureRegions = await AzureService.getAzureRegions();
            setRegions(azureRegions);
            setConnection(true);
        } catch (e) {
            console.log(console.log(`Error when calling get azure regions API: ${e}`));
        }
    };
    const handleConnect = async () => {
        const params: AzureAccountParams = {
            subscriptionId: azureState[STORE_SECTION_FORM].SUBSCRIPTION_ID,
            tenantId: azureState[STORE_SECTION_FORM].TENANT_ID,
            clientId: azureState[STORE_SECTION_FORM].CLIENT_ID,
            clientSecret: azureState[STORE_SECTION_FORM].CLIENT_SECRET,
            azureCloud: azureState[STORE_SECTION_FORM].AZURE_ENVIRONMENT,
        };
        try {
            await AzureService.setAzureEndpoint(params);
            setConnection(true);
            retrieveRegions();
        } catch (e) {
            console.log('Error');
        }
    };
    const onSubmit: SubmitHandler<FormInputs> = () => {
        if (connected && Object.keys(errors).length === 0) {
            if (goToStep && currentStep && submitForm) {
                goToStep(currentStep + 1);
                submitForm(currentStep);
            }
        }
    };
    return (
        <div className="wizard-content-container azure-credential">
            <h2 cds-layout="m-t:lg">Microsoft Azure Credentials</h2>
            <p cds-layout="m-y:lg" className="description">
                Provide the Azure user credentials to create the Management Server on Azure. Don&apos;t have Azure credentials? View our
                guide on&nbsp;
                <a href="/" className="text-blue">
                    creating Microsoft Azure credentials
                </a>
            </p>
            <CdsFormGroup layout="vertical-inline" control-width="shrink">
                <Login connected={connected} methods={methods} handleConnect={handleConnect} handleInputChange={handleInputChange} />
                <div cds-layout="horizontal gap:lg">
                    <CdsSelect layout="compact">
                        <label>Region </label>
                        <select
                            className="select-sm-width"
                            {...register(AZURE_FIELDS.REGION)}
                            onChange={(e: ChangeEvent<HTMLSelectElement>) => handleInputChange(AZURE_FIELDS.REGION, e.target.value)}
                            defaultValue={azureState[STORE_SECTION_FORM].REGION}
                            data-testid="region-select"
                        >
                            <option></option>
                            {regions.map((region) => (
                                <option key={region.name} value={region.name}>
                                    {region.displayName}
                                </option>
                            ))}
                        </select>
                        {errors[AZURE_FIELDS.REGION] && (
                            <CdsControlMessage status="error">{errors[AZURE_FIELDS.REGION]?.message}</CdsControlMessage>
                        )}
                    </CdsSelect>
                    <CdsTextarea status={errors[AZURE_FIELDS.SSH_PUBLIC_KEY] ? 'error' : 'neutral'}>
                        <label>SSH public key</label>
                        <textarea
                            {...register(AZURE_FIELDS.SSH_PUBLIC_KEY)}
                            onChange={(e: ChangeEvent<HTMLTextAreaElement>) =>
                                handleInputChange(AZURE_FIELDS.SSH_PUBLIC_KEY, e.target.value)
                            }
                            defaultValue={azureState[STORE_SECTION_FORM].SSH_PUBLIC_KEY}
                            placeholder={placeholderText}
                        ></textarea>
                        {errors[AZURE_FIELDS.SSH_PUBLIC_KEY] && (
                            <CdsControlMessage status="error">{errors[AZURE_FIELDS.SSH_PUBLIC_KEY]?.message}</CdsControlMessage>
                        )}
                    </CdsTextarea>
                </div>
                <CdsButton onClick={handleSubmit(onSubmit)}>NEXT</CdsButton>
            </CdsFormGroup>
        </div>
    );
}

export default ManagementCredentials;
