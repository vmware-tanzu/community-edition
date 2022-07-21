// React imports
import React, { ChangeEvent, MouseEvent, useContext, useEffect, useState } from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsSelect } from '@cds/react/select';
import { CdsIcon } from '@cds/react/icon';
import { ClarityIcons, refreshIcon, connectIcon, infoCircleIcon } from '@cds/core/icon';
import { CdsRadioGroup, CdsRadio } from '@cds/react/radio';
import { SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';

// App import
import './ManagementCredentials.scss';
import { AwsService } from '../../../../../swagger-api/services/AwsService';
import { AwsStore } from '../../../../../state-management/stores/Store.aws';
import { AWSAccountParams } from '../../../../../swagger-api/models/AWSAccountParams';
import { AWSKeyPair } from '../../../../../swagger-api/models/AWSKeyPair';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { managementCredentialFormSchema } from './management.credential.form.schema';
import ManagementCredentialProfile from './ManagementCredentialProfile';
import ManagementCredentialOneTime from './ManagementCredentialOneTime';
import { StepProps } from '../../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import ConnectionNotification, { CONNECTION_STATUS } from '../../../../../shared/components/ConnectionNotification/ConnectionNotification';
ClarityIcons.addIcons(refreshIcon, connectIcon, infoCircleIcon);

export interface FormInputs {
    PROFILE: string;
    REGION: string;
    SECRET_ACCESS_KEY: string;
    SESSION_TOKEN: string;
    ACCESS_KEY_ID: string;
    EC2_KEY_PAIR: string;
}

/* eslint-disable no-unused-vars */
enum CREDENTIAL_TYPE {
    PROFILE = 'PROFILE',
    ONE_TIME = 'ONE_TIME',
}

function ManagementCredentials(props: Partial<StepProps>) {
    const { handleValueChange, currentStep, goToStep, submitForm } = props;
    const { awsState } = useContext(AwsStore);
    const [connectionStatus, setConnectionStatus] = useState<CONNECTION_STATUS>(CONNECTION_STATUS.DISCONNECTED);
    const [message, setMessage] = useState('');

    const methods = useForm<FormInputs>({
        resolver: yupResolver(managementCredentialFormSchema),
    });
    const {
        register,
        handleSubmit,
        setValue,
        formState: { errors },
    } = methods;

    const [type, setType] = useState<CREDENTIAL_TYPE>(CREDENTIAL_TYPE.PROFILE);

    const [regions, setRegions] = useState<string[]>([]);
    const [keypairs, setKeyPairs] = useState<AWSKeyPair[]>([]);

    useEffect(() => {
        // fetch regions
        AwsService.getAwsRegions().then((data) => setRegions(data));
    }, []);

    useEffect(() => {
        if (connectionStatus === CONNECTION_STATUS.CONNECTED) {
            AwsService.getAwsKeyPairs().then((data) => {
                setKeyPairs(data);
            });
        }
    }, [connectionStatus]);

    const selectCredentialType = (event: ChangeEvent<HTMLSelectElement>) => {
        setConnectionStatus(CONNECTION_STATUS.DISCONNECTED);
        setType(CREDENTIAL_TYPE[event.target.value as CREDENTIAL_TYPE]);
    };

    const onSubmit: SubmitHandler<FormInputs> = (data) => {
        if (connectionStatus === CONNECTION_STATUS.CONNECTED && Object.keys(errors).length === 0) {
            if (goToStep && currentStep && submitForm) {
                goToStep(currentStep + 1);
                submitForm(currentStep);
            }
        }
    };

    const handleConnect = async () => {
        let params: AWSAccountParams = {};
        if (type === CREDENTIAL_TYPE.PROFILE) {
            params = {
                profileName: awsState[STORE_SECTION_FORM].PROFILE,
                region: awsState[STORE_SECTION_FORM].REGION,
            };
        } else {
            params = {
                accessKeyID: awsState[STORE_SECTION_FORM].ACCESS_KEY_ID,
                region: awsState[STORE_SECTION_FORM].REGION,
                secretAccessKey: awsState[STORE_SECTION_FORM].SECRET_ACCESS_KEY,
                sessionToken: awsState[STORE_SECTION_FORM].SESSION_TOKEN,
            };
        }
        try {
            setConnectionStatus(CONNECTION_STATUS.CONNECTING);
            setMessage('Connecting to AWS');
            await AwsService.setAwsEndpoint(params);
            setConnectionStatus(CONNECTION_STATUS.CONNECTED);
            setMessage('Connected to AWS');
        } catch (err: any) {
            setConnectionStatus(CONNECTION_STATUS.ERROR);
            setMessage(`Unable to connect to AWS: ${err.body.message}`);
        }
    };

    const handleSelectProfile = (profile: string) => {
        setConnectionStatus(CONNECTION_STATUS.DISCONNECTED);
        if (handleValueChange) {
            setTimeout(() => {
                handleValueChange(INPUT_CHANGE, 'PROFILE', profile, currentStep, errors);
            });
        }
    };

    const handleSelectRegion = (region: string) => {
        setConnectionStatus(CONNECTION_STATUS.DISCONNECTED);
        if (handleValueChange) {
            setTimeout(() => {
                handleValueChange(INPUT_CHANGE, 'REGION', region, currentStep, errors);
            });
        }
    };

    const handleSelectKeyPair = (event: ChangeEvent<HTMLSelectElement>) => {
        setValue('EC2_KEY_PAIR', event.target.value, { shouldValidate: true });
        if (handleValueChange) {
            setTimeout(() => {
                handleValueChange(INPUT_CHANGE, 'EC2_KEY_PAIR', event.target.value, currentStep, errors);
            });
        }
    };

    const handleInputChange = (field: string, value: string) => {
        setConnectionStatus(CONNECTION_STATUS.DISCONNECTED);
        if (handleValueChange) {
            setTimeout(() => {
                handleValueChange(INPUT_CHANGE, field, value, currentStep, errors);
            });
        }
    };

    const handleRefresh = (event: MouseEvent<HTMLAnchorElement>) => {
        event.preventDefault();
        if (connectionStatus === CONNECTION_STATUS.CONNECTED) {
            AwsService.getAwsKeyPairs().then((data) => {
                setKeyPairs(data);
            });
        }
    };

    return (
        <div className="wizard-content-container">
            <h2 cds-layout="m-t:lg">Amazon Web Services Credentials</h2>
            <CdsRadioGroup layout="vertical-inline" onChange={selectCredentialType}>
                <label cds-text="section medium" cds-layout="m-b:md">
                    Credential Type
                </label>
                <CdsRadio>
                    <label cds-layout="p-r:xxl">AWS credential profile</label>
                    <input type="radio" value={CREDENTIAL_TYPE.PROFILE} checked={type === CREDENTIAL_TYPE.PROFILE} readOnly />
                </CdsRadio>
                <CdsRadio>
                    <label>One-time credential</label>
                    <input type="radio" value={CREDENTIAL_TYPE.ONE_TIME} checked={type === CREDENTIAL_TYPE.ONE_TIME} readOnly />
                </CdsRadio>
            </CdsRadioGroup>
            {type === CREDENTIAL_TYPE.PROFILE && (
                <ManagementCredentialProfile
                    handleSelectProfile={handleSelectProfile}
                    handleSelectRegion={handleSelectRegion}
                    initialProfile={awsState[STORE_SECTION_FORM].PROFILE}
                    initialRegion={awsState[STORE_SECTION_FORM].REGION}
                    regions={regions}
                    methods={methods}
                />
            )}
            {type === CREDENTIAL_TYPE.ONE_TIME && (
                <ManagementCredentialOneTime
                    initialRegion={awsState[STORE_SECTION_FORM].REGION}
                    initialSecretAccessKey={awsState[STORE_SECTION_FORM].SECRET_ACCESS_KEY}
                    initialSessionToken={awsState[STORE_SECTION_FORM].SESSION_TOKEN}
                    initialAccessKeyId={awsState[STORE_SECTION_FORM].ACCESS_KEY_ID}
                    handleSelectRegion={handleSelectRegion}
                    handleInputChange={handleInputChange}
                    regions={regions}
                    methods={methods}
                />
            )}
            <CdsFormGroup layout="vertical-inline" control-width="shrink">
                <div cds-layout="p-t:lg" className="aws-button-container">
                    <CdsButton
                        onClick={handleConnect}
                        disabled={connectionStatus === CONNECTION_STATUS.CONNECTED || !awsState[STORE_SECTION_FORM].REGION}
                    >
                        <CdsIcon shape="connect" size="md"></CdsIcon>
                        CONNECT
                    </CdsButton>
                    <ConnectionNotification message={message} status={connectionStatus}></ConnectionNotification>
                </div>
                <div cds-layout="horizontal gap:lg align:vertical-center">
                    <CdsSelect layout="compact">
                        <label>
                            EC2 key pair <CdsIcon shape="info-circle" size="md"></CdsIcon>
                        </label>
                        <select
                            className="select-md-width"
                            {...register('EC2_KEY_PAIR')}
                            defaultValue={awsState[STORE_SECTION_FORM].EC2_KEY_PAIR}
                            onChange={handleSelectKeyPair}
                            data-testid="keypair-select"
                        >
                            <option></option>
                            {keypairs.map((keypair) => (
                                <option key={keypair.id} value={keypair.name}>
                                    {keypair.name}
                                </option>
                            ))}
                        </select>
                        {errors['EC2_KEY_PAIR'] && (
                            <CdsControlMessage status="error" className="error-height">
                                {errors['EC2_KEY_PAIR'].message}
                            </CdsControlMessage>
                        )}
                        <CdsControlMessage className="control-message-width">
                            Connect with your AWS profile to view available EC2 key pairs.
                        </CdsControlMessage>
                    </CdsSelect>
                    <a
                        href="/Users/miclettej/Dev/miclettej-community-edition/community-edition/cli/cmd/plugin/ui/web/tanzu-ui/public"
                        className="btn-refresh icon-blue"
                        onClick={handleRefresh}
                        cds-text="secondary"
                    >
                        <CdsIcon shape="refresh" size="sm"></CdsIcon>{' '}
                        <span cds-layout="m-t:sm" className="vertical-mid">
                            REFRESH
                        </span>
                    </a>
                </div>
                <CdsButton onClick={handleSubmit(onSubmit)}>NEXT</CdsButton>
            </CdsFormGroup>
        </div>
    );
}

export default ManagementCredentials;
