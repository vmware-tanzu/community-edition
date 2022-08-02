// React imports
import React, { ChangeEvent, MouseEvent, useContext, useEffect, useState } from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import { ClarityIcons, refreshIcon, connectIcon, infoCircleIcon } from '@cds/core/icon';
import { CdsRadioGroup, CdsRadio } from '@cds/react/radio';
import { SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { CdsFormGroup } from '@cds/react/forms';

// App import
import './ManagementCredentials.scss';
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
import SpinnerSelect from '../../../../../shared/components/Select/SpinnerSelect';
import { AWS_ADD_RESOURCES } from '../../../../../state-management/actions/Resources.actions';
import { AwsResourceAction } from '../../../../../shared/types/types';
import { AwsService, AWSVirtualMachine } from '../../../../../swagger-api';

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
    const { awsState, awsDispatch } = useContext(AwsStore);
    const [connectionStatus, setConnectionStatus] = useState<CONNECTION_STATUS>(CONNECTION_STATUS.DISCONNECTED);
    const [message, setMessage] = useState('');
    const [keyPairLoading, setKeyPairLoading] = useState(false);

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

    const fetchKeyPairs = async () => {
        try {
            setKeyPairLoading(true);
            const keyPairs = await AwsService.getAwsKeyPairs();
            setKeyPairs(keyPairs);
        } catch (e: any) {
            console.log(`Unabled to get ec2 key pair: ${e}`);
        } finally {
            setKeyPairLoading(false);
        }
    };

    const [osImages, setOsImages] = useState<AWSVirtualMachine[]>([]);
    const [errorMessage, setErrorMessage] = useState<any>();

    useEffect(() => {
        if (connectionStatus === CONNECTION_STATUS.CONNECTED) {
            fetchKeyPairs();
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
                awsDispatch({
                    type: AWS_ADD_RESOURCES,
                    resourceName: 'osImages',
                    payload: osImages,
                } as AwsResourceAction);
                awsDispatch({
                    type: AWS_ADD_RESOURCES,
                    resourceName: 'errors',
                    payload: errorMessage,
                } as AwsResourceAction);
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
    const resetField = (field: string) => {
        if (handleValueChange && awsState[STORE_SECTION_FORM][field]) {
            handleValueChange(INPUT_CHANGE, field, '', currentStep, errors);
            setValue('EC2_KEY_PAIR', '');
        }
    };
    const handleSelectProfile = (profile: string) => {
        setConnectionStatus(CONNECTION_STATUS.DISCONNECTED);
        resetField('EC2_KEY_PAIR');
        if (handleValueChange) {
            setTimeout(() => {
                handleValueChange(INPUT_CHANGE, 'PROFILE', profile, currentStep, errors);
            });
        }
    };

    const handleSelectRegion = (region: string) => {
        setConnectionStatus(CONNECTION_STATUS.DISCONNECTED);
        resetField('EC2_KEY_PAIR');
        if (handleValueChange) {
            setTimeout(() => {
                handleValueChange(INPUT_CHANGE, 'REGION', region, currentStep, errors);
            });
        }
        retrieveOsImages(region);
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
        resetField('EC2_KEY_PAIR');
        if (handleValueChange) {
            setTimeout(() => {
                handleValueChange(INPUT_CHANGE, field, value, currentStep, errors);
            });
        }
    };

    const handleRefresh = async (event: MouseEvent<HTMLAnchorElement>) => {
        event.preventDefault();
        if (connectionStatus === CONNECTION_STATUS.CONNECTED) {
            fetchKeyPairs();
        }
    };

    function retrieveOsImages(region: string | undefined) {
        try {
            setOsImages([]);
            AwsService.getAwsosImages(region).then((data) => {
                setOsImages(data);
            });
        } catch (e) {
            setErrorMessage(e);
        }
    }

    return (
        <div className="wizard-content-container">
            <h2 cds-layout="m-t:md m-b:xl" cds-text="title">
                Amazon Web Services Credentials
            </h2>
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
                    <SpinnerSelect
                        className="select-md-width"
                        disabled={connectionStatus !== CONNECTION_STATUS.CONNECTED}
                        label="EC2 key pair"
                        handleSelect={handleSelectKeyPair}
                        name="EC2_KEY_PAIR"
                        controlMessage="EC2 key pairs will be retrieved when connected to AWS."
                        isLoading={keyPairLoading}
                        register={register}
                        error={errors['EC2_KEY_PAIR']?.message}
                    >
                        <option></option>
                        {keypairs.map((keypair) => (
                            <option key={keypair.id} value={keypair.name}>
                                {keypair.name}
                            </option>
                        ))}
                    </SpinnerSelect>
                    <a
                        href="/"
                        className={
                            connectionStatus === CONNECTION_STATUS.CONNECTED && !keyPairLoading
                                ? 'btn-refresh icon-blue'
                                : 'btn-refresh disabled'
                        }
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
