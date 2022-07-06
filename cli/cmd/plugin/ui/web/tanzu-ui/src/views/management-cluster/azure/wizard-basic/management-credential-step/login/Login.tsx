// React imports
import React, { ChangeEvent, useContext } from 'react';

// Library imports
import { CdsInput } from '@cds/react/input';
import { CdsIcon } from '@cds/react/icon';
import { UseFormReturn } from 'react-hook-form';
import { CdsControlMessage } from '@cds/react/forms';
import { CdsSelect } from '@cds/react/select';

// App imports
import { FormField, FormInputs } from '../ManagementCredentials';
import { STORE_SECTION_FORM } from '../../../../../../state-management/reducers/Form.reducer';
import { AzureStore } from '../../../../../../state-management/stores/Azure.store';
import { AzureClouds } from '../../../../../../shared/constants/App.constants';
import { CdsButton } from '@cds/react/button';
import { AZURE_FIELDS } from '../../../AzureManagementCluster.constants';

interface Props {
    connected: boolean;
    handleInputChange: (field: FormField, value: string) => void;
    handleConnect: () => void;
    methods: UseFormReturn<FormInputs, any>;
}

export default function Login(props: Props) {
    const { azureState } = useContext(AzureStore);
    const {
        methods: {
            formState: { errors },
            register,
        },
        connected,
        handleInputChange,
        handleConnect,
    } = props;

    const dataEntered = () => {
        const fields = [
            AZURE_FIELDS.SUBSCRIPTION_ID,
            AZURE_FIELDS.TENANT_ID,
            AZURE_FIELDS.CLIENT_ID,
            AZURE_FIELDS.CLIENT_SECRET,
            AZURE_FIELDS.AZURE_ENVIRONMENT,
        ];
        for (let i = 0; i < fields.length; i++) {
            if (azureState[STORE_SECTION_FORM][fields[i]] === '') {
                return false;
            }
        }
        return true;
    };

    return (
        <>
            <div cds-layout="horizontal gap:lg">
                <div cds-layout="horizontal gap:lg">
                    <CdsInput>
                        <label>Tenant ID</label>
                        <input
                            {...register(AZURE_FIELDS.TENANT_ID)}
                            placeholder="Tenant ID"
                            type="text"
                            onChange={(e: ChangeEvent<HTMLInputElement>) => handleInputChange(AZURE_FIELDS.TENANT_ID, e.target.value)}
                            value={azureState[STORE_SECTION_FORM].TENANT_ID}
                            className="large-input"
                        ></input>
                        {errors[AZURE_FIELDS.TENANT_ID] && (
                            <CdsControlMessage status="error">{errors[AZURE_FIELDS.TENANT_ID]?.message}</CdsControlMessage>
                        )}
                    </CdsInput>
                    <CdsInput>
                        <label>Client ID</label>
                        <input
                            {...register(AZURE_FIELDS.CLIENT_ID)}
                            placeholder="Client ID"
                            type="text"
                            onChange={(e: ChangeEvent<HTMLInputElement>) => handleInputChange(AZURE_FIELDS.CLIENT_ID, e.target.value)}
                            value={azureState[STORE_SECTION_FORM].CLIENT_ID}
                            className="large-input"
                        ></input>
                        {errors[AZURE_FIELDS.CLIENT_ID] && (
                            <CdsControlMessage status="error">{errors[AZURE_FIELDS.CLIENT_ID]?.message}</CdsControlMessage>
                        )}
                    </CdsInput>
                </div>
            </div>
            <div cds-layout="horizontal gap:lg">
                <div cds-layout="horizontal gap:md">
                    <CdsInput>
                        <label>Client Secret</label>
                        <input
                            {...register(AZURE_FIELDS.CLIENT_SECRET)}
                            placeholder="Client Secret"
                            type="password"
                            onChange={(e: ChangeEvent<HTMLInputElement>) => handleInputChange(AZURE_FIELDS.CLIENT_SECRET, e.target.value)}
                            value={azureState[STORE_SECTION_FORM].CLIENT_SECRET}
                        ></input>
                        {errors[AZURE_FIELDS.CLIENT_SECRET] && (
                            <CdsControlMessage status="error">{errors[AZURE_FIELDS.CLIENT_SECRET]?.message}</CdsControlMessage>
                        )}
                    </CdsInput>
                    <CdsInput>
                        <label>Subscription ID</label>
                        <input
                            {...register(AZURE_FIELDS.SUBSCRIPTION_ID)}
                            placeholder="Subscription ID"
                            type="text"
                            onChange={(e: ChangeEvent<HTMLInputElement>) => handleInputChange(AZURE_FIELDS.SUBSCRIPTION_ID, e.target.value)}
                            value={azureState[STORE_SECTION_FORM].SUBSCRIPTION_ID}
                        ></input>
                        {errors[AZURE_FIELDS.SUBSCRIPTION_ID] && (
                            <CdsControlMessage status="error">{errors[AZURE_FIELDS.SUBSCRIPTION_ID]?.message}</CdsControlMessage>
                        )}
                    </CdsInput>
                    <CdsSelect layout="compact">
                        <label>Azure environment</label>
                        <select
                            className="select-sm-width"
                            {...register(AZURE_FIELDS.AZURE_ENVIRONMENT)}
                            onChange={(e: ChangeEvent<HTMLSelectElement>) =>
                                handleInputChange(AZURE_FIELDS.AZURE_ENVIRONMENT, e.target.value)
                            }
                            defaultValue={azureState[STORE_SECTION_FORM].AZURE_ENVIRONMENT}
                            data-testid="azure-environment-select"
                        >
                            {AzureClouds.map((azureCloud) => (
                                <option key={azureCloud.name} value={azureCloud.name}>
                                    {azureCloud.displayName}
                                </option>
                            ))}
                        </select>
                        {errors[AZURE_FIELDS.AZURE_ENVIRONMENT] && (
                            <CdsControlMessage status="error">{errors[AZURE_FIELDS.AZURE_ENVIRONMENT]?.message}</CdsControlMessage>
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
        </>
    );
}
