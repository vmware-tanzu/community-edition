// React imports
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage } from '@cds/react/forms';
import React, { ChangeEvent } from 'react';

// Library imports
import styled from 'styled-components';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

// App imports
import './select-management-cluster.scss';
import { ManagementCluster } from '../../shared/models/ManagementCluster';
import RadioButton from '../../shared/components/widgets/RadioButton';
import { StepProps } from '../../shared/components/wizard/Wizard';
import { SubmitHandler, useForm } from 'react-hook-form';

const selectManagementClusterFormSchema = yup.object({
    SELECTED_MANAGEMENT_CLUSTER_NAME: yup.string().nullable().required('Please select a management cluster')
}).required();

interface SelectManagementClusterFormInputs {
    SELECTED_MANAGEMENT_CLUSTER_NAME: string;
}

const Description = styled.p`
    padding: 20px;
    line-height: 24px;
`;

const SubTitle = styled.h3`
    padding-left: 20px;
`;

interface SelectManagementClusterProps extends StepProps {
    retrieveManagementClusters: () => ManagementCluster[],
    selectedManagementCluster: string,
}

function SelectManagementCluster (props: Partial<SelectManagementClusterProps>) {
    // const { state } = useContext(WcStore);
    const methods = useForm<SelectManagementClusterFormInputs>({
        resolver: yupResolver(selectManagementClusterFormSchema),
    });
    const { register, handleSubmit, formState: { errors } } = methods;
    const { currentStep, goToStep, submitForm, retrieveManagementClusters, handleValueChange } = props;
    const onSubmit: SubmitHandler<SelectManagementClusterFormInputs> = (data) => {
        if (Object.keys(errors).length === 0) {
            if (goToStep && currentStep && submitForm) {
                if (handleValueChange) {
                    const clusterName = data.SELECTED_MANAGEMENT_CLUSTER_NAME;
                    const cluster = findClusterFromName(clusterName, clusters);
                    handleValueChange('SELECTED_MANAGEMENT_CLUSTER', cluster, currentStep, errors);
                } else {
                    console.error('no handleChangeValue passed to SelectManagementCluster()?!');
                }
                goToStep(currentStep + 1);
                submitForm(currentStep);
            }
        }
    };

    // TODO: validation on props' retrieveManagementClusters
    const clusters = retrieveManagementClusters ? retrieveManagementClusters() : [];
    const findClusterFromName = (clusterName: string, clusters: ManagementCluster[]) => {
        return clusters.find(cluster => cluster.name === clusterName);
    };

    // If SELECTED_MANAGEMENT_CLUSTER_NAME is already set, select that management cluster
    // TODO: use: const selectedClusterName = state.data?.SELECTED_MANAGEMENT_CLUSTER?.name
    const onSelectManagementCluster = (evt: ChangeEvent<HTMLSelectElement>) => {
        const clusterName = evt.target.value;
        const cluster = findClusterFromName(clusterName, clusters);
        if (handleValueChange) {
            handleValueChange('SELECTED_MANAGEMENT_CLUSTER', cluster, currentStep, errors);
        } else {
            console.error('Unable to record selected management cluster because handleValueChange method is null/undefined')
        }
    }

    return (
        <div className="wizard-content-container" cds-layout="container:fill">
            <Description>
                Select a Management Cluster from which the workload cluster will be provisioned.
                <br/>
                After creation, the workload cluster can be used to run your application workloads.
                <br/>
            </Description>
            <div key="subtitle">
                <SubTitle>Select a Management Cluster</SubTitle>
                {
                    ManagementClusterLayout(clusters, onSelectManagementCluster, register)
                }
            </div>

            <br/>
            { errors.SELECTED_MANAGEMENT_CLUSTER_NAME &&
                <CdsControlMessage status="error">{errors.SELECTED_MANAGEMENT_CLUSTER_NAME.message}</CdsControlMessage>
            }
            <br/>
            <CdsButton onClick={handleSubmit(onSubmit)}>NEXT</CdsButton>
        </div>
    );
}

function ManagementClusterLayout(clusters: ManagementCluster[], onSelectManagementCluster: any, register: any) {
    return <div cds-layout="grid gap:md">
        <div cds-layout="col:12">
            <div cds-layout="grid gap:md">
                { ManagementClusterListHeader() }
                { clusters.map((cluster: ManagementCluster) => {
                    return ManagementClusterInList(cluster, register, false)
                })
                }
            </div>
        </div>
    </div>

}

function ManagementClusterListHeader() {
    return <>
        <div className="text-white header-mc-grid" cds-layout="col:1" key="mc-select-grid-col0"></div>
        <div className="text-white header-mc-grid" cds-layout="col:5" key="mc-select-grid-col1">Cluster Name</div>
        <div className="text-white header-mc-grid" cds-layout="col:1" key="mc-select-grid-col2">Provider</div>
        <div className="text-white header-mc-grid" cds-layout="col:1" key="mc-select-grid-col3">Created</div>
        <div className="text-white header-mc-grid" cds-layout="col:4" key="mc-select-grid-col4">Description</div>
    </>;
}

function ManagementClusterInList(cluster: ManagementCluster, register: any, selected: boolean) {
    return  <>
        <RadioButton name="SELECTED_MANAGEMENT_CLUSTER_NAME" className="input-radio" cdsLayout="col:1"
            checked={selected} register={register} value={cluster.name} />
        <div className="text-white" cds-layout="col:5" key={`${cluster.name}-col1`}>{cluster.name}</div>
        <div className="text-white" cds-layout="col:1" key={`${cluster.name}-col2`}>{cluster.provider}</div>
        <div className="text-white" cds-layout="col:1" key={`${cluster.name}-col3`}>{cluster.created}</div>
        <div className="text-white" cds-layout="col:4" key={`${cluster.name}-col4`}>{cluster.description}</div>
    </>
}

export default SelectManagementCluster;
