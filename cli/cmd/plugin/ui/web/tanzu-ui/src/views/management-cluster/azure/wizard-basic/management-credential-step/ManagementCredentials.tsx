// React imports
import React, { ChangeEvent, useContext, useState } from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsSelect } from '@cds/react/select';
import { CdsIcon } from '@cds/react/icon';
import { SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';
import { CdsTextarea } from '@cds/react/textarea';

// App import
import { AzureAccountParams, AzureLocation, AzureService } from '../../../../../swagger-api';
import { AzureClouds } from '../../../../../shared/constants/App.constants';
import { AzureStore } from '../../../../../state-management/stores/Azure.store';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { StepProps } from '../../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import { managementCredentialFormSchema } from './management.credential.form.schema';
import './ManagementCredentials.scss';

export interface FormInputs {
    TENANT_ID: string;
    CLIENT_ID: string;
    CLIENT_SECRET: string;
    SUBSCRIPTION_ID: string;
    AZURE_ENVIRONMENT: string;
    REGION: string;
    SSH_PUBLIC_KEY: string;
}
type FormField = 'TENANT_ID' | 'CLIENT_ID' | 'CLIENT_SECRET' | 'SUBSCRIPTION_ID' | 'AZURE_ENVIRONMENT' | 'REGION' | 'SSH_PUBLIC_KEY';
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
        if (field !== 'SSH_PUBLIC_KEY') {
            setConnection(false);
        }
        setValue(field, value, { shouldValidate: true });
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, field, value, currentStep, errors);
        }
    };

    const handleSelectChange = (field: string, value: string) => {
        if (field === 'AZURE_ENVIRONMENT') {
            setConnection(false);
        }
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

    const dataEntered = () => {
        const fields = ['SUBSCRIPTION_ID', 'TENANT_ID', 'CLIENT_ID', 'CLIENT_SECRET', 'AZURE_ENVIRONMENT'];
        for (let i = 0; i < fields.length; i++) {
            if (azureState[STORE_SECTION_FORM][fields[i]] === '') {
                return false;
            }
        }
        return true;
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
                <div cds-layout="horizontal gap:lg">
                    <div cds-layout="horizontal gap:lg">
                        <CdsInput>
                            <label>Tenant ID</label>
                            <input
                                {...register('TENANT_ID')}
                                placeholder="Tenant ID"
                                type="text"
                                onChange={(e: ChangeEvent<HTMLInputElement>) => handleInputChange('TENANT_ID', e.target.value)}
                                value={azureState[STORE_SECTION_FORM].TENANT_ID}
                                className="large-input"
                            ></input>
                            {errors['TENANT_ID'] && <CdsControlMessage status="error">{errors['TENANT_ID'].message}</CdsControlMessage>}
                        </CdsInput>
                        <CdsInput>
                            <label>Client ID</label>
                            <input
                                {...register('CLIENT_ID')}
                                placeholder="Client ID"
                                type="text"
                                onChange={(e: ChangeEvent<HTMLInputElement>) => handleInputChange('CLIENT_ID', e.target.value)}
                                value={azureState[STORE_SECTION_FORM].CLIENT_ID}
                                className="large-input"
                            ></input>
                            {errors['CLIENT_ID'] && <CdsControlMessage status="error">{errors['CLIENT_ID'].message}</CdsControlMessage>}
                        </CdsInput>
                    </div>
                </div>
                <div cds-layout="horizontal gap:lg">
                    <div cds-layout="horizontal gap:md">
                        <CdsInput>
                            <label>Client Secret</label>
                            <input
                                {...register('CLIENT_SECRET')}
                                placeholder="Client Secret"
                                type="password"
                                onChange={(e: ChangeEvent<HTMLInputElement>) => handleInputChange('CLIENT_SECRET', e.target.value)}
                                value={azureState[STORE_SECTION_FORM].CLIENT_SECRET}
                            ></input>
                            {errors['CLIENT_SECRET'] && (
                                <CdsControlMessage status="error">{errors['CLIENT_SECRET'].message}</CdsControlMessage>
                            )}
                        </CdsInput>
                        <CdsInput>
                            <label>Subscription ID</label>
                            <input
                                {...register('SUBSCRIPTION_ID')}
                                placeholder="Subscription ID"
                                type="text"
                                onChange={(e: ChangeEvent<HTMLInputElement>) => handleInputChange('SUBSCRIPTION_ID', e.target.value)}
                                value={azureState[STORE_SECTION_FORM].SUBSCRIPTION_ID}
                            ></input>
                            {errors['SUBSCRIPTION_ID'] && (
                                <CdsControlMessage status="error">{errors['SUBSCRIPTION_ID'].message}</CdsControlMessage>
                            )}
                        </CdsInput>
                        <CdsSelect layout="compact">
                            <label>Azure environment</label>
                            <select
                                className="select-sm-width"
                                {...register('AZURE_ENVIRONMENT')}
                                onChange={(e: ChangeEvent<HTMLSelectElement>) => handleSelectChange('AZURE_ENVIRONMENT', e.target.value)}
                                defaultValue={azureState[STORE_SECTION_FORM].AZURE_ENVIRONMENT}
                                data-testid="azure-environment-select"
                            >
                                {AzureClouds.map((azureCloud) => (
                                    <option key={azureCloud.name} value={azureCloud.name}>
                                        {azureCloud.displayName}
                                    </option>
                                ))}
                            </select>
                            {errors['AZURE_ENVIRONMENT'] && (
                                <CdsControlMessage status="error">{errors['AZURE_ENVIRONMENT'].message}</CdsControlMessage>
                            )}
                        </CdsSelect>
                    </div>
                </div>
                <div cds-layout="p-t:lg">
                    <CdsButton onClick={handleConnect} disabled={connected || !dataEntered()}>
                        <CdsIcon shape="connect" size="md"></CdsIcon>
                        {connected ? 'CONNECTED' : 'CONNECT'}
                    </CdsButton>
                </div>
                <div cds-layout="horizontal gap:lg">
                    <CdsSelect layout="compact">
                        <label>Region </label>
                        <select
                            className="select-sm-width"
                            {...register('REGION')}
                            onChange={(e: ChangeEvent<HTMLSelectElement>) => handleSelectChange('REGION', e.target.value)}
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
                        {errors['REGION'] && <CdsControlMessage status="error">{errors['REGION'].message}</CdsControlMessage>}
                    </CdsSelect>
                    <CdsTextarea status={errors['SSH_PUBLIC_KEY'] ? 'error' : 'neutral'}>
                        <label>SSH public key</label>
                        <textarea
                            {...register('SSH_PUBLIC_KEY')}
                            onChange={(e: ChangeEvent<HTMLTextAreaElement>) => handleInputChange('SSH_PUBLIC_KEY', e.target.value)}
                            defaultValue={azureState[STORE_SECTION_FORM].SSH_PUBLIC_KEY}
                            placeholder={placeholderText}
                        ></textarea>
                        {errors['SSH_PUBLIC_KEY'] && (
                            <CdsControlMessage status="error">{errors['SSH_PUBLIC_KEY'].message}</CdsControlMessage>
                        )}
                    </CdsTextarea>
                </div>
                <CdsButton onClick={handleSubmit(onSubmit)}>NEXT</CdsButton>
            </CdsFormGroup>
        </div>
    );
}

export default ManagementCredentials;
