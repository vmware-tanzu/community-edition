// React imports
import React, { ChangeEvent, useContext, useEffect, useState } from 'react';
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
import { createSchema } from './vsphere.credential.form.schema';
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { IPFAMILIES, VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';
import { isValidFqdn, isValidIp4, isValidIp6 } from '../../../../shared/validations/Validation.service';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
import { ThumbprintDisplay } from './ThumbprintDisplay';
import { VSPHERE_ADD_RESOURCES, VSPHERE_DELETE_RESOURCES } from '../../../../state-management/actions/Resources.actions';
import { VSphereCredentials, VSphereDatacenter, VsphereService, VSphereVirtualMachine } from '../../../../swagger-api';
import { VsphereResourceAction } from '../../../../shared/types/types';
import { VsphereStore } from '../Store.vsphere.mc';

export interface FormInputs {
    [VSPHERE_FIELDS.DATACENTER]: string;
    [VSPHERE_FIELDS.IPFAMILY]: string;
    [VSPHERE_FIELDS.PASSWORD]: string;
    [VSPHERE_FIELDS.SERVERNAME]: string;
    [VSPHERE_FIELDS.USERNAME]: string;
    [VSPHERE_FIELDS.USETHUMBPRINT]: boolean;
}

const SERVER_RESPONSE_BAD_CREDENTIALS = 403;

export function VsphereCredentialsStep(props: Partial<StepProps>) {
    const { handleValueChange, currentStep, goToStep, submitForm } = props;
    const { vsphereState, vsphereDispatch } = useContext(VsphereStore);

    const [connected, setConnection] = useState(false);
    const [connectionErrorMessage, setConnectionErrorMessage] = useState('');
    const [datacenters, setDatacenters] = useState<VSphereDatacenter[]>([]);
    const [loadingDatacenters, setLoadingDatacenters] = useState(false);
    const [selectedDatacenter, setSelectedDatacenter] = useState<string>();
    const [dcOsImages, setDcOsImages] = useState<VSphereVirtualMachine[]>([]);
    const [osImageMessage, setOsImageMessage] = useState<string>('');
    const [loadingOsImages, setLoadingOsImages] = useState(false);
    const [thumbprint, setThumbprint] = useState('');
    const [thumbprintServer, setThumbprintServer] = useState('');
    const [thumbprintErrorMessage, setThumbprintErrorMessage] = useState('');
    const [useThumbprint, setUseThumbprint] = useState(true);
    const [ipFamily, setIpFamily] = useState(vsphereState[VSPHERE_FIELDS.IPFAMILY] || IPFAMILIES.IPv4);
    const [selectedDcHasTemplate, setSelectedDcHasTemplate] = useState<boolean>(false);
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

    const clearOsImages = (oldDatacenter: string | undefined) => {
        setDcOsImages([]);
        if (oldDatacenter) {
            vsphereDispatch({
                type: VSPHERE_DELETE_RESOURCES,
                resourceName: 'osImages',
                datacenter: oldDatacenter,
            } as VsphereResourceAction);
        }
    };

    const clearDatacenters = () => {
        setDatacenters([]); // important to clear data centers before clearing os images, or os image error msg will be set
        if (selectedDatacenter) {
            clearOsImages(selectedDatacenter);
        }
        setSelectedDatacenter('');
        handleValueChange && handleValueChange(INPUT_CHANGE, VSPHERE_FIELDS.DATACENTER, '', currentStep, errors);
        delete errors[VSPHERE_FIELDS.DATACENTER];
    };

    const handleFieldChange = (event: ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const fieldName = event.target.name as VSPHERE_FIELDS;
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
        const oldSelectedDatacenter = selectedDatacenter;
        const newSelectedDatacenter = event.target.value;
        handleFieldChange(event);
        setSelectedDatacenter(newSelectedDatacenter);
        if (newSelectedDatacenter) {
            retrieveOsImages(newSelectedDatacenter, oldSelectedDatacenter);
        } else {
            // clear values related to (possible) previous datacenter selection
            clearOsImages(oldSelectedDatacenter);
        }
    };

    const handleCredentialsFieldChange = (event: ChangeEvent<HTMLInputElement>) => {
        setConnection(false);
        setConnectionErrorMessage('');
        handleFieldChange(event);
    };

    const handleCredentialsFieldBlur = (event: ChangeEvent<HTMLInputElement>) => {
        if (event.target.name === VSPHERE_FIELDS.SERVERNAME) {
            verifyVsphereThumbprint(event.target.value);
        }
    };

    // toggle the value of ipFamily
    const handleIpFamilyChange = (event: ChangeEvent<HTMLInputElement>) => {
        const newValue = event.target['checked'] ? IPFAMILIES.IPv6 : IPFAMILIES.IPv4;
        handleValueChange && handleValueChange(INPUT_CHANGE, VSPHERE_FIELDS.IPFAMILY, newValue, currentStep, errors);
        setIpFamily(newValue);
    };

    const handleUseThumbprintChange = (event: ChangeEvent<HTMLInputElement>) => {
        const newValue = event.target['checked'];
        handleValueChange && handleValueChange(INPUT_CHANGE, VSPHERE_FIELDS.USETHUMBPRINT, newValue, currentStep, errors);
        setUseThumbprint(newValue);
    };

    const canContinue = (): boolean => {
        return connected && selectedDcHasTemplate;
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
            vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.SERVERNAME] &&
            vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.USERNAME] &&
            vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.PASSWORD]
        );
    };

    const login = () => {
        const vSphereCredentials = {
            username: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.USERNAME],
            host: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.SERVERNAME],
            password: vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.PASSWORD],
            thumbprint: useThumbprint ? thumbprint : '',
            insecure: !useThumbprint,
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

    const isValidServerName = (serverName: string): boolean => {
        return (
            isValidFqdn(serverName) ||
            (ipFamily === IPFAMILIES.IPv6 && isValidIp6(serverName)) ||
            (ipFamily === IPFAMILIES.IPv4 && isValidIp4(serverName))
        );
    };

    const verifyVsphereThumbprint = (serverName: string) => {
        if (isValidServerName(serverName)) {
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
    };

    const handleConnect = () => {
        setConnectionErrorMessage('');
        setConnection(false);
        login();
    };

    useEffect(() => {
        clearDatacenters();
        if (connected) {
            setLoadingDatacenters(true);
            VsphereService.getVSphereDatacenters().then((datacenters) => {
                setDatacenters(datacenters);
                setLoadingDatacenters(false);
            });
        } else {
            setLoadingDatacenters(false);
        }
    }, [connected]);

    useEffect(() => {
        // If the user has entered a value for the server name, its validity will change when the IP family selection changes.
        // To make sure the user sees that reflected in the UI, we reset the value (to the same thing) and ask the framework to validate
        const existingServerValue = getValues(VSPHERE_FIELDS.SERVERNAME);
        if (existingServerValue) {
            setValue(VSPHERE_FIELDS.SERVERNAME, existingServerValue, { shouldTouch: true, shouldValidate: true });
        }
        // There is a special case where the user typed in a server name that was erroneous (so we did not get SSL thumbprint),
        // but by changing the IP FAMILY, the server name has become valid (without actually changing value). So now we want to get the
        // thumbprint of the now-valid server
        verifyVsphereThumbprint(vsphereState[STORE_SECTION_FORM][VSPHERE_FIELDS.SERVERNAME]);
    }, [ipFamily]);

    useEffect(() => {
        const nOsImages = dcOsImages?.length;
        const nTemplates = dcOsImages.reduce<number>((accum, image) => accum + (image.isTemplate ? 1 : 0), 0);
        if (nOsImages === 0) {
            // There may be no OS images because no data center has been selected (or because the selected dc has no OS images)
            setOsImageMessage(
                datacenters && selectedDatacenter && !loadingOsImages
                    ? `No OS images are available! Please select a different data center or add an OS image to ${selectedDatacenter}`
                    : ''
            );
            setSelectedDcHasTemplate(false);
        } else if (nTemplates === 0) {
            const describeNumTemplates = nOsImages === 1 ? 'There is one OS image' : `There are ${nOsImages} OS images`;
            const notATemplate = nOsImages === 1 ? 'it is not a template' : 'none of them are templates';
            // TODO: get URL for how to convert an OS image to a template
            setOsImageMessage(
                `${describeNumTemplates} on data center ${selectedDatacenter}, but ${notATemplate}.` +
                    `For information on how to convert an OS image to a template, see URL`
            );
            setSelectedDcHasTemplate(false);
        } else {
            setOsImageMessage('');
            setSelectedDcHasTemplate(true);
        }
        if (loadingOsImages) {
            console.log('Loading OS images...');
        } else {
            console.log(
                `There are ${nOsImages} OS images on datacenter ${selectedDatacenter}, of which ${nTemplates} ${
                    nTemplates === 1 ? 'is a template' : 'are templates'
                }`
            );
        }
    }, [dcOsImages, loadingOsImages]);

    return (
        <>
            <div className="wizard-content-container" cds-layout="vertical gap:lg">
                <div cds-layout="grid gap:sm">
                    {Title()}
                    {IntroSection()}
                    {IPFamilySection(ipFamily)}
                    {Credentials()}
                    {ThumbprintVerification()}
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
            <div cds-layout="m-b:xs col:8 align:top">
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

    function ConnectionSection(dataEntered: boolean, isConnected: boolean, errMessage: string) {
        return (
            <div cds-layout="col:12">
                <div cds-layout="grid align:vertical-center gap:md">
                    <div cds-layout="col:2">
                        <CdsButton onClick={handleConnect} disabled={isConnected || !dataEntered}>
                            <CdsIcon shape="connect" size="md"></CdsIcon>
                            {isConnected ? 'CONNECTED' : 'CONNECT'}
                        </CdsButton>
                    </div>
                    <div></div>
                    <div cds-layout="col:9 p-b:sm">
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
                <div cds-layout="col:6">
                    {osImageMessage && (
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

    function retrieveOsImages(newDatacenter: string | undefined, prevDatacenter: string | undefined) {
        setLoadingOsImages(true);
        clearOsImages(prevDatacenter);

        if (newDatacenter) {
            VsphereService.getVSphereOsImages(newDatacenter).then((osImages) => {
                setDcOsImages(osImages);
                setLoadingOsImages(false);
                vsphereDispatch({
                    type: VSPHERE_ADD_RESOURCES,
                    datacenter: newDatacenter,
                    resourceName: 'osImages',
                    payload: osImages,
                } as VsphereResourceAction);
            });
        }
    }

    function handleRecheckOsImages() {
        retrieveOsImages(selectedDatacenter, selectedDatacenter);
    }
}
