// React imports
import React, { ChangeEvent, useContext, useEffect, useState } from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsTextarea } from '@cds/react/textarea';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';

// App import
import { AzureAccountParams, AzureLocation, AzureService, ApiError } from '../../../../../swagger-api';
import { AzureOrchestrator } from '../azure-orchestrator/AzureOrchestrator.service';
import { AzureStore } from '../../store/Azure.store.mc';
import { AZURE_FIELDS } from '../../azure-mc-basic/AzureManagementClusterBasic.constants';
import { CONNECTION_STATUS } from '../../../../../shared/components/ConnectionNotification/ConnectionNotification';
import { FormAction } from '../../../../../shared/types/types';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { managementCredentialFormSchema } from './management.credential.form.schema';
import ManagementCredentialsLogin from './ManagementCredentialsLogin';
import { RESOURCE } from '../../../../../state-management/actions/Resources.actions';
import SpinnerSelect from '../../../../../shared/components/Select/SpinnerSelect';
import { StepProps } from '../../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import UseUpdateTabStatus from '../../../../../shared/components/wizard/UseUpdateTabStatus.hooks';
import { clearPreviousResourceData } from '../../../default-orchestrator/DefaultOrchestrator';
export interface FormInputs {
    [AZURE_FIELDS.TENANT_ID]: string;
    [AZURE_FIELDS.CLIENT_ID]: string;
    [AZURE_FIELDS.CLIENT_SECRET]: string;
    [AZURE_FIELDS.SUBSCRIPTION_ID]: string;
    [AZURE_FIELDS.AZURE_ENVIRONMENT]: string;
    [AZURE_FIELDS.REGION]: string;
    [AZURE_FIELDS.SSH_PUBLIC_KEY]: string;
}
const placeholderText = `Begins with 'ssh-rsa', 'ecdsa-sha2-nistp256', 'ecdsa-sha2-
nistp384', 'ecdsa-sha2-nistp521', 'ssh-ed25519', 'sk-ecdsa-sha2-nistp256@openssh.com', or 'sk-ssh-ed25519@openssh.com'`;

function ManagementCredentials(props: Partial<StepProps>) {
    const { currentStep, goToStep, submitForm, updateTabStatus } = props;
    const { azureState, azureDispatch } = useContext(AzureStore);
    const methods = useForm<FormInputs>({
        resolver: yupResolver(managementCredentialFormSchema),
        mode: 'all',
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
    const [errorObject, setErrorObject] = useState<{ [fieldName: string]: ApiError }>({});

    // update tab status bar
    if (updateTabStatus) {
        UseUpdateTabStatus(errors, currentStep, updateTabStatus);
    }

    const resetField = (field: string) => {
        setValue(AZURE_FIELDS.REGION, '');
        azureDispatch({
            type: INPUT_CHANGE,
            field,
            payload: '',
        } as FormAction);
    };

    const handleInputChange = (field: string, value: string) => {
        if (field !== AZURE_FIELDS.SSH_PUBLIC_KEY && field !== AZURE_FIELDS.REGION) {
            setConnectionStatus(CONNECTION_STATUS.DISCONNECTED);
            if (azureState[STORE_SECTION_FORM].REGION) {
                resetField(AZURE_FIELDS.REGION);
            }
        }
        azureDispatch({
            type: INPUT_CHANGE,
            field,
            payload: value,
        } as FormAction);
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
            subscriptionId: azureState[STORE_SECTION_FORM][AZURE_FIELDS.SUBSCRIPTION_ID],
            tenantId: azureState[STORE_SECTION_FORM][AZURE_FIELDS.SUBSCRIPTION_ID],
            clientId: azureState[STORE_SECTION_FORM][AZURE_FIELDS.SUBSCRIPTION_ID],
            clientSecret: azureState[STORE_SECTION_FORM][AZURE_FIELDS.SUBSCRIPTION_ID],
            azureCloud: azureState[STORE_SECTION_FORM][AZURE_FIELDS.SUBSCRIPTION_ID],
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

    const isConnected = () => {
        return connectionStatus === CONNECTION_STATUS.CONNECTED;
    };

    const hasError = () => {
        return Object.keys(errorObject).length > 0;
    };

    function showErrorInfo() {
        if (isConnected() && hasError()) {
            return (
                <div>
                    <div className="error-text">Error Occurred</div>
                    <br />
                    {Object.keys(errorObject).map((errorField) => {
                        return (
                            <CdsControlMessage status="error" key={errorField}>
                                {errorObject[errorField].status}
                                &nbsp;
                                {errorObject[errorField].statusText}
                                &nbsp;
                                {errorObject[errorField].body.message}
                                <br />
                            </CdsControlMessage>
                        );
                    })}
                    <br />
                </div>
            );
        }
        return;
    }

    useEffect(() => {
        if (azureState[STORE_SECTION_FORM][AZURE_FIELDS.REGION]) {
            AzureOrchestrator.initOsImages({ azureState, azureDispatch, errorObject, setErrorObject });
        } else {
            clearPreviousResourceData(azureDispatch, RESOURCE.AZURE_ADD_RESOURCES, AZURE_FIELDS.OS_IMAGE);
        }
    }, [azureState[STORE_SECTION_FORM][AZURE_FIELDS.REGION]]);

    return (
        <div className="wizard-content-container azure-credential">
            <h2 cds-layout="m-t:md m-b:xl" cds-text="title">
                Microsoft Azure Credentials
            </h2>
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
            <FormProvider {...methods}>
                <CdsFormGroup layout="vertical-inline" control-width="shrink">
                    <ManagementCredentialsLogin
                        status={connectionStatus}
                        message={message}
                        handleConnect={handleConnect}
                        handleInputChange={handleInputChange}
                    />
                    <div cds-layout="horizontal gap:lg">
                        <SpinnerSelect
                            label="Region"
                            className="select-sm-width"
                            disabled={connectionStatus !== CONNECTION_STATUS.CONNECTED}
                            handleSelect={(e: ChangeEvent<HTMLSelectElement>) => handleInputChange(AZURE_FIELDS.REGION, e.target.value)}
                            name={AZURE_FIELDS.REGION}
                            isLoading={regionLoading}
                            register={register}
                            error={errors[AZURE_FIELDS.REGION]?.message}
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
                        <div cds-layout="col:6 align:horizontal-center">{showErrorInfo()}</div>
                    </div>
                    <CdsButton onClick={handleSubmit(onSubmit)}>NEXT</CdsButton>
                </CdsFormGroup>
            </FormProvider>
        </div>
    );
}

export default ManagementCredentials;
