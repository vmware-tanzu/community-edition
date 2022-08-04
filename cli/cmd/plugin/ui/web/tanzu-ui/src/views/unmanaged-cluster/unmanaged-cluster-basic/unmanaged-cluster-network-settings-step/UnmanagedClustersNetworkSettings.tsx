// React imports
import React, { useContext, useState } from 'react';

// Library imports
import { blockIcon, blocksGroupIcon, ClarityIcons, clusterIcon } from '@cds/core/icon';
import { CdsSelect } from '@cds/react/select';
import { FormProvider, useForm } from 'react-hook-form';
import { CdsRadio, CdsRadioGroup } from '@cds/react/radio';
import { CdsIcon } from '@cds/react/icon';
import { CdsButton } from '@cds/react/button';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

// App imports
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
import { UnmanagedService } from '../../../../swagger-api/services/UnmanagedService';
import { useNavigate } from 'react-router-dom';
import UseUpdateTabStatus from '../../../../shared/components/wizard/UseUpdateTabStatus.hooks';
import TextInputWithError from '../../../../shared/components/Input/TextInputWithError';
import { FormAction } from '../../../../shared/types/types';

ClarityIcons.addIcons(blockIcon, blocksGroupIcon, clusterIcon);

enum UNMANAGED_NETWORK_FIELDS {
    CLUSTER_SERVICE_CIDR = 'CLUSTER_SERVICE_CIDR',
    CLUSTER_POD_CIDR = 'CLUSTER_POD_CIDR',
    IP_ADDRESS = 'IP_ADDRESS',
    HOST_PORT_MAPPING = 'HOST_PORT_MAPPING',
    NODE_PORT_MAPPING = 'NODE_PORT_MAPPING',
}

interface FormInputs {
    [UNMANAGED_NETWORK_FIELDS.IP_ADDRESS]: string;
    [UNMANAGED_NETWORK_FIELDS.HOST_PORT_MAPPING]: string;
    [UNMANAGED_NETWORK_FIELDS.NODE_PORT_MAPPING]: string;
    [UNMANAGED_NETWORK_FIELDS.CLUSTER_SERVICE_CIDR]: string;
    [UNMANAGED_NETWORK_FIELDS.CLUSTER_POD_CIDR]: string;
}

const unmanagedClusterNetworkSettingStepFormSchema = yup
    .object({
        [UNMANAGED_NETWORK_FIELDS.CLUSTER_SERVICE_CIDR]: yup
            .string()
            .nullable()
            .required('Please enter a CIDR for your cluster service')
            .test('', 'Must be valid CIDR', (value) => value !== null && isValidCidr(value)),
        [UNMANAGED_NETWORK_FIELDS.CLUSTER_POD_CIDR]: yup
            .string()
            .nullable()
            .required('Please enter a CIDER for your cluster pod')
            .test('', 'Must be valid CIDR', (value) => value !== null && isValidCidr(value)),
        [UNMANAGED_NETWORK_FIELDS.IP_ADDRESS]: yup
            .string()
            .nullable()
            .required('Please enter an IP Address for your cluster pod')
            .test('', 'Must be valid IP Address', (value) => value !== null && isValidIp(value)),
    })
    .required();

const unmanagedClusterProviders = [
    {
        label: 'calico',
        value: 'calico',
    },
    {
        label: 'antrea',
        value: 'antrea',
    },
    {
        label: 'none',
        value: 'none',
    },
];

const unmanagedClusterProtocol = [
    {
        label: 'tcp',
        value: 'tcp',
    },
    {
        label: 'udp',
        value: 'udp',
    },
    {
        label: 'sctp',
        value: 'sctp',
    },
];

