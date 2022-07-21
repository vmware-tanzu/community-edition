// React imports
import React, { ChangeEvent, useContext, useState } from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
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
import ManagementCredentialsLogin from './ManagementCredentialsLogin';
import './ManagementCredentials.scss';
import { AZURE_FIELDS } from '../../AzureManagementCluster.constants';
import { CONNECTION_STATUS } from '../../../../../shared/components/ConnectionNotification/ConnectionNotification';
import SpinnerSelect from '../../../../../shared/components/Select/SpinnerSelect';

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
    const [regions, setRegions] = useState<AzureLocation[]>([]);
    const [connectionStatus, setConnectionStatus] = useState<CONNECTION_STATUS>(CONNECTION_STATUS.DISCONNECTED);
    const [message, setMessage] = useState('');
    const [regionLoading, setRegionLoading] = useState(false);

    const resetField = (field: FormField) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, field, '', currentStep, errors);
            setValue(field, '');
        }
    };

    const handleInputChange = (field: FormField, value: string) => {
        if (field !== AZURE_FIELDS.SSH_PUBLIC_KEY && field !== AZURE_FIELDS.REGION) {
            setConnectionStatus(CONNECTION_STATUS.DISCONNECTED);
            if (azureState[STORE_SECTION_FORM].REGION) {
                resetField(AZURE_FIELDS.REGION);
            }
        }
        setValue(field, value, { shouldValidate: true });
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, field, value, currentStep, errors);
        }
    };

    const retrieveRegions = async () => {
        try {
            setRegionLoading(true);
            const azureRegions = await AzureService.getAzureRegions();
            setRegions(azureRegions);
        } catch (e) {
            console.log(console.log(`Error when calling get azure regions API: ${e}`));
        } finally {
            setRegionLoading(false);
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
            setConnectionStatus(CONNECTION_STATUS.CONNECTING);
            setMessage('Connecting to Azure');
            await AzureService.setAzureEndpoint(params);
            setConnectionStatus(CONNECTION_STATUS.CONNECTED);
            setMessage('Connected to Azure');
            retrieveRegions();
        } catch (err: any) {
            setConnectionStatus(CONNECTION_STATUS.ERROR);
            setMessage(`Unable to connect to Azure: ${err.body.message}`);
        }
    };
    const onSubmit: SubmitHandler<FormInputs> = () => {
        if (CONNECTION_STATUS.CONNECTED && Object.keys(errors).length === 0) {
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
                <a
                    href="/Users/miclettej/Dev/miclettej-community-edition/community-edition/cli/cmd/plugin/ui/web/tanzu-ui/public"
                    className="text-blue"
                >
                    creating Microsoft Azure credentials
                </a>
            </p>
            <CdsFormGroup layout="vertical-inline" control-width="shrink">
                <ManagementCredentialsLogin
                    status={connectionStatus}
                    message={message}
                    methods={methods}
                    handleConnect={handleConnect}
                    handleInputChange={handleInputChange}
                />
                <div cds-layout="horizontal gap:lg">
                    <SpinnerSelect
                        label="Region"
                        className="select-sm-width"
                        disabled={connectionStatus !== CONNECTION_STATUS.CONNECTED}
                        handleSelect={(e: ChangeEvent<HTMLSelectElement>) => handleInputChange(AZURE_FIELDS.REGION, e.target.value)}
                        name="REGION"
                        isLoading={regionLoading}
                        register={register}
                        error={errors['REGION']?.message}
                    >
                        <option></option>
                        {regions.map((region) => (
                            <option key={region.name} value={region.name}>
                                {region.displayName}
                            </option>
                        ))}
                    </SpinnerSelect>
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
