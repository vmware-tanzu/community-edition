//React imports
import React, { useContext } from 'react';
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { SubmitHandler, useForm } from 'react-hook-form';

// Library imports
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { CdsAlert, CdsAlertGroup } from '@cds/react/alert';
import { CdsInput } from '@cds/react/input';
import { CdsIcon } from '@cds/react/icon';
import { ClarityIcons, computerIcon, connectIcon, cpuIcon, flaskIcon, memoryIcon } from '@cds/core/icon';

// App imports
import './WorkloadClusterWizard.scss';
import { isValidClusterName } from '../../shared/validations/Validation.service';
import { ManagementCluster } from '../../shared/models/ManagementCluster';
import RadioButton from '../../shared/components/widgets/RadioButton';
import { StepProps } from '../../shared/components/wizard/Wizard';
import { WcStore } from '../../state-management/stores/Store.wc';

interface ClusterTopologyStepFormInputs {
    WORKLOAD_CLUSTER_NAME: string;
    SELECTED_WORKER_NODE_INSTANCE_TYPE: string;
}

const clusterTopologyStepFormSchema = yup.object({
    WORKLOAD_CLUSTER_NAME: yup.string().nullable().required('Please enter a name for your workload cluster').test('', 'Cluster name must contain only lower case letters and hyphen', value => value !== null && isValidClusterName(value)),
    SELECTED_WORKER_NODE_INSTANCE_TYPE: yup.string().nullable().required('Please select an instance type for your workload cluster nodes')
}).required();

interface WorkerNodeInstanceType {
    id: string,
    icon: string,
    name: string,
    description: string,
}

// NOTE: icons must be imported
const workerNodeInstanceTypes: WorkerNodeInstanceType[] = [
    {
        id: 'basic-demo', name: 'Basic demo', icon: 'flask',
        description: 'Virtual machines with a range of compute and memory resources. Intended for small projects and development' +
            ' environments.'
    },
    {
        id: 'general-purpose', name: 'General purpose', icon: 'computer',
        description: 'General purpose instances powered by multi-threaded CPUs. Balanced, high performance, compute and memory for' +
            ' production workloads.'
    },
    {
        id: 'compute-optimized', name: 'Compute optimized', icon: 'cpu',
        description: 'Compute optimized instances suited for CPU-intensive workloads such as CI/CD, machine learning, and data' +
            ' processing.'
    },
    {
        id: 'memory-optimized', name: 'Memory optimized', icon: 'memory',
        description: 'Memory optimized instances best suited for in-memory operations such as big-data and performant databases.'
    }
];

ClarityIcons.addIcons(flaskIcon, computerIcon, cpuIcon, memoryIcon);

function ClusterTopologyStep(props: Partial<StepProps>) {
    const {handleValueChange, currentStep, goToStep, submitForm, getValue} = props;
    const {state, dispatch} = useContext(WcStore);
    const methods = useForm<ClusterTopologyStepFormInputs>({
        resolver: yupResolver(clusterTopologyStepFormSchema),
    });
    const {
        register,
        handleSubmit,
        setValue,
        formState: {errors},
    } = methods;

    const onSubmit: SubmitHandler<ClusterTopologyStepFormInputs> = (data) => {
        if (Object.keys(errors).length === 0) {
            if (goToStep && currentStep && submitForm && handleValueChange) {
                handleValueChange('WORKLOAD_CLUSTER_NAME', data.WORKLOAD_CLUSTER_NAME, currentStep, errors);
                handleValueChange('SELECTED_WORKER_NODE_INSTANCE_TYPE', data.SELECTED_WORKER_NODE_INSTANCE_TYPE, currentStep, errors);
                goToStep(currentStep + 1);
                submitForm(currentStep);
            }
        }
    };

    let cluster = state.data.SELECTED_MANAGEMENT_CLUSTER
    if (!state.data.SELECTED_MANAGEMENT_CLUSTER) {
        console.log('ClusterTopologyState did not receive a selected cluster')
        cluster = getValue ? getValue('SELECTED_MANAGEMENT_CLUSTER') : undefined
    }
    return (<div className="wizard-content-container" key="cluster-topology">
        <p cds-text="heading">Workload Topology Settings</p>
        <br/>
        {ManagementClusterInfoBanner(cluster)}
        <br/>
        <div cds-layout="grid gap:md" key="section-holder">
            <div cds-layout="col:6" key="cluster-name-section"> {ClusterNameSection(errors, register)} </div>
            <div cds-layout="col:6" key="instance-type-section"> {WorkerNodeInstanceTypeSection(errors, register)} </div>
        </div>
        <br/>
        <CdsButton onClick={handleSubmit(onSubmit)}>NEXT</CdsButton>
    </div>);
}

function WorkerNodeInstanceTypeSection(errors: any, register: any) {
    return <div cds-layout="vertical gap:lg gap@md:lg col@sm:6 col:6">
        <div cds-layout="cols:6">Select a worker node instance type</div>
        <div cds-layout="grid gap:md align:fill">
            {
                workerNodeInstanceTypes.map(instanceType => {
                    return InstanceTypeInList(instanceType, register);
                })
            }
        </div>
        { errors.SELECTED_WORKER_NODE_INSTANCE_TYPE &&
            <CdsControlMessage status="error">{errors.SELECTED_WORKER_NODE_INSTANCE_TYPE.message}</CdsControlMessage>
        }
    </div>;
}

function InstanceTypeInList(instance: WorkerNodeInstanceType, register: any) {
    return <>
        <div className="text-white" cds-layout="col:1"><CdsIcon shape={instance.icon}></CdsIcon></div>
        <RadioButton className="input-radio" cdsLayout="col:1" value={instance.id} register={register}
                     name="SELECTED_WORKER_NODE_INSTANCE_TYPE" />
        <div className="text-white" cds-layout="col:10">{instance.name} {instance.description}</div>
    </>
        ;
}

function ClusterNameSection(errors: any, register: any) {
    return <div cds-layout="vertical gap:lg gap@md:lg col@sm:6 col:6" >
        <CdsFormGroup layout="vertical">
            <CdsInput layout="vertical">
                <label>Cluster Name</label>
                <input placeholder="workload-cluster-name" {...register("WORKLOAD_CLUSTER_NAME")} />
                { errors.WORKLOAD_CLUSTER_NAME && <CdsControlMessage status="error">{errors.WORKLOAD_CLUSTER_NAME.message}</CdsControlMessage> }
            </CdsInput>
        </CdsFormGroup>
        <div>Can only contain lowercase alphanumeric characters and dashes. </div>
        <div>You will use this workload cluster name when using the Tanzu CLI and kubectl utilities.</div>
    </div>;
}

function ManagementClusterInfoBanner(managementCluster: ManagementCluster) {
    if (!managementCluster) {
        return <></>
    }
    return <CdsAlertGroup
        type="banner"
        status="success"
        aria-label={`This workload cluster will be provisioned on ${managementCluster.provider} using ${managementCluster.name}`}
    >
        <CdsAlert closable>
            This workload cluster will be provisioned on {managementCluster.provider} using <b>{managementCluster.name}</b>
        </CdsAlert>
    </CdsAlertGroup>;
}
export default ClusterTopologyStep;
