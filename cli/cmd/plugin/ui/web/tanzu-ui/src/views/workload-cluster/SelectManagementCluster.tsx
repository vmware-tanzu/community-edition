// React imports
import { CdsButton } from '@cds/react/button';
import { CdsCard } from '@cds/react/card';
import React, { useEffect, useState } from 'react';

// Library imports
import { CdsProgressCircle } from '@cds/react/progress-circle';
import styled from 'styled-components';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

// App imports
import './SelectManagementCluster.scss';
import { CancelablePromise, ManagementCluster } from '../../swagger-api';
import { INPUT_CHANGE } from '../../state-management/actions/Form.actions';
import { StepProps } from '../../shared/components/wizard/Wizard';
import { useForm } from 'react-hook-form';
import { ProviderData, retrieveProviderInfo } from '../../shared/services/Provider.service';

const selectManagementClusterFormSchema = yup
    .object({
        SELECTED_MANAGEMENT_CLUSTER_NAME: yup.string().nullable().required('Please select a management cluster'),
    })
    .required();

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
    retrieveManagementClusters: () => CancelablePromise<Array<ManagementCluster>>;
    selectedManagementCluster: string;
}

function SelectManagementCluster(props: Partial<SelectManagementClusterProps>) {
    const [clusters, setClusters] = useState<ManagementCluster[]>([]);
    const [loadingClusters, setLoadingClusters] = useState<boolean>(false);
    const methods = useForm<SelectManagementClusterFormInputs>({
        resolver: yupResolver(selectManagementClusterFormSchema),
    });
    const {
        formState: { errors },
    } = methods;
    const { currentStep, goToStep, submitForm, retrieveManagementClusters, handleValueChange } = props;

    const onSelectManagementCluster = (cluster: ManagementCluster) => {
        handleValueChange && handleValueChange(INPUT_CHANGE, 'SELECTED_MANAGEMENT_CLUSTER', cluster, currentStep, errors);
        if (currentStep) {
            submitForm && submitForm(currentStep);
            goToStep && goToStep(currentStep + 1);
        }
    };

    // fxn to either return the management clusters, or log a problem
    const retrieveMC = retrieveManagementClusters
        ? () => {
              retrieveManagementClusters().then((data) => {
                  setClusters(data);
                  setLoadingClusters(false);
              });
          }
        : () => {
              console.log('Wanted to call retrieveManagementClusters(), but no method passed to SelectManagementCluster!');
          };

    // Load the management cluster data when component first initialized
    useEffect(() => {
        setLoadingClusters(true);
        // FOR TESTING, SIMULATE 2-second call
        setTimeout(retrieveMC, 1000);
    }, []); // eslint-disable-line react-hooks/exhaustive-deps

    return (
        <div className="wizard-content-container" cds-layout="container:fill">
            <Description>
                Select a Management Cluster from which the workload cluster will be provisioned.
                <br />
                After creation, the workload cluster can be used to run your application workloads.
                <br />
            </Description>
            <div key="subtitle">
                <SubTitle>Select a Management Cluster</SubTitle>
                {loadingClusters && ManagementClusterLoading()}
                {!loadingClusters && ManagementClusterLayout(clusters, onSelectManagementCluster)}
            </div>
        </div>
    );
}

function ManagementClusterLoading() {
    return (
        <div cds-layout="grid gap:md">
            <div cds-layout="col:12">
                <div cds-layout="grid gap:md">{ManagementClusterListLoading()}</div>
            </div>
        </div>
    );
}

function ManagementClusterLayout(clusters: ManagementCluster[], onSelectManagementCluster: (mc: ManagementCluster) => void) {
    return (
        <>
            {clusters.map((cluster: ManagementCluster) => {
                return ManagementClusterInList(cluster, onSelectManagementCluster);
            })}
        </>
    );
}

function ManagementClusterInList(cluster: ManagementCluster, onSelectManagementCluster: (mc: ManagementCluster) => void) {
    const onButtonClick = () => {
        onSelectManagementCluster(cluster);
    };
    const provider: ProviderData = retrieveProviderInfo(cluster.provider || 'unknown');
    return (
        <>
            <CdsCard className="section-raised">
                <div cds-layout="grid gap:md">
                    <div className="text-white" cds-layout="col:1">
                        <img src={provider.logo} className="logo logo-42" cds-layout="m-r:md" alt={`${provider.name} logo`} />
                    </div>
                    <div cds-layout="col:10">
                        <div cds-layout="grid gap:md">
                            <div className="text-white" cds-layout="col:12">
                                {cluster.name} <span className="text-context-label">| MGMT CLUSTER</span>
                            </div>
                            <div className="text-context-info" cds-layout="col:12">
                                <div cds-layout="grid gap:sm">
                                    <div cds-layout="col:1">
                                        <span className="text-context-label">Created:</span>
                                    </div>
                                    <div cds-layout="col:11">
                                        <span className="text-context-info">{cluster.created}</span>
                                    </div>
                                    <div cds-layout="col:1">
                                        <span className="text-context-label">Description:</span>
                                    </div>
                                    <div cds-layout="col:11">
                                        <span className="text-context-info">{cluster.description}</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div className="text-white" cds-layout="col:1">
                        <CdsButton onClick={onButtonClick}>SELECT</CdsButton>
                    </div>
                </div>
            </CdsCard>
        </>
    );
}

function ManagementClusterListLoading() {
    return (
        <>
            <div className="text-white" cds-layout="col:1"></div>
            <div cds-layout="horizontal gap:sm col:11" cds-theme="dark">
                <CdsProgressCircle size="xl" status="info"></CdsProgressCircle>
            </div>
        </>
    );
}

export default SelectManagementCluster;
