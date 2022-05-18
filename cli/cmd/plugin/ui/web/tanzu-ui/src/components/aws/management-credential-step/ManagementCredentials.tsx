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
import { AwsService } from '../../../swagger-api/services/AwsService';
import { AwsStore } from '../../../state-management/stores/Store.aws';
import { AWSAccountParams } from '../../../swagger-api/models/AWSAccountParams';
import { AWSKeyPair } from '../../../swagger-api/models/AWSKeyPair';
import { INPUT_CHANGE } from '../../../state-management/actions/Form.actions';
import { managementCredentialFormSchema } from './management.credential.form.schema';
import ManagementCredentialProfile from './ManagementCredentialProfile';
import ManagementCredentialOneTime from './ManagementCredentialOneTime';
import { StepProps } from '../../../shared/components/wizard/Wizard';

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
    const [connected, setConnection] = useState(false);

    const [regions, setRegions] = useState<string[]>([]);
    const [keypairs, setKeyPairs] = useState<AWSKeyPair[]>([]);
    // const [amiObjects, setAmiObjects] = useState<AWSKeyPair[]>([]);

    useEffect(() => {
        // fetch regions
        AwsService.getAwsRegions().then((data) => setRegions(data));
    }, []);

    useEffect(() => {
        if (connected) {
            AwsService.getAwsKeyPairs().then((data) => {
                setKeyPairs(data);
            });

            AwsService.getAwsosImages(awsState.data.REGION).then((data) => {
                console.log(data);
                // setAmiObjects(data);
            });
        }
    }, [connected]); // eslint-disable-line react-hooks/exhaustive-deps

    const selectCredentialType = (event: ChangeEvent<HTMLSelectElement>) => {
        setConnection(false);
        setType(CREDENTIAL_TYPE[event.target.value as CREDENTIAL_TYPE]);
    };

    const onSubmit: SubmitHandler<FormInputs> = (data) => {
        if (connected && Object.keys(errors).length === 0) {
            if (goToStep && currentStep && submitForm) {
                goToStep(currentStep + 1);
                submitForm(currentStep);
            }
        }
    };

    const handleConnect = () => {
        let params: AWSAccountParams = {};
        if (type === CREDENTIAL_TYPE.PROFILE) {
            params = {
                profileName: awsState.data.PROFILE,
                region: awsState.data.REGION,
            };
        } else {
            params = {
                accessKeyID: awsState.data.ACCESS_KEY_ID,
                region: awsState.data.REGION,
                secretAccessKey: awsState.data.SECRET_ACCESS_KEY,
                sessionToken: awsState.data.SESSION_TOKEN,
            };
        }
        AwsService.setAwsEndpoint(params).then(() => {
            setConnection(true);
        });
    };

    const handleSelectProfile = (profile: string) => {
        setConnection(false);
        if (handleValueChange) {
            setTimeout(() => {
                handleValueChange(INPUT_CHANGE, 'PROFILE', profile, currentStep, errors);
            });
        }
    };

    const handleSelectRegion = (region: string) => {
        setConnection(false);
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
        setConnection(false);
        if (handleValueChange) {
            setTimeout(() => {
                handleValueChange(INPUT_CHANGE, field, value, currentStep, errors);
            });
        }
    };

    const handleRefresh = (event: MouseEvent<HTMLAnchorElement>) => {
        event.preventDefault();
        if (connected) {
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
                    initialProfile={awsState.data.PROFILE}
                    initialRegion={awsState.data.REGION}
                    regions={regions}
                    methods={methods}
                />
            )}
            {type === CREDENTIAL_TYPE.ONE_TIME && (
                <ManagementCredentialOneTime
                    initialRegion={awsState.data.REGION}
                    initialSecretAccessKey={awsState.data.SECRET_ACCESS_KEY}
                    initialSessionToken={awsState.data.SESSION_TOKEN}
                    initialAccessKeyId={awsState.data.ACCESS_KEY_ID}
                    handleSelectRegion={handleSelectRegion}
                    handleInputChange={handleInputChange}
                    regions={regions}
                    methods={methods}
                />
            )}
            <CdsFormGroup layout="vertical-inline" control-width="shrink">
                <div cds-layout="p-t:lg">
                    <CdsButton onClick={handleConnect} disabled={connected || !awsState.data.REGION}>
                        <CdsIcon shape="connect" size="md"></CdsIcon>
                        {connected ? 'CONNECTED' : 'CONNECT'}
                    </CdsButton>
                </div>
                <div cds-layout="horizontal gap:lg align:vertical-center">
                    <CdsSelect layout="compact">
                        <label>
                            EC2 key pair <CdsIcon shape="info-circle" size="md"></CdsIcon>
                        </label>
                        <select
                            className="select-md-width"
                            {...register('EC2_KEY_PAIR')}
                            defaultValue={awsState.data.EC2_KEY_PAIR}
                            onChange={handleSelectKeyPair}
                        >
                            <option></option>
                            {keypairs.map((keypair) => (
                                <option key={keypair.id}>{keypair.name}</option>
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
                    <a href="/" className="btn-refresh icon-blue" onClick={handleRefresh} cds-text="secondary">
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
