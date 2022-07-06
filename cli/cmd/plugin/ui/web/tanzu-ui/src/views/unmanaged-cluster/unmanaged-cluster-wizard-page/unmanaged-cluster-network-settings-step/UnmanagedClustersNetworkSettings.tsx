// React imports
import React, { ChangeEvent, useState, useContext } from 'react';

// Library imports
import { ClarityIcons, blockIcon, blocksGroupIcon, clusterIcon } from '@cds/core/icon';
import { CdsControlMessage } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';
import { CdsSelect } from '@cds/react/select';
import { useForm } from 'react-hook-form';
import { CdsRadio, CdsRadioGroup } from '@cds/react/radio';
import { CdsIcon } from '@cds/react/icon';
import { CdsButton } from '@cds/react/button';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

// App imports
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { UmcStore } from '../../../../state-management/stores/Store.umc';
import { isValidCidr, isValidIp } from '../../../../shared/validations/Validation.service';
import { STORE_SECTION_FORM } from '../../../../state-management/reducers/Form.reducer';

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
        value: 'CALICO',
    },
    {
        label: 'antrea',
        value: 'ANTREA',
    },
    {
        label: 'none',
        value: 'NONE',
    },
];

const unmanagedClusterProtocol = [
    {
        label: 'tcp',
        value: 'TCP',
    },
    {
        label: 'udp',
        value: 'UDP',
    },
    {
        label: 'sctp',
        value: 'SCTP',
    },
];

function UnmanagedClusterNetworkSettings(props: Partial<StepProps>) {
    const { handleValueChange, currentStep } = props;

    const { umcState } = useContext(UmcStore);

    const {
        register,
        formState: { errors },
    } = useForm<FormInputs>({ resolver: yupResolver(unmanagedClusterNetworkSettingStepFormSchema) });

    const handleFieldChange = (event: ChangeEvent<HTMLInputElement>) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, event.target.name, event.target.value, currentStep, errors);
        }
    };

    const [selectedProvider, setSelectedProvider] = useState(umcState[STORE_SECTION_FORM].CLUSTER_PROVIDER);

    const handleProviderChange = (event: ChangeEvent<HTMLSelectElement>) => {
        setSelectedProvider(event.target.value);
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, 'CLUSTER_PROVIDER', event.target.value, currentStep, errors);
        }
    };

    const handleProtocolChange = (event: ChangeEvent<HTMLSelectElement>) => {
        if (handleValueChange) {
            handleValueChange(INPUT_CHANGE, 'CLUSTER_PROTOCOL', event.target.value, currentStep, errors);
        }
    };

    return (
        <div className="cluster-settings-container" cds-layout="m:lg">
            <div cds-layout="p-b:lg" cds-text="title">
                Network settings
            </div>
            <div cds-layout="grid">
                <div cds-layout="col@sm:8">
                    <div cds-layout="vertical gap:lg">
                        {ClusterProvider()}
                        {ClusterCidr()}
                        {NodeHostPortMapping()}
                        <CdsButton cds-layout="col:start-1" status="success">
                            <CdsIcon shape="cluster" size="sm"></CdsIcon>
                            Create Unmanaged cluster
                        </CdsButton>
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
        return (
            <div cds-layout="horizontal gap:xl">
                <div cds-layout="col:4">
                    <CdsInput layout="vertical" control-width="shrink">
                        <label>
                            Cluster service CIDR <CdsIcon shape="info-circle" size="md" status="info"></CdsIcon>
                        </label>
                        <input
                            {...register(UNMANAGED_NETWORK_FIELDS.CLUSTER_SERVICE_CIDR)}
                            placeholder="CLUSTER SERVICE CIDR"
                            onChange={handleFieldChange}
                            defaultValue="100.64.0.0/13"
                        ></input>
                        {errors['CLUSTER_SERVICE_CIDR'] && (
                            <CdsControlMessage status="error">{errors['CLUSTER_SERVICE_CIDR'].message}</CdsControlMessage>
                        )}
                    </CdsInput>
                </div>
                <div cds-layout="col:4">
                    <CdsInput layout="vertical" control-width="shrink">
                        <label>
                            Cluster pod CIDR <CdsIcon shape="info-circle" size="md" status="info"></CdsIcon>
                        </label>
                        <input
                            {...register(UNMANAGED_NETWORK_FIELDS.CLUSTER_POD_CIDR)}
                            placeholder="Cluster POD CIDR"
                            onChange={handleFieldChange}
                            defaultValue={'100.96.0.0/11'}
                        ></input>
                        {errors['CLUSTER_POD_CIDR'] && (
                            <CdsControlMessage status="error">{errors['CLUSTER_POD_CIDR'].message}</CdsControlMessage>
                        )}
                    </CdsInput>
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
                        <CdsInput layout="vertical" control-width="shrink" cds-layout="col:2">
                            <label>IP Address</label>
                            <input
                                {...register(UNMANAGED_NETWORK_FIELDS.IP_ADDRESS)}
                                placeholder="127.0.0.1"
                                onChange={handleFieldChange}
                                defaultValue={umcState[STORE_SECTION_FORM].IP_ADDRESS}
                            ></input>
                        </CdsInput>
                        <CdsInput layout="vertical" control-width="shrink" cds-layout="col:2">
                            <label>Host Port</label>
                            <input
                                {...register(UNMANAGED_NETWORK_FIELDS.HOST_PORT_MAPPING)}
                                placeholder="80"
                                onChange={handleFieldChange}
                                defaultValue={umcState[STORE_SECTION_FORM].HOST_PORT_MAPPING}
                            ></input>
                        </CdsInput>
                        <CdsInput layout="vertical" control-width="shrink" cds-layout="col:2">
                            <label>Node Port</label>
                            <input
                                {...register(UNMANAGED_NETWORK_FIELDS.NODE_PORT_MAPPING)}
                                placeholder="80"
                                onChange={handleFieldChange}
                                defaultValue={umcState[STORE_SECTION_FORM].NODE_PORT_MAPPING}
                            ></input>
                        </CdsInput>
                        <CdsSelect cds-layout="align:shrink" onChange={handleProtocolChange}>
                            <label>Protocol</label>
                            <select>
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
                    <CdsControlMessage>{`${umcState[STORE_SECTION_FORM].IP_ADDRESS}/${umcState[STORE_SECTION_FORM].HOST_PORT_MAPPING}/${umcState[STORE_SECTION_FORM].NODE_PORT_MAPPING}/${umcState[STORE_SECTION_FORM].CLUSTER_PROTOCOL}`}</CdsControlMessage>
                </div>
            </div>
        );
    }
}

export default UnmanagedClusterNetworkSettings;
