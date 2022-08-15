// React imports
import React, { ChangeEvent, useContext } from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage } from '@cds/react/forms';
import { CdsIcon } from '@cds/react/icon';
import { CdsSelect } from '@cds/react/select';
import { useFormContext } from 'react-hook-form';

// App imports
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import { AzureStore } from '../../../../../state-management/stores/Azure.store';
import { AzureClouds } from '../../../../../shared/constants/App.constants';
import { AZURE_FIELDS } from '../../azure-mc-basic/AzureManagementClusterBasic.constants';
import ConnectionNotification, { CONNECTION_STATUS } from '../../../../../shared/components/ConnectionNotification/ConnectionNotification';
import './ManagementCredentialsLogin.scss';
import TextInputWithError from '../../../../../shared/components/Input/TextInputWithError';

interface Props {
    status: CONNECTION_STATUS;
    message: string;
    handleInputChange: (field: string, value: string) => void;
    handleConnect: () => void;
}

export default function ManagementCredentialsLogin(props: Props) {
    const { azureState } = useContext(AzureStore);
    const { status, message, handleInputChange, handleConnect } = props;
    const {
        register,
        formState: { errors },
    } = useFormContext();

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
                    <TextInputWithError
                        defaultValue={azureState[STORE_SECTION_FORM][AZURE_FIELDS.TENANT_ID]}
                        label="Tenant ID"
                        name={AZURE_FIELDS.TENANT_ID}
                        handleInputChange={handleInputChange}
                        className="large-input"
                    />
                    <TextInputWithError
                        defaultValue={azureState[STORE_SECTION_FORM][AZURE_FIELDS.CLIENT_ID]}
                        label="Client ID"
                        name={AZURE_FIELDS.CLIENT_ID}
                        handleInputChange={handleInputChange}
                        className="large-input"
                    />
                </div>
            </div>
            <div cds-layout="horizontal gap:lg">
                <div cds-layout="horizontal gap:md">
                    <TextInputWithError
                        label="Client Secret"
                        name={AZURE_FIELDS.CLIENT_SECRET}
                        handleInputChange={handleInputChange}
                        defaultValue={azureState[STORE_SECTION_FORM][AZURE_FIELDS.CLIENT_SECRET]}
                    />
                    <TextInputWithError
                        label="Subscription ID"
                        name={AZURE_FIELDS.SUBSCRIPTION_ID}
                        handleInputChange={handleInputChange}
                        defaultValue={azureState[STORE_SECTION_FORM][AZURE_FIELDS.SUBSCRIPTION_ID]}
                    />
                    <CdsSelect layout="compact">
                        <label>Azure environment</label>
                        <select
                            className="select-sm-width"
                            {...register(AZURE_FIELDS.AZURE_ENVIRONMENT, {
                                onChange: (e: ChangeEvent<HTMLSelectElement>) => {
                                    handleInputChange(AZURE_FIELDS.AZURE_ENVIRONMENT, e.target.value);
                                },
                            })}
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
            <div cds-layout="p-t:lg" className="azure-button-container">
                <CdsButton onClick={handleConnect} disabled={status === CONNECTION_STATUS.CONNECTED || !dataEntered()}>
                    <CdsIcon shape="connect" size="md"></CdsIcon>
                    CONNECT
                </CdsButton>
                <ConnectionNotification message={message} status={status}></ConnectionNotification>
            </div>
        </>
    );
}