function UnmanagedClusterNetworkSettings(props: Partial<StepProps>) {
    const { updateTabStatus, currentStep } = props;

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

    const methods = useForm<FormInputs>({ resolver: yupResolver(unmanagedClusterNetworkSettingStepFormSchema), mode: 'all' });
    const {
        formState: { errors },
    } = methods;

    const { umcDispatch } = useContext(UmcStore);

    // update tab status bar
    if (updateTabStatus) {
        UseUpdateTabStatus(errors, currentStep, updateTabStatus);
    }

    const handleFieldChange = (field: string, value: string) => {
        umcDispatch({
            type: INPUT_CHANGE,
            field,
            payload: value,
        } as FormAction);
    };

    const [selectedProvider] = useState(umcState[STORE_SECTION_FORM].CLUSTER_PROVIDER);

    const deployUnmanagedCluster = async () => {
        setConnectionStatus(CONNECTION_STATUS.CONNECTING);
        setConnectionMessage('Attempting to start unmanaged cluster creation.');
        const unmanagedClusterParams: CreateUnmanagedClusterParams = {
            name: umcState[STORE_SECTION_FORM].CLUSTER_NAME,
            provider: K8sProviders.KIND,
            cni: umcState[STORE_SECTION_FORM].CLUSTER_PROVIDER,
            podcidr: umcState[STORE_SECTION_FORM].POD_CIDR,
            servicecidr: umcState[STORE_SECTION_FORM].SERVICE_CIDR,
            portmappings: [
                combineNodeToHostPortMapping(
                    umcState[STORE_SECTION_FORM].IP_ADDRESS,
                    umcState[STORE_SECTION_FORM].NODE_PORT_MAPPING,
                    umcState[STORE_SECTION_FORM].HOST_PORT_MAPPING,
                    umcState[STORE_SECTION_FORM].CLUSTER_PROTOCOL
                ),
            ],
            controlplanecount: umcState[STORE_SECTION_FORM].CONTROL_PLANE_NODES_COUNT,
            workernodecount: umcState[STORE_SECTION_FORM].WORKER_NODES_COUNT,
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
        <FormProvider {...methods}>
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
        </FormProvider>
    );

    function ClusterProvider() {
        return (
            <CdsRadioGroup layout="vertical-inline" onChange={(e: any) => handleFieldChange('CLUSTER_PROVIDER', e.target.value)}>
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
        return (
            <div cds-layout="horizontal gap:xl">
                <div cds-layout="col:4">
                    <TextInputWithError
                        label="Cluster service CIDR"
                        name="CLUSTER_SERVICE_CIDR"
                        handleInputChange={handleFieldChange}
                        placeholder="CLUSTER SERVICE CIDR"
                        defaultValue={umcState[STORE_SECTION_FORM].CLUSTER_SERVICE_CIDR}
                    />
                </div>
                <div cds-layout="col:4">
                    <TextInputWithError
                        label="Cluster POD CIDR"
                        name="CLUSTER_POD_CIDR"
                        handleInputChange={handleFieldChange}
                        placeholder="CLUSTER POD CIDR"
                        defaultValue={umcState[STORE_SECTION_FORM].CLUSTER_POD_CIDR}
                    />
                </div>
            </div>
        );
    }

    function NodeHostPortMapping() {
        //TODO: refactor in the future const errorNodeHost = errors[UNMANAGED_NETWORK_FIELDS.NODE_HOST_PORT_MAPPING];
        return (
            <div cds-layout="grid">
                <div cds-layout="col:10">
                    <label cds-text="body extrabold" cds-layout="p-t:sm" slot="label">
                        Node to host port mapping <CdsIcon shape="info-circle" size="md" status="info"></CdsIcon>
                    </label>
                    <div cds-layout="grid gap:sm p-y:md">
                        <div cds-layout="col:2">
                            <TextInputWithError
                                label="IP Address"
                                name={UNMANAGED_NETWORK_FIELDS.IP_ADDRESS}
                                handleInputChange={handleFieldChange}
                                placeholder="127.0.0.1"
                                defaultValue={umcState[STORE_SECTION_FORM].IP_ADDRESS}
                            />
                        </div>
                        <div cds-layout="col:2">
                            <TextInputWithError
                                label="Node Port"
                                name={UNMANAGED_NETWORK_FIELDS.NODE_PORT_MAPPING}
                                handleInputChange={handleFieldChange}
                                placeholder="80"
                                defaultValue={umcState[STORE_SECTION_FORM].NODE_PORT_MAPPING}
                            />
                        </div>
                        <div cds-layout="col:2">
                            <TextInputWithError
                                label="Host Port"
                                name={UNMANAGED_NETWORK_FIELDS.HOST_PORT_MAPPING}
                                handleInputChange={handleFieldChange}
                                placeholder="80"
                                defaultValue={umcState[STORE_SECTION_FORM].HOST_PORT_MAPPING}
                            />
                        </div>
                        <CdsSelect cds-layout="align:shrink" onChange={(e: any) => handleFieldChange('CLUSTER_PROTOCOL', e.target.value)}>
                            <label>Protocol</label>
                            <select>
                                {unmanagedClusterProtocol.map((unmanagedClusterProtocol, index) => {
                                    return (
                                        <option cds-layout="m:md m-l:none" key={index}>
                                            {unmanagedClusterProtocol.label}
                                        </option>
                                    );
                                })}
                            </select>
                        </CdsSelect>
                    </div>
                    {combineNodeToHostPortMapping(
                        umcState[STORE_SECTION_FORM].IP_ADDRESS,
                        umcState[STORE_SECTION_FORM].NODE_PORT_MAPPING,
                        umcState[STORE_SECTION_FORM].HOST_PORT_MAPPING,
                        umcState[STORE_SECTION_FORM].CLUSTER_PROTOCOL
                    )}
                </div>
            </div>
        );
    }
}

export default UnmanagedClusterNetworkSettings;
