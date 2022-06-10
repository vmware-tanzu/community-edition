// React imports
import React, { ChangeEvent, useContext, useEffect, useState } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';
// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsIcon } from '@cds/react/icon';
import { CdsInput } from '@cds/react/input';
import { CdsSelect } from '@cds/react/select';
import { CdsToggle } from '@cds/react/toggle';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';
// App imports
import { createSchema } from './vsphere.credential.form.schema';
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { IPFAMILIES, VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { VSphereCredentials, VSphereDatacenter, VsphereService } from '../../../../swagger-api';
import { VsphereStore } from '../Store.vsphere.mc';

export interface FormInputs {
    [VSPHERE_FIELDS.DATACENTER]: string;
    [VSPHERE_FIELDS.IPFAMILY]: string;
    [VSPHERE_FIELDS.PASSWORD]: string;
    [VSPHERE_FIELDS.SERVERNAME]: string;
    [VSPHERE_FIELDS.USERNAME]: string;
}

const SERVER_RESPONSE_BAD_CREDENTIALS = 403;

export function VsphereCredentialsStep(props: Partial<StepProps>) {
    const { handleValueChange, currentStep, goToStep, submitForm } = props;
    const { vsphereState } = useContext(VsphereStore);

    const [connected, setConnection] = useState(false);
    const [connectionErrorMessage, setConnectionErrorMessage] = useState('');
    const [datacenters, setDatacenters] = useState([] as VSphereDatacenter[]);
    const [loadingDatacenters, setLoadingDatacenters] = useState(false);
    const [ipFamily, setIpFamily] = useState(vsphereState.data[VSPHERE_FIELDS.IPFAMILY] || IPFAMILIES.IPv4);
    const methods = useForm<FormInputs>({
        resolver: yupResolver(createSchema(ipFamily)),
    });
    const {
        register,
        handleSubmit,
        setValue,
        formState: { errors },
        getValues,
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

    // toggle the value of ipFamily
    const handleIpFamilyChange = (event: ChangeEvent<HTMLInputElement>) => {
        const newValue = event.target['checked'] ? IPFAMILIES.IPv6 : IPFAMILIES.IPv4;
        handleValueChange && handleValueChange(INPUT_CHANGE, VSPHERE_FIELDS.IPFAMILY, newValue, currentStep, errors);
        setIpFamily(newValue);
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

    useEffect(() => {
        // If the user has entered a value for the server name, its validity will change when the IP family selection changes.
        // To make sure the user sees that reflected in the UI, we reset the value (to the same thing) and ask the framework to validate
        const existingServerValue = getValues(VSPHERE_FIELDS.SERVERNAME);
        if (existingServerValue) {
            setValue(VSPHERE_FIELDS.SERVERNAME, existingServerValue, { shouldTouch: true, shouldValidate: true });
        }
    }, [ipFamily]);

    // TODO: add thumbprint verification
    return (
        <>
            <div className="wizard-content-container" cds-layout="vertical gap:lg">
                <div cds-layout="grid gap:sm">
                    {Title()}
                    {IntroSection()}
                    {IPFamilySection(ipFamily)}
                </div>
                <div>
                    <CdsFormGroup layout="vertical-inline" control-width="shrink">
                        <div cds-layout="horizontal gap:lg align:vertical-center p-b:sm">
                            {CredentialsField('vSphere server', VSPHERE_FIELDS.SERVERNAME, 'vSphere server')}
                            {CredentialsField('Username', VSPHERE_FIELDS.USERNAME, 'username')}
                            {CredentialsField('Password', VSPHERE_FIELDS.PASSWORD, 'password', true)}
                        </div>
                    </CdsFormGroup>
                    {ConnectionSection(connectionDataEntered(), connected, connectionErrorMessage)}
                </div>
                {DatacenterSection()}
                <div>
                    <CdsButton onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                        NEXT
                    </CdsButton>
                </div>
            </div>
        </>
    );

    function Title() {
        return (
            <div cds-layout="col:12">
                <h2>vSphere Credentials</h2>
            </div>
        );
    }
    function IntroSection() {
        return (
            <div cds-layout="m-b:xs col:9 align:top">
                Provide the vCenter server user credentials to create the Management Servicer on vSphere.
                <p cds-layout="m-t:lg" className="description">
                    Don&apos;t have vSphere credentials? View our guide on{' '}
                    <a href="/" className="text-blue">
                        creating vSphere credentials
                    </a>
                    .
                </p>
            </div>
        );
    }

    function IPFamilySection(ipFam: IPFAMILIES) {
        return (
            <div cds-layout="col:3">
                IP family (currently {ipFam})
                <CdsFormGroup layout="vertical">
                    <CdsToggle>
                        <label>Use IPv6</label>
                        <input type="checkbox" {...register(VSPHERE_FIELDS.IPFAMILY)} onChange={handleIpFamilyChange} />
                    </CdsToggle>
                </CdsFormGroup>
            </div>
        );
    }

    function ConnectionSection(dataEntered: boolean, isConnected: boolean, errMessage: string) {
        return (
            <div cds-layout="grid align:vertical-center gap:md">
                <div cds-layout="col:2">
                    <CdsButton onClick={handleConnect} disabled={isConnected || !dataEntered}>
                        <CdsIcon shape="connect" size="md"></CdsIcon>
                        {isConnected ? 'CONNECTED' : 'CONNECT'}
                    </CdsButton>
                </div>
                <div cds-layout="col:10 p-b:sm">
                    {errMessage && (
                        <div>
                            <CdsControlMessage status="error">{errMessage}</CdsControlMessage>
                        </div>
                    )}
                    {isConnected && (
                        <div>
                            <CdsControlMessage status="success">Connection established</CdsControlMessage>
                        </div>
                    )}
                    {!errMessage && !isConnected && (
                        <div>
                            <CdsControlMessage status="neutral">&nbsp;</CdsControlMessage>
                        </div>
                    )}
                </div>
            </div>
        );
    }

    function CredentialsField(label: string, fieldName: VSPHERE_FIELDS, placeholder: string, isPassword = false) {
        const err = errors[fieldName];
        return (
            <CdsInput layout="compact">
                <label cds-layout="p-b:xs">{label}</label>
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
            <div>
                <CdsSelect layout="vertical" controlWidth="shrink">
                    <label cds-layout="p-b:xs">Datacenter</label>
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
                <div cds-layout="p-t:md">
                    {errNoDataCentersFound() && <CdsControlMessage status="error">No data centers found on server!</CdsControlMessage>}
                    {errDataCenter() && <CdsControlMessage status="error">{errDataCenterMsg()}</CdsControlMessage>}
                    {!errNoDataCentersFound() && !errDataCenter() && <CdsControlMessage status="neutral">&nbsp;</CdsControlMessage>}
                </div>
            </div>
        );
    }
}
