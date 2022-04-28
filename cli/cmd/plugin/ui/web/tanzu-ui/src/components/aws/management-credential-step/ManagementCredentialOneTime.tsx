// React imports
import React, { ChangeEvent, useEffect } from 'react';

// Library imports
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';
import { CdsSelect } from '@cds/react/select';

// App imports
import { FormInputs } from './ManagementCredentials';
import { UseFormReturn } from 'react-hook-form';
import './ManagementCredentialOneTime.scss';

interface Props {
    initialRegion: string;
    initialSecretAccessKey: string;
    initialSessionToken: string;
    initialAccessKeyId: string;
    regions: string[];
    handleSelectRegion: (region: string) => void;
    handleInputChange: (field: string, value: string) => void;
    methods: UseFormReturn<FormInputs, any>;
}

function ManagementCredentialOneTime(props: Props) {
    const {
        methods: {
            formState: { errors },
            register,
            setValue,
        },
        regions,
        initialRegion,
        initialSecretAccessKey,
        initialSessionToken,
        initialAccessKeyId,
        handleInputChange,
        handleSelectRegion,
    } = props;

    useEffect(() => {
        setValue('REGION', initialRegion);
    }, []);

    const handleRegionChange = (event: ChangeEvent<HTMLSelectElement>) => {
        setValue('REGION', event.target.value, { shouldValidate: true });
        handleSelectRegion(event.target.value);
    };

    const handleSecretAccessKeyChange = (
        event: ChangeEvent<HTMLInputElement>
    ) => {
        handleInputChange('SECRET_ACCESS_KEY', event.target.value);
    };

    const handleSessionTokenChange = (event: ChangeEvent<HTMLInputElement>) => {
        handleInputChange('SESSION_TOKEN', event.target.value);
    };

    const handleAccessKeyIdChange = (event: ChangeEvent<HTMLInputElement>) => {
        handleInputChange('ACCESS_KEY_ID', event.target.value);
    };
    
    return (
        <div className="credential-one-time-container">
            <p cds-layout="m-y:lg">
                Enter AWS account credentials directly in the Access Key ID and
                Secret Access Key fields for your Amazon Web Services account.
                Optionally specify an AWS session token in Session Token if your
                AWS account is configured to require temporary credentials.
            </p>
            <CdsFormGroup layout="vertical-inline" control-width="shrink">
                <div cds-layout="horizontal gap:lg align:vertical-top">
                    <div cds-layout="vertical gap:lg align:vertical-center">
                        <CdsInput>
                            <label>Secret access key</label>
                            <input
                                {...register('SECRET_ACCESS_KEY')}
                                placeholder="Secret access key"
                                type="password"
                                onChange={handleSecretAccessKeyChange}
                                value={initialSecretAccessKey}
                            ></input>
                            {errors['SECRET_ACCESS_KEY'] && (
                                <CdsControlMessage status="error">
                                    {errors['SECRET_ACCESS_KEY'].message}
                                </CdsControlMessage>
                            )}
                        </CdsInput>
                        <CdsInput>
                            <label>Session token</label>
                            <input
                                {...register('SESSION_TOKEN')}
                                placeholder="Session token"
                                type="password"
                                onChange={handleSessionTokenChange}
                                value={initialSessionToken}
                            ></input>
                            {errors['SESSION_TOKEN'] && (
                                <CdsControlMessage status="error">
                                    {errors['SESSION_TOKEN'].message}
                                </CdsControlMessage>
                            )}
                        </CdsInput>
                        <CdsInput>
                            <label>Access key ID</label>
                            <input
                                {...register('ACCESS_KEY_ID')}
                                placeholder="Access key ID"
                                type="password"
                                onChange={handleAccessKeyIdChange}
                                value={initialAccessKeyId}
                            ></input>
                            {errors['ACCESS_KEY_ID'] && (
                                <CdsControlMessage status="error">
                                    {errors['ACCESS_KEY_ID'].message}
                                </CdsControlMessage>
                            )}
                        </CdsInput>
                    </div>
                    <div cds-layout="vertical gap:lg align:vertical-top">
                        <CdsSelect layout="compact">
                            <label>AWS Region </label>
                            <select
                                className="select-sm-width"
                                {...register('REGION')}
                                onChange={handleRegionChange}
                                defaultValue={initialRegion}
                            >
                                <option></option>
                                {regions.map((region) => (
                                    <option key={region}> {region} </option>
                                ))}
                            </select>
                            {errors['REGION'] && (
                                <CdsControlMessage status="error">
                                    {errors['REGION'].message}
                                </CdsControlMessage>
                            )}
                        </CdsSelect>
                    </div>
                </div>
            </CdsFormGroup>
        </div>
    );
}

export default ManagementCredentialOneTime;
