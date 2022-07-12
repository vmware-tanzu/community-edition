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
import { CONNECTION_STATUS, ConnectionNotification } from '../../../../shared/components/ConnectionNotification/ConnectionNotification';
import { createSchema } from './vsphere.credential.form.schema';
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { IP_FAMILIES, VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';
import { isValidFqdn, isValidIp4, isValidIp6 } from '../../../../shared/validations/Validation.service';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
import { ThumbprintDisplay } from './ThumbprintDisplay';
import { VSPHERE_ADD_RESOURCES } from '../../../../state-management/actions/Resources.actions';
import { VSphereCredentials, VSphereDatacenter, VsphereService, VSphereVirtualMachine } from '../../../../swagger-api';
import { VsphereResourceAction } from '../../../../shared/types/types';
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
    const { handleValueChange, currentStep, goToStep, submitForm } = props;
    const { vsphereState, vsphereDispatch } = useContext(VsphereStore);

    const [connectionMessage, setConnectionMessage] = useState('');
    const [connectionStatus, setConnectionStatus] = useState(CONNECTION_STATUS.DISCONNECTED);
    const [datacenters, setDatacenters] = useState<VSphereDatacenter[]>([]);
    const [loadingDatacenters, setLoadingDatacenters] = useState(false);
    const [selectedDatacenter, setSelectedDatacenter] = useState<string>();
    const [osImages, setOsImages] = useState<VSphereVirtualMachine[]>([]);
    const [osImageMessage, setOsImageMessage] = useState<string>('');
    const [serverNameAtBlur, setServerNameAtBlur] = useState<string>('');
    const [thumbprint, setThumbprint] = useState('');
    const [thumbprintServer, setThumbprintServer] = useState('');
    const [thumbprintErrorMessage, setThumbprintErrorMessage] = useState('');
    const [useThumbprint, setUseThumbprint] = useState(true);
    const [ipFamily, setIpFamily] = useState(vsphereState[VSPHERE_FIELDS.IPFAMILY] || IP_FAMILIES.IPv4);
    const [selectedDcHasTemplate, setSelectedDcHasTemplate] = useState<boolean>(false);
    const methods = useForm<VsphereCredentialsStepInputs>({
        resolver: yupResolver(createSchema(ipFamily)),
    });
    const {
        register,
        handleSubmit,
        setValue,
        formState: { errors },
    } = methods;

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
                        setThumbprint(response.thumbprint || '');
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

    const handleFieldChange = (event: ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const fieldName = event.target.name as VSPHERE_CREDENTIALS_STEP_FIELDS;
        const value = event.target.value;
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, fieldName, value, currentStep, errors);
            setValue(fieldName, value, { shouldValidate: true });
        }
        if (fieldName === VSPHERE_FIELDS.SERVERNAME) {
            setThumbprint('');
            setThumbprintErrorMessage('');
        }
    };

    const handleDatacenterChange = (event: ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const newSelectedDatacenter = event.target.value;
        handleFieldChange(event);
        setSelectedDatacenter(newSelectedDatacenter);
        setOsImages([]);
        setOsImageMessage('');
        setSelectedDcHasTemplate(false);
        if (newSelectedDatacenter) {
            retrieveOsImages(newSelectedDatacenter);
        }
    };

    const handleCredentialsFieldChange = (event: ChangeEvent<HTMLInputElement>) => {
        setConnectionStatus(CONNECTION_STATUS.DISCONNECTED);
        setConnectionMessage('');
        handleFieldChange(event);
    };

    const handleCredentialsFieldBlur = (event: ChangeEvent<HTMLInputElement>) => {
        if (event.target.name === VSPHERE_FIELDS.SERVERNAME) {
            setServerNameAtBlur(event.target.value);
        }
    };

    // toggle the value of ipFamily
    const handleIpFamilyChange = (event: ChangeEvent<HTMLInputElement>) => {
        const newValue = event.target['checked'] ? IP_FAMILIES.IPv6 : IP_FAMILIES.IPv4;
        handleValueChange && handleValueChange(INPUT_CHANGE, VSPHERE_FIELDS.IPFAMILY, newValue, currentStep, errors);
        setIpFamily(newValue);
    };

    const handleUseThumbprintChange = (event: ChangeEvent<HTMLInputElement>) => {
        const newValue = event.target['checked'];
        handleValueChange && handleValueChange(INPUT_CHANGE, VSPHERE_FIELDS.USETHUMBPRINT, newValue, currentStep, errors);
        setUseThumbprint(newValue);
    };

    const canContinue = (): boolean => {
        return connectionStatus === CONNECTION_STATUS.CONNECTED && selectedDcHasTemplate;
    };

    const onSubmit: SubmitHandler<VsphereCredentialsStepInputs> = (data) => {
        if (canContinue() && goToStep && currentStep && submitForm) {
            goToStep(currentStep + 1);
            submitForm(currentStep);
            vsphereDispatch({
                type: VSPHERE_ADD_RESOURCES,
                datacenter: selectedDatacenter,
                resourceName: 'osImages',
                payload: osImages,
            } as VsphereResourceAction);
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
            thumbprint: useThumbprint ? thumbprint : '',
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
        setOsImages([]);
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
        setThumbprint('');
        setThumbprintErrorMessage('');
        if (serverNameAtBlur) {
            // If the user has already entered a value for the server name, its validity may change when the IP family selection changes.
            // To make sure the user sees that newly-valid or newly-invalid status reflected in the UI,
            // we reset the value (to the same thing) and ask the framework to validate
            setValue(VSPHERE_FIELDS.SERVERNAME, serverNameAtBlur, { shouldTouch: true, shouldValidate: true });
            verifyVsphereThumbprint(serverNameAtBlur);
        }
    }, [setValue, serverNameAtBlur, verifyVsphereThumbprint]);

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
                <h2>vSphere Credentials</h2>
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
                        <input type="checkbox" {...register(VSPHERE_FIELDS.IPFAMILY)} onChange={handleIpFamilyChange} />
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
            <input type="checkbox" {...register(VSPHERE_FIELDS.USETHUMBPRINT)} onChange={handleUseThumbprintChange} checked />
        ) : (
            <input type="checkbox" {...register(VSPHERE_FIELDS.USETHUMBPRINT)} onChange={handleUseThumbprintChange} />
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
                    <ThumbprintDisplay thumbprint={thumbprint} errorMessage={thumbprintErrorMessage} serverName={thumbprintServer} />
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
                    <div cds-layout="col:9 p-b:sm">{ConnectionNotification(connectionStatus, connectionMessage)}</div>
                </div>
            </div>
        );
    }

    function CredentialsField(label: string, fieldName: VSPHERE_CREDENTIALS_STEP_FIELDS, placeholder: string, isPassword = false) {
        const err = errors[fieldName];
        return (
            <CdsInput layout="compact">
                <label cds-layout="p-b:xs">{label}</label>
                <input
                    {...register(fieldName)}
                    placeholder={placeholder}
                    type={isPassword ? 'password' : 'text'}
                    onChange={handleCredentialsFieldChange}
                    onBlurCapture={handleCredentialsFieldBlur}
                    defaultValue={vsphereState[STORE_SECTION_FORM][fieldName]}
                />
                {err && <CdsControlMessage status="error">{err.message}</CdsControlMessage>}
                {!err && <CdsControlMessage status="neutral">&nbsp;</CdsControlMessage>}
            </CdsInput>
        );
    }

    function DatacenterSection() {
        return (
            <div cds-layout="grid horizontal gap:md">
                <div cds-layout="col:6">
                    <CdsSelect layout="vertical" controlWidth="shrink">
                        <label cds-layout="p-b:xs">Datacenter</label>
                        <select
                            {...register(VSPHERE_FIELDS.DATACENTER)}
                            onChange={handleDatacenterChange}
                            disabled={connectionStatus !== CONNECTION_STATUS.CONNECTED || !datacenters || datacenters.length === 0}
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

    function retrieveOsImages(datacenter: string | undefined) {
        setOsImages([]);
        setSelectedDcHasTemplate(false);
        if (datacenter) {
            VsphereService.getVSphereOsImages(datacenter).then((fetchedOsImages) => {
                setOsImages(fetchedOsImages);
                // TODO: get good URL for how to convert OSImage to template
                const { msg, nImages, nTemplates } = analyzeOsImages(datacenter, 'URL', fetchedOsImages);
                setOsImageMessage(msg);
                setSelectedDcHasTemplate(nTemplates > 0);
                console.log(`After retrieving os images, nImages=${nImages} and nTemplates=${nTemplates}`);
            });
        }
    }

    function handleRecheckOsImages() {
        retrieveOsImages(selectedDatacenter);
    }
}
