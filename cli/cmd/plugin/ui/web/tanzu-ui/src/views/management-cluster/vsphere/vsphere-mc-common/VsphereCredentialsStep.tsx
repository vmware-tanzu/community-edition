// React imports
import React, { ChangeEvent, useCallback, useContext, useEffect, useState } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsCheckbox } from '@cds/react/checkbox';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsIcon } from '@cds/react/icon';
import { CdsInput } from '@cds/react/input';
import { CdsSelect } from '@cds/react/select';
import { CdsToggle } from '@cds/react/toggle';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';

// App imports
import '../VsphereManagementCluster.scss';
import { analyzeOsImages } from './VsphereOsImageUtil';
import ConnectionNotification, { CONNECTION_STATUS } from '../../../../shared/components/ConnectionNotification/ConnectionNotification';
import { createSchema } from './vsphere.credential.form.schema';
import { DefaultOrchestrator } from '../../default-orchestrator/DefaultOrchestrator';
import { FormAction } from '../../../../shared/types/types';
import { initDefaults, initOsImages } from './VsphereOrchestrator.service';
import { INPUT_CHANGE, INPUT_CLEAR } from '../../../../state-management/actions/Form.actions';
import { IP_FAMILIES, VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';
import { isValidFqdn, isValidIp4, isValidIp6 } from '../../../../shared/validations/Validation.service';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
import { ThumbprintDisplay } from './ThumbprintDisplay';
import UseUpdateTabStatus from '../../../../shared/components/wizard/UseUpdateTabStatus.hooks';
import { VSphereCredentials, VSphereDatacenter, VsphereService, VSphereVirtualMachine } from '../../../../swagger-api';
import { VsphereStore } from '../Store.vsphere.mc';

export interface VsphereCredentialsStepInputs {
    [VSPHERE_FIELDS.DATACENTER]: string;
    [VSPHERE_FIELDS.IPFAMILY]: string;
    [VSPHERE_FIELDS.PASSWORD]: string;
    [VSPHERE_FIELDS.SERVERNAME]: string;
    [VSPHERE_FIELDS.USERNAME]: string;
    [VSPHERE_FIELDS.USETHUMBPRINT]: boolean;
}
type VSPHERE_CREDENTIALS_STEP_FIELDS =
    | VSPHERE_FIELDS.DATACENTER
    | VSPHERE_FIELDS.IPFAMILY
    | VSPHERE_FIELDS.PASSWORD
    | VSPHERE_FIELDS.SERVERNAME
    | VSPHERE_FIELDS.USERNAME
    | VSPHERE_FIELDS.USETHUMBPRINT;

const SERVER_RESPONSE_BAD_CREDENTIALS = 403;

export function VsphereCredentialsStep(props: Partial<StepProps>) {
    const { currentStep, goToStep, submitForm, updateTabStatus } = props;
    const { vsphereState, vsphereDispatch } = useContext(VsphereStore);

    const [connectionMessage, setConnectionMessage] = useState('');
    const [connectionStatus, setConnectionStatus] = useState(CONNECTION_STATUS.DISCONNECTED);
    const [datacenters, setDatacenters] = useState<VSphereDatacenter[]>([]);
    const [errorObject, setErrorObject] = useState({});
    const [loadingDatacenters, setLoadingDatacenters] = useState(false);
    const [selectedDatacenter, setSelectedDatacenter] = useState<string>();
    const [osImageMessage, setOsImageMessage] = useState<string>('');
    const [serverNameAtBlur, setServerNameAtBlur] = useState<string>('');
    const [thumbprint, setThumbprint] = useState('');
    const [thumbprintServer, setThumbprintServer] = useState('');
    const [thumbprintErrorMessage, setThumbprintErrorMessage] = useState('');
    const [useThumbprint, setUseThumbprint] = useState(true);
    const [ipFamily, setIpFamily] = useState(vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.IPFAMILY] || IP_FAMILIES.IPv4);
    const [selectedDcHasTemplate, setSelectedDcHasTemplate] = useState<boolean>(false);
    const methods = useForm<VsphereCredentialsStepInputs>({
        resolver: yupResolver(createSchema(ipFamily)),
        mode: 'all',
    });
    const {
        register,
        handleSubmit,
        setValue,
        formState: { errors },
    } = methods;

    // update tab status bar
    if (updateTabStatus) {
        UseUpdateTabStatus(errors, currentStep, updateTabStatus);
    }

    const errDataCenter = () => connectionStatus === CONNECTION_STATUS.CONNECTED && errors[VSPHERE_FIELDS.DATACENTER];
    const errNoDataCentersFound = () => {
        return connectionStatus === CONNECTION_STATUS.CONNECTED && !loadingDatacenters && !datacenters?.length;
    };
    const errDataCenterMsg = () => errors[VSPHERE_FIELDS.DATACENTER]?.message || '';

    const verifyVsphereThumbprint = useCallback(
        (serverName: string) => {
            function isValidServerName(serverName: string): boolean {
                return (
                    isValidFqdn(serverName) ||
                    (ipFamily === IP_FAMILIES.IPv6 && isValidIp6(serverName)) ||
                    (ipFamily === IP_FAMILIES.IPv4 && isValidIp4(serverName))
                );
            }

            function doThumbprintCheck(serverName: string) {
                setThumbprintServer(serverName);
                VsphereService.getVsphereThumbprint(serverName).then(
                    (response) => {
                        console.log(`thumbprint response: ${JSON.stringify(response)}`);
                        vsphereDispatch({
                            type: INPUT_CHANGE,
                            field: VSPHERE_FIELDS.THUMBPRINT,
                            payload: response.thumbprint,
                        } as FormAction);
                    },
                    (reasonRejected) => {
                        setThumbprintErrorMessage(`Unable to obtain thumbprint: ${reasonRejected.message}`);
                    }
                );
            }

            if (isValidServerName(serverName)) {
                doThumbprintCheck(serverName);
            }
        },
        [ipFamily]
    );

    const handleFieldChange = (fieldName: VSPHERE_FIELDS, value: any) => {
        vsphereDispatch({
            type: INPUT_CHANGE,
            field: fieldName,
            payload: value,
        } as FormAction);

        if (fieldName === VSPHERE_FIELDS.SERVERNAME) {
            vsphereDispatch({ type: INPUT_CLEAR, field: VSPHERE_FIELDS.THUMBPRINT } as FormAction);
            setThumbprintErrorMessage('');
        }
    };

    const handleFieldChangeEvent = (event: ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const fieldName = event.target.name as VSPHERE_CREDENTIALS_STEP_FIELDS;
        const value = event.target.value;
        handleFieldChange(fieldName, value);
    };

    const handleDatacenterChange = async (event: ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const newSelectedDatacenter = event.target.value;
        handleFieldChangeEvent(event);
        setSelectedDatacenter(newSelectedDatacenter);
        DefaultOrchestrator.clearResourceData(vsphereDispatch, VSPHERE_FIELDS.VMTEMPLATE);
        setOsImageMessage('');
        setSelectedDcHasTemplate(false);
        if (newSelectedDatacenter) {
            await retrieveOsImages(newSelectedDatacenter);
        }
    };

    const handleCredentialsFieldChange = (event: ChangeEvent<HTMLInputElement>) => {
        setConnectionStatus(CONNECTION_STATUS.DISCONNECTED);
        setConnectionMessage('');
        handleFieldChangeEvent(event);
    };

    const handleCredentialsFieldBlur = (event: ChangeEvent<HTMLInputElement>) => {
        if (event.target.name === VSPHERE_FIELDS.SERVERNAME) {
            setServerNameAtBlur(event.target.value);
        }
    };

    // toggle the value of ipFamily
    const handleIpFamilyChange = (event: ChangeEvent<HTMLInputElement>) => {
        const newValue = event.target['checked'] ? IP_FAMILIES.IPv6 : IP_FAMILIES.IPv4;
        vsphereDispatch({
            type: INPUT_CHANGE,
            field: VSPHERE_FIELDS.IPFAMILY,
            payload: newValue,
        } as FormAction);
        setIpFamily(newValue);
    };

    const handleUseThumbprintChange = (event: ChangeEvent<HTMLInputElement>) => {
        const newValue = event.target['checked'];
        vsphereDispatch({
            type: INPUT_CHANGE,
            field: VSPHERE_FIELDS.USETHUMBPRINT,
            payload: newValue,
        } as FormAction);
        setUseThumbprint(newValue);
    };

    const canContinue = (): boolean => {
        return connectionStatus === CONNECTION_STATUS.CONNECTED && selectedDcHasTemplate;
    };

    const onSubmit: SubmitHandler<VsphereCredentialsStepInputs> = (data) => {
        if (canContinue() && goToStep && currentStep && submitForm) {
            goToStep(currentStep + 1);
            submitForm(currentStep);
        }
    };

    const connectionDataEntered = (): boolean => {
        return (
            !errors[VSPHERE_FIELDS.SERVERNAME] &&
            vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.SERVERNAME] &&
            vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.USERNAME] &&
            vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.PASSWORD]
        );
    };

    const login = () => {
        setConnectionStatus(CONNECTION_STATUS.CONNECTING);
        setConnectionMessage(`Connecting to ${vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.SERVERNAME]}`);
        const vSphereCredentials = {
            username: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.USERNAME],
            host: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.SERVERNAME],
            password: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.PASSWORD],
            thumbprint: useThumbprint ? vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.THUMBPRINT] : '',
            insecure: !useThumbprint,
        } as VSphereCredentials;
        // TODO: remove setTimeout(), which is just here to simulate a backend call delay
        setTimeout(() => {
            VsphereService.setVSphereEndpoint(vSphereCredentials).then(
                () => {
                    setConnectionStatus(CONNECTION_STATUS.CONNECTED);
                    setConnectionMessage(`Connected to ${vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.SERVERNAME]}`);
                },
                (reason: any) => {
                    setConnectionStatus(CONNECTION_STATUS.ERROR);
                    console.warn(`When trying to connect to the server, encountered error ${JSON.stringify(reason)}`);
                    let msg = 'Unable to connect to server! (See console for details.)';
                    if (reason?.status === SERVER_RESPONSE_BAD_CREDENTIALS) {
                        msg = 'Incorrect username/password combination! Please try again. (See console for technical details.)';
                    }
                    setConnectionMessage(msg);
                }
            );
        }, 1600);
    };

    const handleConnect = () => {
        login();
    };

    useEffect(() => {
        setSelectedDatacenter('');
        setDatacenters([]);
        DefaultOrchestrator.clearResourceData(vsphereDispatch, VSPHERE_FIELDS.VMTEMPLATE);
        setOsImageMessage('');
        if (connectionStatus === CONNECTION_STATUS.CONNECTED) {
            setLoadingDatacenters(true);
            VsphereService.getVSphereDatacenters().then((datacenters) => {
                setDatacenters(datacenters);
                setLoadingDatacenters(false);
            });
        } else {
            setLoadingDatacenters(false);
        }
    }, [connectionStatus]);

    // NOTE: this effect is primarily used to get the server thumbprint when the server name change happens (on blur, not every
    // time the user types a character!). However, there is a special case where the user typed in a server name that was erroneous
    // (so we did not get an SSL thumbprint), but by changing the IP FAMILY, the server name has become valid (without actually changing
    // value). So now we want to get the thumbprint of the now-valid server. This is taken care of by having verifyVsphereThumbprint in
    // the dependencies list, because that function is re-assigned whenever ipFamily changes, which will trigger this effect.
    useEffect(() => {
        vsphereDispatch({ type: INPUT_CLEAR, field: VSPHERE_FIELDS.THUMBPRINT } as FormAction);
        setThumbprintErrorMessage('');
        if (serverNameAtBlur) {
            // If the user has already entered a value for the server name, its validity may change when the IP family selection changes.
            // To make sure the user sees that newly-valid or newly-invalid status reflected in the UI,
            // we reset the value (to the same thing) and ask the framework to validate
            setValue(VSPHERE_FIELDS.SERVERNAME, serverNameAtBlur, { shouldTouch: true, shouldValidate: true });
            verifyVsphereThumbprint(serverNameAtBlur);
        }
    }, [setValue, serverNameAtBlur, verifyVsphereThumbprint]);

    useEffect(() => {
        initDefaults(vsphereDispatch);
    }, []);

    return (
        <>
            <div className="wizard-content-container" cds-layout="vertical gap:lg">
                <div cds-layout="grid gap:sm">
                    {Title()}
                    {IntroSection()}
                    {IPFamilySection(ipFamily)}
                    {Credentials()}
                    {ThumbprintVerification()}
                    {ConnectionSection(connectionDataEntered(), connectionStatus, connectionMessage)}
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
                <h2 cds-layout="m-t:md m-b:xl" cds-text="title">
                    vSphere Credentials
                </h2>
            </div>
        );
    }
    function IntroSection() {
        return (
            <div cds-layout="m-b:xs col:8 align:top">
                Provide the vCenter server user credentials to create the Management Cluster on vSphere.
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

    function IPFamilySection(ipFam: IP_FAMILIES) {
        return (
            <div cds-layout="col:4">
                IP family (currently {ipFam})
                <CdsFormGroup layout="vertical">
                    <CdsToggle>
                        <label>Use IPv6</label>
                        <input
                            type="checkbox"
                            {...register(VSPHERE_FIELDS.IPFAMILY, {
                                onChange: handleIpFamilyChange,
                            })}
                        />
                    </CdsToggle>
                </CdsFormGroup>
            </div>
        );
    }

    function Credentials() {
        return (
            <div cds-layout="col:8 p-t:xxs">
                <CdsFormGroup layout="vertical-inline" control-width="shrink">
                    <div cds-layout="horizontal gap:lg align:vertical-center p-b:sm">
                        {CredentialsField('vSphere server', VSPHERE_FIELDS.SERVERNAME, 'vSphere server')}
                        {CredentialsField('Username', VSPHERE_FIELDS.USERNAME, 'username')}
                        {CredentialsField('Password', VSPHERE_FIELDS.PASSWORD, 'password', true)}
                    </div>
                </CdsFormGroup>
            </div>
        );
    }

    function ThumbprintVerification() {
        const box = useThumbprint ? (
            <input type="checkbox" {...register(VSPHERE_FIELDS.USETHUMBPRINT, { onChange: handleUseThumbprintChange })} checked />
        ) : (
            <input type="checkbox" {...register(VSPHERE_FIELDS.USETHUMBPRINT, { onChange: handleUseThumbprintChange })} />
        );
        return (
            <div cds-layout="col:4 vertical gap:md">
                <div>
                    <CdsFormGroup layout="vertical">
                        <CdsCheckbox layout="horizontal">
                            <label>Use SSL thumbprint for secure login</label>
                            {box}
                        </CdsCheckbox>
                    </CdsFormGroup>
                </div>
                <div>
                    <ThumbprintDisplay
                        thumbprint={vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.THUMBPRINT]}
                        errorMessage={thumbprintErrorMessage}
                        serverName={thumbprintServer}
                    />
                </div>
            </div>
        );
    }

    function ConnectionSection(dataEntered: boolean, connectionStatus: CONNECTION_STATUS, connectionMessage: string) {
        return (
            <div cds-layout="col:12">
                <div cds-layout="grid align:vertical-center gap:md">
                    <div cds-layout="col:2">
                        <CdsButton onClick={handleConnect} disabled={connectionStatus === CONNECTION_STATUS.CONNECTED || !dataEntered}>
                            <CdsIcon shape="connect" size="md"></CdsIcon>
                            CONNECT
                        </CdsButton>
                    </div>
                    <div></div>
                    <div cds-layout="col:9 p-b:sm">
                        <ConnectionNotification status={connectionStatus} message={connectionMessage}></ConnectionNotification>
                    </div>
                </div>
            </div>
        );
    }

    function CredentialsField(label: string, fieldName: VSPHERE_CREDENTIALS_STEP_FIELDS, placeholder: string, isPassword = false) {
        const err = errors[fieldName];
        return (
            <CdsInput layout="vertical">
                <label cds-layout="p-b:xs">{label}</label>
                <input
                    {...register(fieldName, {
                        onChange: handleCredentialsFieldChange,
                    })}
                    placeholder={placeholder}
                    type={isPassword ? 'password' : 'text'}
                    onBlurCapture={handleCredentialsFieldBlur}
                    defaultValue={vsphereState[STORE_SECTION_FORM][fieldName]}
                />
                {err && <CdsControlMessage status="error">{err.message}</CdsControlMessage>}
                {!err && <CdsControlMessage status="neutral">&nbsp;</CdsControlMessage>}
            </CdsInput>
        );
    }

    // NOTE: because of the way the "disabled" attribute works, we have two functions to render the Datacenter CdsSelect control.
    // Other more elegant solutions do not appear to work correctly.
    function DatacenterSelectWithDatacenters() {
        return (
            <CdsSelect layout="vertical" controlWidth="shrink">
                <label cds-layout="p-b:xs">Datacenter</label>
                <select {...register(VSPHERE_FIELDS.DATACENTER, { onChange: handleDatacenterChange })}>
                    <option />
                    {datacenters.map((dc) => (
                        <option key={dc.moid}>{dc.name}</option>
                    ))}
                </select>
            </CdsSelect>
        );
    }

    function DatacenterSelectWithoutDatacenters() {
        return (
            <CdsSelect layout="vertical" className="min-width-200">
                <label cds-layout="p-b:xs">Datacenter</label>
                <select {...register(VSPHERE_FIELDS.DATACENTER)} disabled>
                    <option />
                </select>
            </CdsSelect>
        );
    }

    function DatacenterSection() {
        const hasDatacenters = datacenters && datacenters.length > 0;
        return (
            <div cds-layout="grid horizontal gap:md">
                <div cds-layout="col:6">
                    {hasDatacenters && DatacenterSelectWithDatacenters()}
                    {!hasDatacenters && DatacenterSelectWithoutDatacenters()}
                    <div cds-layout="p-t:md">
                        {errNoDataCentersFound() && <CdsControlMessage status="error">No data centers found on server!</CdsControlMessage>}
                        {errDataCenter() && <CdsControlMessage status="error">{errDataCenterMsg()}</CdsControlMessage>}
                        {!errNoDataCentersFound() && !errDataCenter() && <CdsControlMessage status="neutral">&nbsp;</CdsControlMessage>}
                    </div>
                </div>
                <div cds-layout="col:6">
                    {connectionStatus === CONNECTION_STATUS.CONNECTED && selectedDatacenter && osImageMessage && (
                        <div>
                            <div className="error-text">Unable to use this Datacenter unless changes are made</div>
                            <br />
                            <CdsControlMessage status="error">{osImageMessage}</CdsControlMessage>
                            <br />
                            <CdsButton action="outline" onClick={handleRecheckOsImages}>
                                Refresh the OS image check
                            </CdsButton>
                        </div>
                    )}
                </div>
            </div>
        );
    }

    async function retrieveOsImages(datacenter: string | undefined) {
        setOsImageMessage('');
        setSelectedDcHasTemplate(false);
        if (!datacenter) {
            DefaultOrchestrator.clearResourceData(vsphereDispatch, VSPHERE_FIELDS.VMTEMPLATE);
            return;
        }
        const osImages = await initOsImages(datacenter, {
            errorObject,
            setErrorObject,
            vsphereDispatch,
            vsphereState,
        });

        // TODO: get good URL for how to convert OSImage to template
        const { msg, nImages, nTemplates } = analyzeOsImages(datacenter, '', osImages);
        setOsImageMessage(msg);
        setSelectedDcHasTemplate(nTemplates > 0);
        console.log(`After retrieving os images, nImages=${nImages} and nTemplates=${nTemplates}`);
    }

    function handleRecheckOsImages() {
        retrieveOsImages(selectedDatacenter);
    }
}
