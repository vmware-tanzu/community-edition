//React imports
import React, { useContext } from 'react';
import { CdsButton } from '@cds/react/button';
import { SubmitHandler, useForm } from 'react-hook-form';

// Library imports
import { ClarityIcons, computerIcon, cpuIcon, flaskIcon, memoryIcon } from '@cds/core/icon';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

// App imports
import { ClusterName, clusterNameValidation } from '../../shared/components/FormInputComponents/ClusterName/ClusterName';
import { getSelectedManagementCluster } from './WorkloadClusterUtility';
import ManagementClusterInfoBanner from './ManagementClusterInfoBanner';
import {
    NodeInstanceType,
    nodeInstanceTypeValidation,
    NodeProfile,
} from '../../shared/components/FormInputComponents/NodeProfile/NodeProfile';
import { StepProps } from '../../shared/components/wizard/Wizard';
import { VSPHERE_FIELDS } from '../management-cluster/vsphere/VsphereManagementCluster.constants';
import { WcStore } from './Store.wc';
import './WorkloadClusterWizard.scss';

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
        [FIELD_NAME_WORKLOAD_CLUSTER_NAME]: clusterNameValidation(),
        [FIELD_NAME_WORKER_NODE_INSTANCE_TYPE]: nodeInstanceTypeValidation(),
    })
    .required();

// NOTE: icons must be imported
const workerNodeInstanceTypes: NodeInstanceType[] = [
    {
        id: 'basic-demo',
        label: 'Basic demo',
        icon: 'flask',
        description:
            'Virtual machines with a range of compute and memory resources. Intended for small projects and development environments.',
    },
    {
        id: 'general-purpose',
        label: 'General purpose',
        icon: 'computer',
        description:
            'General purpose instances powered by multi-threaded CPUs. Balanced, high performance, compute and memory for' +
            ' production workloads.',
    },
    {
        id: 'compute-optimized',
        label: 'Compute optimized',
        icon: 'cpu',
        description: 'Compute optimized instances suited for CPU-intensive workloads such as CI/CD, machine learning, and data processing.',
    },
    {
        id: 'memory-optimized',
        label: 'Memory optimized',
        icon: 'memory',
        description: 'Memory optimized instances best suited for in-memory operations such as big-data and performant databases.',
    },
];

ClarityIcons.addIcons(flaskIcon, computerIcon, cpuIcon, memoryIcon);

function ClusterTopologyStep(props: Partial<StepProps>) {
    const { currentStep, goToStep, submitForm } = props;
    const { state } = useContext(WcStore);
    const cluster = getSelectedManagementCluster(state);
    const methods = useForm<ClusterTopologyStepFormInputs>({ resolver: yupResolver(clusterTopologyStepFormSchema) });
    const {
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

    const onFieldChange = (value: string, fieldName?: string | undefined) => {
        console.log('This part need to be refactored');
    };

    return (
        <div className="wizard-content-container" key="cluster-topology">
            <p cds-text="heading">Workload Topology Settings</p>
            <br />
            {ManagementClusterInfoBanner(cluster)}
            <br />
            <div cds-layout="grid gap:xxl" key="section-holder">
                <div cds-layout="col:6" key="cluster-name-section">
                    <ClusterName field={VSPHERE_FIELDS.CLUSTERNAME} clusterNameChange={onFieldChange} />
                </div>
                <div cds-layout="col:6" key="instance-type-section">
                    <NodeProfile
                        field={VSPHERE_FIELDS.INSTANCETYPE}
                        nodeInstanceTypes={workerNodeInstanceTypes}
                        nodeInstanceTypeChange={onFieldChange}
                    />
                </div>
            </div>
            <br />
            <CdsButton onClick={handleSubmit(onSubmit)}>NEXT</CdsButton>
        </div>
    );
}

export default ClusterTopologyStep;
