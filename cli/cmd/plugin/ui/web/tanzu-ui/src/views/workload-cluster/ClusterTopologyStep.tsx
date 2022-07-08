//React imports
import React, { ChangeEvent, useContext } from 'react';
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { SubmitHandler, useForm } from 'react-hook-form';

// Library imports
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { CdsIcon } from '@cds/react/icon';
import { ClarityIcons, computerIcon, cpuIcon, flaskIcon, memoryIcon } from '@cds/core/icon';

// App imports
import './WorkloadClusterWizard.scss';
import { CCVAR_CHANGE } from '../../state-management/actions/Form.actions';
import { ClusterNameSection } from '../../shared/components/FormInputSections/ClusterNameSection';
import { getSelectedManagementCluster } from './WorkloadClusterUtility';
import { isK8sCompliantString } from '../../shared/validations/Validation.service';
import ManagementClusterInfoBanner from './ManagementClusterInfoBanner';
import RadioButton from '../../shared/components/widgets/RadioButton';
import { StepProps } from '../../shared/components/wizard/Wizard';
import { WcStore } from './Store.wc';
import { VSPHERE_FIELDS } from '../management-cluster/vsphere/VsphereManagementCluster.constants';

// NOTE: field names expected to start [category name]___ (because the Form reducer strips that to find the data path)
// that way these field names match the ones in the ClusterAttributeStep and are stored with the same mechanism
const FIELD_NAME_WORKLOAD_CLUSTER_NAME = 'TopologyStep___WORKLOAD_CLUSTER_NAME';
const FIELD_NAME_WORKER_NODE_INSTANCE_TYPE = 'TopologyStep___WORKER_NODE_INSTANCE_TYPE';

interface ClusterTopologyStepFormInputs {
    [FIELD_NAME_WORKLOAD_CLUSTER_NAME]: string;
    [FIELD_NAME_WORKER_NODE_INSTANCE_TYPE]: string;
}

const clusterTopologyStepFormSchema = yup
    .object({
        [FIELD_NAME_WORKLOAD_CLUSTER_NAME]: yup
            .string()
            .nullable()
            .required('Please enter a name for your workload cluster')
            .test(
                '',
                'Cluster name must contain only lower case letters and hyphen',
                (value) => value !== null && isK8sCompliantString(value)
            ),
        [FIELD_NAME_WORKER_NODE_INSTANCE_TYPE]: yup
            .string()
            .nullable()
            .required('Please select an instance type for your workload cluster nodes'),
    })
    .required();

interface WorkerNodeInstanceType {
    id: string;
    icon: string;
    name: string;
    description: string;
}

// NOTE: icons must be imported
const workerNodeInstanceTypes: WorkerNodeInstanceType[] = [
    {
        id: 'basic-demo',
        name: 'Basic demo',
        icon: 'flask',
        description:
            'Virtual machines with a range of compute and memory resources. Intended for small projects and development environments.',
    },
    {
        id: 'general-purpose',
        name: 'General purpose',
        icon: 'computer',
        description:
            'General purpose instances powered by multi-threaded CPUs. Balanced, high performance, compute and memory for' +
            ' production workloads.',
    },
    {
        id: 'compute-optimized',
        name: 'Compute optimized',
        icon: 'cpu',
        description: 'Compute optimized instances suited for CPU-intensive workloads such as CI/CD, machine learning, and data processing.',
    },
    {
        id: 'memory-optimized',
        name: 'Memory optimized',
        icon: 'memory',
        description: 'Memory optimized instances best suited for in-memory operations such as big-data and performant databases.',
    },
];

ClarityIcons.addIcons(flaskIcon, computerIcon, cpuIcon, memoryIcon);

function ClusterTopologyStep(props: Partial<StepProps>) {
    const { handleValueChange, currentStep, goToStep, submitForm } = props;
    const { state } = useContext(WcStore);
    const cluster = getSelectedManagementCluster(state);
    const methods = useForm<ClusterTopologyStepFormInputs>({ resolver: yupResolver(clusterTopologyStepFormSchema) });
    const {
        register,
        handleSubmit,
        formState: { errors },
    } = methods;
    const onSubmit: SubmitHandler<ClusterTopologyStepFormInputs> = (data) => {
        if (Object.keys(errors).length === 0) {
            if (goToStep && currentStep && submitForm) {
                goToStep(currentStep + 1);
                submitForm(currentStep);
            }
        }
    };

    const onClusterNameChange = (clusterName: string | undefined) => {
        if (handleValueChange) {
            handleValueChange(CCVAR_CHANGE, VSPHERE_FIELDS.CLUSTERNAME, clusterName, currentStep, errors, { clusterName: cluster.name });
        }
    };

    const onValueChange = (evt: ChangeEvent<HTMLSelectElement>) => {
        if (handleValueChange) {
            const value = evt.target.value;
            const key = evt.target.name;
            handleValueChange(CCVAR_CHANGE, key, value, currentStep, errors, { clusterName: cluster.name });
        }
    };

    return (
        <div className="wizard-content-container" key="cluster-topology">
            <p cds-text="heading">Workload Topology Settings</p>
            <br />
            {ManagementClusterInfoBanner(cluster)}
            <br />
            <div cds-layout="grid gap:xxl" key="section-holder">
                <div cds-layout="col:6" key="cluster-name-section">
                    {ClusterNameSection(VSPHERE_FIELDS.CLUSTERNAME, errors, register, onClusterNameChange)}
                </div>
                <div cds-layout="col:6" key="instance-type-section">
                    {WorkerNodeInstanceTypeSection(errors, register, onValueChange)}
                </div>
            </div>
            <br />
            <CdsButton onClick={handleSubmit(onSubmit)}>NEXT</CdsButton>
        </div>
    );
}

function WorkerNodeInstanceTypeSection(
    errors: any,
    register: any,
    onSelectNodeInstanceType: (evt: ChangeEvent<HTMLSelectElement>) => void
) {
    return (
        <div cds-layout="vertical gap:lg gap@md:lg col@sm:6 col:6">
            <div cds-layout="cols:6">Select a worker node instance type</div>
            <div cds-layout="grid gap:md align:fill">
                {workerNodeInstanceTypes.map((instanceType) => {
                    return InstanceTypeInList(instanceType, register, onSelectNodeInstanceType);
                })}
            </div>
            {errors[FIELD_NAME_WORKER_NODE_INSTANCE_TYPE] && (
                <CdsControlMessage status="error">{errors[FIELD_NAME_WORKER_NODE_INSTANCE_TYPE].message}</CdsControlMessage>
            )}
        </div>
    );
}

function InstanceTypeInList(
    instance: WorkerNodeInstanceType,
    register: any,
    onSelectNodeInstanceType: (evt: ChangeEvent<HTMLSelectElement>) => void
) {
    return (
        <>
            <div className="text-white" cds-layout="col:1">
                <CdsIcon shape={instance.icon}></CdsIcon>
            </div>
            <RadioButton
                className="input-radio"
                cdsLayout="col:1"
                value={instance.id}
                register={register}
                name={FIELD_NAME_WORKER_NODE_INSTANCE_TYPE}
                onChange={onSelectNodeInstanceType}
            />
            <div className="text-white" cds-layout="col:10">
                {instance.name} {instance.description}
            </div>
        </>
    );
}

export default ClusterTopologyStep;
