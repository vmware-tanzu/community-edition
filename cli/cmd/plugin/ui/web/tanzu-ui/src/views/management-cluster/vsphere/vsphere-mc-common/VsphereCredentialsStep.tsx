// React imports
import React, { ChangeEvent, useContext, useEffect, useState } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';
// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsIcon } from '@cds/react/icon';
import { CdsInput } from '@cds/react/input';
import { CdsSelect } from '@cds/react/select';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';
// App imports
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';
import { vsphereCredentialFormSchema } from './vsphere.credential.form.schema';
import { VSphereCredentials, VSphereDatacenter, VsphereService } from '../../../../swagger-api';
import { VsphereStore } from '../Store.vsphere.mc';

export interface FormInputs {
    [VSPHERE_FIELDS.SERVERNAME]: string;
    [VSPHERE_FIELDS.USERNAME]: string;
    [VSPHERE_FIELDS.PASSWORD]: string;
    [VSPHERE_FIELDS.DATACENTER]: string;
}

const SERVER_RESPONSE_BAD_CREDENTIALS = 403;

export function VsphereCredentialsStep(props: Partial<StepProps>) {
    const { vsphereState } = useContext(VsphereStore);
    const { handleValueChange, currentStep, goToStep, submitForm } = props;
    const [connected, setConnection] = useState(false);
    const [datacenters, setDatacenters] = useState([] as VSphereDatacenter[]);
    const [loadingDatacenters, setLoadingDatacenters] = useState(false);
    const [connectionErrorMessage, setConnectionErrorMessage] = useState('');
    const methods = useForm<FormInputs>({
        resolver: yupResolver(vsphereCredentialFormSchema),
    });

    const {
        register,
        handleSubmit,
        setValue,
        formState: { errors },
    } = methods;

    const errDataCenter = () => connected && errors[VSPHERE_FIELDS.DATACENTER];
    const errNoDataCentersFound = () => {
        return connected && !loadingDatacenters && !datacenters?.length;
    };
    const errDataCenterMsg = () => errors[VSPHERE_FIELDS.DATACENTER]?.message || '';

    const handleFieldChange = (event: ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const fieldName = event.target.name as VSPHERE_FIELDS;
        const value = event.target.value;
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, fieldName, value, currentStep, errors);
            setValue(fieldName, value, { shouldValidate: true });
        }
    };

    const handleCredentialsFieldChange = (event: ChangeEvent<HTMLInputElement>) => {
        setConnection(false);
        setConnectionErrorMessage('');
        handleFieldChange(event);
    };

    const canContinue = (): boolean => {
        return connected && Object.keys(errors).length === 0 && vsphereState.data[VSPHERE_FIELDS.DATACENTER];
    };

    const onSubmit: SubmitHandler<FormInputs> = (data) => {
        if (canContinue() && goToStep && currentStep && submitForm) {
            goToStep(currentStep + 1);
            submitForm(currentStep);
        }
    };

    const connectionDataEntered = (): boolean => {
        return (
            !errors[VSPHERE_FIELDS.SERVERNAME] &&
            vsphereState.data[VSPHERE_FIELDS.SERVERNAME] &&
            vsphereState.data[VSPHERE_FIELDS.USERNAME] &&
            vsphereState.data[VSPHERE_FIELDS.PASSWORD]
        );
    };

    const handleConnect = () => {
        // TODO: use verifyThumbprint
        setConnectionErrorMessage('');
        setConnection(false);
        const vSphereCredentials = {
            username: vsphereState.data[VSPHERE_FIELDS.USERNAME],
            host: vsphereState.data[VSPHERE_FIELDS.SERVERNAME],
            password: vsphereState.data[VSPHERE_FIELDS.PASSWORD],
        } as VSphereCredentials;
        VsphereService.setVSphereEndpoint(vSphereCredentials).then(
            () => {
                setConnection(true);
            },
            (reason: any) => {
                console.warn(`When trying to connect to the server, encountered error ${JSON.stringify(reason)}`);
                let msg = 'Unable to connect to server! (See console for details.)';
                if (reason?.status === SERVER_RESPONSE_BAD_CREDENTIALS) {
                    msg = 'Incorrect username/password combination! Please try again. (See console for technical details.)';
                }
                setConnectionErrorMessage(msg);
            }
        );
    };

    useEffect(() => {
        if (connected) {
            setLoadingDatacenters(true);
            VsphereService.getVSphereDatacenters().then((datacenters) => {
                setDatacenters(datacenters);
                setLoadingDatacenters(false);
            });
        } else {
            setDatacenters([]);
            setLoadingDatacenters(false);
            delete errors[VSPHERE_FIELDS.DATACENTER];
        }
    }, [connected]);

    // TODO: add IP family, thumbprint verification, datacenter
    return (
        <div>
            <div className="wizard-content-container">
                {IntroSection()}
                <div cds-layout="p-t:lg">
                    <CdsFormGroup layout="vertical-inline" control-width="shrink">
                        <div cds-layout="horizontal gap:lg align:vertical-center">
                            {CredentialsField('vSphere server', VSPHERE_FIELDS.SERVERNAME, 'vSphere server')}
                            {CredentialsField('Username', VSPHERE_FIELDS.USERNAME, 'username')}
                            {CredentialsField('Password', VSPHERE_FIELDS.PASSWORD, 'password', true)}
                        </div>
                    </CdsFormGroup>
                    {VerticalSpacer()}
                    {ConnectionSection(connectionDataEntered(), connected, connectionErrorMessage)}
                </div>
                {VerticalSpacer()}
                {DatacenterSection()}
                {VerticalSpacer()}
                <CdsButton onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                    NEXT
                </CdsButton>
            </div>
        </div>
    );

    function VerticalSpacer() {
        return (
            <div>
                <br />
            </div>
        );
    }

    function IntroSection() {
        return (
            <>
                <h2 cds-layout="m-t:lg">vSphere Credentials</h2>
                <p cds-layout="m-y:lg" className="description">
                    Provide the vCenter server user credentials to create the Management Servicer on vSphere.
                </p>
                <p cds-layout="m-y:lg" className="description">
                    Don&apos;t have vSphere credentials? View our guide on{' '}
                    <a href="/" className="text-blue">
                        creating vSphere credentials
                    </a>
                    .
                </p>
            </>
        );
    }

    function ConnectionSection(dataEntered: boolean, isConnected: boolean, errMessage: string) {
        return (
            <>
                <div>
                    {errMessage && (
                        <div>
                            <CdsControlMessage status="error">{errMessage}</CdsControlMessage>
                            <br />
                        </div>
                    )}
                    {isConnected && (
                        <div>
                            <CdsControlMessage status="success">Connection established</CdsControlMessage>
                            <br />
                        </div>
                    )}
                    {!errMessage && !isConnected && (
                        <div>
                            <CdsControlMessage status="neutral">&nbsp;</CdsControlMessage>
                            <br />
                        </div>
                    )}
                </div>
                <CdsButton onClick={handleConnect} disabled={isConnected || !dataEntered}>
                    <CdsIcon shape="connect" size="md"></CdsIcon>
                    {isConnected ? 'CONNECTED' : 'CONNECT'}
                </CdsButton>
            </>
        );
    }

    function CredentialsField(label: string, fieldName: VSPHERE_FIELDS, placeholder: string, isPassword = false) {
        const err = errors[fieldName];
        return (
            <CdsInput layout="compact">
                <label cds-layout="p-b:md">{label}</label>
                <input
                    {...register(fieldName)}
                    placeholder={placeholder}
                    type={isPassword ? 'password' : 'text'}
                    onChange={handleCredentialsFieldChange}
                    defaultValue={vsphereState.data[fieldName]}
                />
                {err && <CdsControlMessage status="error">{err.message}</CdsControlMessage>}
                {!err && <CdsControlMessage status="neutral">&nbsp;</CdsControlMessage>}
            </CdsInput>
        );
    }

    function DatacenterSection() {
        return (
            <>
                <CdsSelect layout="vertical">
                    <label cds-layout="p-b:md">Datacenter</label>
                    <select
                        {...register(VSPHERE_FIELDS.DATACENTER)}
                        onChange={handleFieldChange}
                        disabled={!connected || !datacenters || datacenters.length === 0}
                    >
                        <option />
                        {datacenters.map((dc) => (
                            <option key={dc.moid}>{dc.name}</option>
                        ))}
                    </select>
                </CdsSelect>
                <div>
                    <br />
                </div>

                {errNoDataCentersFound() && <CdsControlMessage status="error">No data centers found on server!</CdsControlMessage>}
                {errDataCenter() && <CdsControlMessage status="error">{errDataCenterMsg()}</CdsControlMessage>}
                {!errNoDataCentersFound() && !errDataCenter() && <CdsControlMessage status="neutral">&nbsp;</CdsControlMessage>}
            </>
        );
    }
}
