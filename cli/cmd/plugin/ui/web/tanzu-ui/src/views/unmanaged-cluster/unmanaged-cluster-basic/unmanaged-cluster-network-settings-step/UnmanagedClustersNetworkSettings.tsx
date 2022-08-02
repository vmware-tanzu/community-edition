// React imports
import React, { ChangeEvent, useContext, useState } from 'react';
import { useForm } from 'react-hook-form';

// Library imports
import { blockIcon, blocksGroupIcon, ClarityIcons, clusterIcon } from '@cds/core/icon';
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage } from '@cds/react/forms';
import { CdsIcon } from '@cds/react/icon';
import { CdsInput } from '@cds/react/input';
import { CdsRadio, CdsRadioGroup } from '@cds/react/radio';
import { CdsSelect } from '@cds/react/select';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

// App imports
import { ClusterProtocols } from '../../../../shared/constants/ClusterProtocols.constants';
import { CniProviders } from '../../../../shared/constants/CniProviders.constants';
import ConnectionNotification, { CONNECTION_STATUS } from '../../../../shared/components/ConnectionNotification/ConnectionNotification';
import { CreateUnmanagedClusterParams } from '../../../../swagger-api';
import { DeploymentStates, DeploymentTypes } from '../../../../shared/constants/Deployment.constants';
import { DEPLOYMENT_STATUS_CHANGED } from '../../../../state-management/actions/Deployment.actions';
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { isValidCidr, isValidIp } from '../../../../shared/validations/Validation.service';
import { K8sProviders } from '../../../../shared/constants/K8sProviders.constants';
import { NavRoutes } from '../../../../shared/constants/NavRoutes.constants';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { Store } from '../../../../state-management/stores/Store';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';
import { UmcStore } from '../../../../state-management/stores/Store.umc';
import { UNMANAGED_CLUSTER_FIELDS } from '../../unmanaged-cluster-common/UnmanagedCluster.constants';
import { UNMANAGED_PLACEHOLDER_VALUES } from '../../unmanaged-cluster-common/unmanaged.defaults';
import { UnmanagedService } from '../../../../swagger-api/services/UnmanagedService';
import { useNavigate } from 'react-router-dom';

ClarityIcons.addIcons(blockIcon, blocksGroupIcon, clusterIcon);

interface FormInputs {
    [UNMANAGED_CLUSTER_FIELDS.IP_ADDRESS]: string;
    [UNMANAGED_CLUSTER_FIELDS.HOST_PORT_MAPPING]: string;
    [UNMANAGED_CLUSTER_FIELDS.NODE_PORT_MAPPING]: string;
    [UNMANAGED_CLUSTER_FIELDS.CLUSTER_SERVICE_CIDR]: string;
    [UNMANAGED_CLUSTER_FIELDS.CLUSTER_POD_CIDR]: string;
}

const unmanagedClusterNetworkSettingStepFormSchema = yup
    .object({
        [UNMANAGED_CLUSTER_FIELDS.CLUSTER_SERVICE_CIDR]: yup
            .string()
            .nullable()
            .required('Please enter a CIDR for your cluster service')
            .test('', 'Must be valid CIDR', (value) => value !== null && isValidCidr(value)),
        [UNMANAGED_CLUSTER_FIELDS.CLUSTER_POD_CIDR]: yup
            .string()
            .nullable()
            .required('Please enter a CIDER for your cluster pod')
            .test('', 'Must be valid CIDR', (value) => value !== null && isValidCidr(value)),
        [UNMANAGED_CLUSTER_FIELDS.IP_ADDRESS]: yup
            .string()
            .nullable()
            .required('Please enter an IP Address for your cluster pod')
            .test('', 'Must be valid IP Address', (value) => value !== null && isValidIp(value)),
    })
    .required();

const unmanagedClusterProviders = [
    {
        label: CniProviders.CALICO,
        value: CniProviders.CALICO,
    },
    {
        label: CniProviders.ANTREA,
        value: CniProviders.ANTREA,
    },
    {
        label: CniProviders.NONE,
        value: CniProviders.NONE,
    },
];

const unmanagedClusterProtocol = [
    {
        label: ClusterProtocols.TCP,
        value: ClusterProtocols.TCP,
    },
    {
        label: ClusterProtocols.UDP,
        value: ClusterProtocols.UDP,
    },
    {
        label: ClusterProtocols.SCTP,
        value: ClusterProtocols.SCTP,
    },
];

