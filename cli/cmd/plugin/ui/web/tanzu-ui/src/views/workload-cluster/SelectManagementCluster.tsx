// React imports
import React, { useContext, useState } from 'react';

// Library imports
import styled from 'styled-components';
import { CdsButton } from '@cds/react/button';
import { CdsControlMessage } from '@cds/react/forms';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

// App imports
import { ManagementCluster } from '../../shared/models/ManagementCluster';
import './select-management-cluster.scss';
import { StepProps } from '../../shared/components/wizard/Wizard';
import { SubmitHandler, useForm } from 'react-hook-form';
import { WcStore } from '../../state-management/stores/Store.wc';

const selectManagementClusterFormSchema = yup.object({
    SELECTED_MANAGEMENT_CLUSTER: yup.string().nullable().required('Please select a management cluster')
}).required();

interface SelectManagementClusterFormInputs {
    SELECTED_MANAGEMENT_CLUSTER: string;
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
    const { state, dispatch } = useContext(WcStore);
    const methods = useForm<SelectManagementClusterFormInputs>({
        resolver: yupResolver(selectManagementClusterFormSchema),
    });
    const { register, handleSubmit, formState: { errors } } = methods;
    const { currentStep, goToStep, submitForm, retrieveManagementClusters, handleValueChange } = props;
    const onSubmit: SubmitHandler<SelectManagementClusterFormInputs> = (data) => {
        if (Object.keys(errors).length === 0) {
            if (goToStep && currentStep && submitForm) {
                console.log(`submitting form is handleValueChange a SUBMIT_FORM with SELECTED_MANAGEMENT_CLUSTER=${data.SELECTED_MANAGEMENT_CLUSTER}`);
                handleValueChange && handleValueChange('SELECTED_MANAGEMENT_CLUSTER', data.SELECTED_MANAGEMENT_CLUSTER, currentStep, errors);
                goToStep(currentStep + 1);
                submitForm(currentStep);
            }
        }
    };

    // TODO: validation on props' retrieveManagementClusters
    const clusters = retrieveManagementClusters ? retrieveManagementClusters() : [];

    return (
        <div >
            <Description>
                Select a Management Cluster from which the workload cluster will be provisioned.
                <br/>
                After creation, the workload cluster can be used to run your application workloads.
                <br/>
                state={JSON.stringify(state)}
            </Description>
            <div>
            <SubTitle>Select a Management Cluster</SubTitle>
                <div cds-layout="grid gap:md" key="header-mc-grid">
                    <div className="header-mc-grid" cds-layout="col:1"></div>
                    <div className="text-white header-mc-grid" cds-layout="col:5">Cluster Name</div>
                    <div className="text-white header-mc-grid" cds-layout="col:1">Provider</div>
                    <div className="text-white header-mc-grid" cds-layout="col:1">Created</div>
                    <div className="text-white header-mc-grid" cds-layout="col:4">Description</div>

                    {
                        clusters.map((cluster: ManagementCluster) => { return ManagementClusterInList(cluster, register); })
                    }
                </div>
            </div>

            <br/>
            { errors.SELECTED_MANAGEMENT_CLUSTER && <CdsControlMessage status="error">{errors.SELECTED_MANAGEMENT_CLUSTER.message}</CdsControlMessage> }
            <br/>
            <CdsButton onClick={handleSubmit(onSubmit)}>NEXT</CdsButton>
        </div>
    );
}

function ManagementClusterInList(cluster: ManagementCluster, register: any) {
    return         <><input className="inputradio" cds-layout="col:1" type="radio" value={cluster.name} name="{`management-cluster-radio`}"
                 {...register("SELECTED_MANAGEMENT_CLUSTER")} />
            <div className="text-white" cds-layout="col:5">{cluster.name}</div>
            <div className="text-white" cds-layout="col:1">{cluster.provider}</div>
            <div className="text-white" cds-layout="col:1">{cluster.created}</div>
            <div className="text-white" cds-layout="col:4">{cluster.description}</div>
        </>
        ;
}

export default SelectManagementCluster;
