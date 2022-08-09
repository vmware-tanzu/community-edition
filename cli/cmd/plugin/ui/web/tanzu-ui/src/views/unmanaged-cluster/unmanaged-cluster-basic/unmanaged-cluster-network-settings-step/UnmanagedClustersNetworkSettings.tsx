// React imports
import React, { useContext, useState } from 'react';

// Library imports
import { blockIcon, blocksGroupIcon, ClarityIcons, clusterIcon } from '@cds/core/icon';
import { CdsSelect } from '@cds/react/select';
import { FormProvider, useForm } from 'react-hook-form';
import { CdsRadio, CdsRadioGroup } from '@cds/react/radio';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

// App imports
import { CdsIcon } from '@cds/react/icon';
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage } from '@cds/react/forms';
import { ClusterProtocols } from '../../../../shared/constants/ClusterProtocols.constants';
import { CniProviders } from '../../../../shared/constants/CniProviders.constants';
import ConnectionNotification, { CONNECTION_STATUS } from '../../../../shared/components/ConnectionNotification/ConnectionNotification';
import { CreateUnmanagedClusterParams } from '../../../../swagger-api';
import { DeploymentStates, DeploymentTypes } from '../../../../shared/constants/Deployment.constants';
import { DEPLOYMENT_STATUS_CHANGED } from '../../../../state-management/actions/Deployment.actions';
import { FormAction } from '../../../../shared/types/types';
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
import UseUpdateTabStatus from '../../../../shared/components/wizard/UseUpdateTabStatus.hooks';
import TextInputWithError from '../../../../shared/components/Input/TextInputWithError';

ClarityIcons.addIcons(blockIcon, blocksGroupIcon, clusterIcon);

interface FormInputs {
    [UNMANAGED_CLUSTER_FIELDS.IP_ADDRESS]: string;
    [UNMANAGED_CLUSTER_FIELDS.HOST_PORT_MAPPING]: string;
    [UNMANAGED_CLUSTER_FIELDS.NODE_PORT_MAPPING]: string;
    [UNMANAGED_CLUSTER_FIELDS.CLUSTER_SERVICE_CIDR]: string;
    [UNMANAGED_CLUSTER_FIELDS.CLUSTER_POD_CIDR]: string;
}

function createYupSchemaObject() {
    return {
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
    };
}
const unmanagedClusterNetworkSettingStepFormSchema = yup.object(createYupSchemaObject()).required();

const unmanagedCniProviders = [
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
        <FormProvider {...methods}>
            <div className="cluster-settings-container" cds-layout="m:lg">
                <div cds-layout="p-b:lg" cds-text="title">
                    Network settings
                </div>
                <div cds-layout="grid">
                    <div cds-layout="col:12">
                        <div cds-layout="vertical gap:lg">
                            {CniProvider()}
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
                                        <CdsIcon shape="cluster" size="sm" data-testid="create-cluster-btn"></CdsIcon>
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

    function CniProvider() {
        return (
            <CdsRadioGroup
                layout="vertical-inline"
                onChange={(e: any) => handleFieldChange(UNMANAGED_CLUSTER_FIELDS.CNI_PROVIDER, e.target.value)}
            >
                <label>
                    Container Network Interface (CNI) provider <CdsIcon shape="info-circle" size="md" status="info"></CdsIcon>
                </label>
                {unmanagedCniProviders.map((unmanagedCniProviders, index) => {
                    return (
                        <CdsRadio cds-layout="m-t:md" key={index}>
                            <label>{unmanagedCniProviders.label}</label>
                            <input
                                type="radio"
                                key={index}
                                value={unmanagedCniProviders.value}
                                checked={
                                    umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CNI_PROVIDER] === unmanagedCniProviders.value
                                }
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
            <div cds-layout="horizontal gap:xl p-y:md">
                <div cds-layout="col:4">
                    <TextInputWithError
                        label="Cluster service CIDR"
                        name={UNMANAGED_CLUSTER_FIELDS.CLUSTER_SERVICE_CIDR}
                        handleInputChange={handleFieldChange}
                        placeholder={UNMANAGED_PLACEHOLDER_VALUES.CLUSTER_SERVICE_CIDR}
                        defaultValue={umcState[STORE_SECTION_FORM].CLUSTER_SERVICE_CIDR}
                    />
                </div>
                <div cds-layout="col:4">
                    <TextInputWithError
                        label="Cluster POD CIDR"
                        name={UNMANAGED_CLUSTER_FIELDS.CLUSTER_POD_CIDR}
                        handleInputChange={handleFieldChange}
                        placeholder={UNMANAGED_PLACEHOLDER_VALUES.CLUSTER_POD_CIDR}
                        defaultValue={umcState[STORE_SECTION_FORM].CLUSTER_POD_CIDR}
                    />
                </div>
            </div>
        );
    }

    function NodeHostPortMapping() {
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
                                name={UNMANAGED_CLUSTER_FIELDS.IP_ADDRESS}
                                handleInputChange={handleFieldChange}
                                placeholder={UNMANAGED_PLACEHOLDER_VALUES.IP_ADDRESS}
                                defaultValue={umcState[STORE_SECTION_FORM].IP_ADDRESS}
                            />
                        </div>
                        <div cds-layout="col:2">
                            <TextInputWithError
                                label="Node Port"
                                name={UNMANAGED_CLUSTER_FIELDS.NODE_PORT_MAPPING}
                                handleInputChange={handleFieldChange}
                                placeholder={UNMANAGED_PLACEHOLDER_VALUES.NODE_PORT_MAPPING}
                                defaultValue={umcState[STORE_SECTION_FORM].NODE_PORT_MAPPING}
                            />
                        </div>
                        <div cds-layout="col:2">
                            <TextInputWithError
                                label="Host Port"
                                name={UNMANAGED_CLUSTER_FIELDS.HOST_PORT_MAPPING}
                                handleInputChange={handleFieldChange}
                                placeholder={UNMANAGED_PLACEHOLDER_VALUES.HOST_PORT_MAPPING}
                                defaultValue={umcState[STORE_SECTION_FORM].HOST_PORT_MAPPING}
                            />
                        </div>
                        <CdsSelect
                            cds-layout="align:shrink"
                            onChange={(e: any) => handleFieldChange(UNMANAGED_CLUSTER_FIELDS.CLUSTER_PROTOCOL, e.target.value)}
                        >
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
                    <CdsControlMessage>
                        {combineNodeToHostPortMapping(
                            umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.IP_ADDRESS],
                            umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.NODE_PORT_MAPPING],
                            umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.HOST_PORT_MAPPING],
                            umcState[STORE_SECTION_FORM][UNMANAGED_CLUSTER_FIELDS.CLUSTER_PROTOCOL]
                        )}
                    </CdsControlMessage>
                </div>
            </div>
        );
    }
}

export default UnmanagedClusterNetworkSettings;