function UnmanagedClusterNetworkSettings(props: Partial<StepProps>) {
    const { handleValueChange, currentStep } = props;

    const { dispatch } = useContext(Store);
    const { umcState } = useContext(UmcStore);

    const [connectionMessage, setConnectionMessage] = useState('');
    const [connectionStatus, setConnectionStatus] = useState(CONNECTION_STATUS.DISCONNECTED);

    const navigate = useNavigate();
    const navigateToProgress = (): void => {
        navigate('/' + NavRoutes.DEPLOY_PROGRESS);
    };

    function combineNodeToHostPortMapping(ipAddress: string, nodePort: string, hostPort: string, protocol: string) {
        return `${ipAddress}:${nodePort}:${hostPort}/${protocol}`;
    }

    const {
        register,
        formState: { errors },
    } = useForm<FormInputs>({ resolver: yupResolver(unmanagedClusterNetworkSettingStepFormSchema) });

    const handleFieldChange = (event: ChangeEvent<HTMLInputElement>) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, event.target.name, event.target.value, currentStep, errors);
        }
    };

    const [selectedProvider, setSelectedProvider] = useState(umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CNI_PROVIDER]);

    const handleProviderChange = (event: ChangeEvent<HTMLSelectElement>) => {
        setSelectedProvider(event.target.value);
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, UNMANAGED_CLUSTER_FIELDS.CNI_PROVIDER, event.target.value, currentStep, errors);
        }
    };

    const handleProtocolChange = (event: ChangeEvent<HTMLSelectElement>) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, UNMANAGED_CLUSTER_FIELDS.CLUSTER_PROTOCOL, event.target.value, currentStep, errors);
        }
    };

    const deployUnmanagedCluster = async () => {
        setConnectionStatus(CONNECTION_STATUS.CONNECTING);
        setConnectionMessage('Attempting to start unmanaged cluster creation.');
        const unmanagedClusterParams: CreateUnmanagedClusterParams = {
            name: umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CLUSTER_NAME],
            provider: K8sProviders.KIND,
            cni: umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CLUSTER_PROVIDER],
            podcidr: umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CLUSTER_POD_CIDR],
            servicecidr: umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CLUSTER_SERVICE_CIDR],
            portmappings: [
                combineNodeToHostPortMapping(
                    umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.IP_ADDRESS],
                    umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.NODE_PORT_MAPPING],
                    umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.HOST_PORT_MAPPING],
                    umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CLUSTER_PROTOCOL]
                ),
            ],
            controlplanecount: umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CONTROL_PLANE_NODE_COUNT],
            workernodecount: umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.WORKER_NODE_COUNT],
        };

        try {
            await UnmanagedService.createUnmanagedCluster(unmanagedClusterParams);

            dispatch({
                type: DEPLOYMENT_STATUS_CHANGED,
                payload: {
                    type: DeploymentTypes.UNMANAGED_CLUSTER,
                    status: DeploymentStates.RUNNING,
                },
            });

            setConnectionStatus(CONNECTION_STATUS.DISCONNECTED);
            navigateToProgress();
        } catch (e) {
            setConnectionStatus(CONNECTION_STATUS.ERROR);
            const msg = 'Error starting unmanaged cluster creation. Please see browser console for more details.';
            console.warn(`Error calling unmanaged cluster create API: ${e}`);
            setConnectionMessage(msg);
        }
    };

    return (
        <div className="cluster-settings-container" cds-layout="m:lg">
            <div cds-layout="p-b:lg" cds-text="title">
                Network settings
            </div>
            <div cds-layout="grid">
                <div cds-layout="col:12">
                    <div cds-layout="vertical gap:lg">
                        {ClusterProvider()}
                        {ClusterCidr()}
                        {NodeHostPortMapping()}
                        <div cds-layout="grid align:vertical-center gap:md">
                            <div cds-layout="col:3">
                                <CdsButton
                                    cds-layout="col:start-1"
                                    status="success"
                                    onClick={deployUnmanagedCluster}
                                    disabled={connectionStatus === CONNECTION_STATUS.CONNECTING}
                                >
                                    <CdsIcon shape="cluster" size="sm"></CdsIcon>
                                    Create Unmanaged cluster
                                </CdsButton>
                            </div>
                            <div></div>
                            <div cds-layout="col:8 p-b:sm">
                                <ConnectionNotification status={connectionStatus} message={connectionMessage} />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );

    function ClusterProvider() {
        return (
            <CdsRadioGroup layout="vertical-inline" onChange={handleProviderChange}>
                <label>
                    Container Network Interface (CNI) provider <CdsIcon shape="info-circle" size="md" status="info"></CdsIcon>
                </label>
                {unmanagedClusterProviders.map((unmanagedClusterProviders, index) => {
                    return (
                        <CdsRadio cds-layout="m:md m-l:none" key={index}>
                            <label>{unmanagedClusterProviders.label}</label>
                            <input
                                type="radio"
                                key={index}
                                value={unmanagedClusterProviders.value}
                                checked={selectedProvider === unmanagedClusterProviders.value}
                                readOnly
                            />
                        </CdsRadio>
                    );
                })}
            </CdsRadioGroup>
        );
    }

    function ClusterCidr() {
        const errorClusterServiceCidr = errors[UNMANAGED_CLUSTER_FIELDS.CLUSTER_SERVICE_CIDR];
        const errorClusterPodCidr = errors[UNMANAGED_CLUSTER_FIELDS.CLUSTER_POD_CIDR];

        return (
            <div cds-layout="grid gap:xl wrap:none m-b:md">
                <div cds-layout="col:3">
                    <CdsInput layout="vertical">
                        <label>
                            Cluster service CIDR <CdsIcon shape="info-circle" size="md" status="info"></CdsIcon>
                        </label>
                        <input
                            cds-layout="p:none"
                            {...register(UNMANAGED_CLUSTER_FIELDS.CLUSTER_SERVICE_CIDR)}
                            placeholder={UNMANAGED_PLACEHOLDER_VALUES[UNMANAGED_CLUSTER_FIELDS.CLUSTER_SERVICE_CIDR]}
                            onChange={handleFieldChange}
                            defaultValue={umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CLUSTER_SERVICE_CIDR]}
                        ></input>
                        {errorClusterServiceCidr && <CdsControlMessage status="error">{errorClusterServiceCidr.message}</CdsControlMessage>}
                    </CdsInput>
                </div>
                <div cds-layout="col:3">
                    <CdsInput layout="vertical">
                        <label>
                            Cluster pod CIDR <CdsIcon shape="info-circle" size="md" status="info"></CdsIcon>
                        </label>
                        <input
                            cds-layout="p:none"
                            {...register(UNMANAGED_CLUSTER_FIELDS.CLUSTER_POD_CIDR)}
                            placeholder={UNMANAGED_PLACEHOLDER_VALUES[UNMANAGED_CLUSTER_FIELDS.CLUSTER_POD_CIDR]}
                            onChange={handleFieldChange}
                            defaultValue={umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CLUSTER_POD_CIDR]}
                        ></input>
                        {errorClusterPodCidr && <CdsControlMessage status="error">{errorClusterPodCidr.message}</CdsControlMessage>}
                    </CdsInput>
                </div>
            </div>
        );
    }

    function NodeHostPortMapping() {
        //TODO: refactor in the future const errorNodeHost = errors[UNMANAGED_NETWORK_FIELDS.NODE_HOST_PORT_MAPPING];
        const errorIpAddress = errors[UNMANAGED_CLUSTER_FIELDS.IP_ADDRESS];
        const errorNodePortMapping = errors[UNMANAGED_CLUSTER_FIELDS.NODE_PORT_MAPPING];
        const errorHostPortMapping = errors[UNMANAGED_CLUSTER_FIELDS.HOST_PORT_MAPPING];

        return (
            <div cds-layout="m-b:md">
                <div cds-layout="horizontal gap:sm p-y:md wrap:none">
                    <CdsInput layout="vertical" control-width="shrink" cds-layout="align:shrink">
                        <label cds-layout="wrap:none">IP Address</label>
                        <input
                            cds-layout="p:none"
                            {...register(UNMANAGED_CLUSTER_FIELDS.IP_ADDRESS)}
                            placeholder={UNMANAGED_PLACEHOLDER_VALUES[UNMANAGED_CLUSTER_FIELDS.IP_ADDRESS]}
                            onChange={handleFieldChange}
                            defaultValue={umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.IP_ADDRESS]}
                        ></input>
                        {errorIpAddress && <CdsControlMessage status="error">{errorIpAddress.message}</CdsControlMessage>}
                    </CdsInput>
                    <CdsInput layout="vertical" control-width="shrink">
                        <label>Node Port</label>
                        <input
                            cds-layout="p:none"
                            {...register(UNMANAGED_CLUSTER_FIELDS.NODE_PORT_MAPPING)}
                            placeholder={UNMANAGED_PLACEHOLDER_VALUES[UNMANAGED_CLUSTER_FIELDS.NODE_PORT_MAPPING]}
                            onChange={handleFieldChange}
                            defaultValue={umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.NODE_PORT_MAPPING]}
                        ></input>
                        {errorNodePortMapping && <CdsControlMessage status="error">{errorNodePortMapping.message}</CdsControlMessage>}
                    </CdsInput>
                    <CdsInput layout="vertical" control-width="shrink">
                        <label>Host Port</label>
                        <input
                            cds-layout="p:none"
                            {...register(UNMANAGED_CLUSTER_FIELDS.HOST_PORT_MAPPING)}
                            placeholder={UNMANAGED_PLACEHOLDER_VALUES[UNMANAGED_CLUSTER_FIELDS.HOST_PORT_MAPPING]}
                            onChange={handleFieldChange}
                            defaultValue={umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.HOST_PORT_MAPPING]}
                        ></input>
                        {errorHostPortMapping && <CdsControlMessage status="error">{errorHostPortMapping.message}</CdsControlMessage>}
                    </CdsInput>
                    <CdsSelect layout="vertical" control-width="shrink" onChange={handleProtocolChange}>
                        <label>Protocol</label>
                        <select cds-layout="p:none">
                            {unmanagedClusterProtocol.map((unmanagedClusterProtocol, index) => {
                                return (
                                    <option cds-layout="m:md m-l:none" key={index}>
                                        <label>{unmanagedClusterProtocol.label}</label>
                                    </option>
                                );
                            })}
                        </select>
                    </CdsSelect>
                </div>
                <CdsControlMessage>
                    {combineNodeToHostPortMapping(
                        umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.IP_ADDRESS],
                        umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.NODE_PORT_MAPPING],
                        umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.HOST_PORT_MAPPING],
                        umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CLUSTER_PROTOCOL]
                    )}
                </CdsControlMessage>
            </div>
        );
    }
}

export default UnmanagedClusterNetworkSettings;
