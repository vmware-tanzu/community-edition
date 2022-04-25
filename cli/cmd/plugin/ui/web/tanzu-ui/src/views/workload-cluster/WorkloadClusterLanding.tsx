// React imports
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import styled from 'styled-components';
import { CdsCard } from '@cds/react/card';

import ManagementClusterLogo from '../../assets/management-cluster.svg';
import { CdsButton } from '@cds/react/button';
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';

const Section = styled.section`
    padding: 20px;
`;

const Header = styled.div`
    margin-top: 25px;
    display: flex;
    align-items: center;
`;

const Title = styled.span`
    padding-left: 20px;
    font-size: 28px;
`;

const Description = styled.p`
    padding: 20px;
    line-height: 24px;
`;

const SubTitle = styled.h3`
    padding-left: 20px;
`;

const fakeServiceRetrievesManagmentClusterObjects = () => {
    return [
        { name: 'shimon-test-cluster-1', provider: 'aws', dateCreated: '10/22/2021', description: 'just fooling around'},
        { name: 'some-other-cluster', provider: 'vsphere', dateCreated: '1/13/2022', description: 'a very serious cluster'}
    ];
};

const WorkloadClusterLanding: React.FC = () => {
    const navigate = useNavigate();
    const clusters = fakeServiceRetrievesManagmentClusterObjects();
    const [selectedManagementCluster, setSelectedManagementCluster ] = useState('(not yet selected)');

    const onSelectManagementCluster = (selectedManagementCluster: string) => {
        setSelectedManagementCluster(selectedManagementCluster);
        alert(`You selected ${selectedManagementCluster}`);
        navigate(NavRoutes.WELCOME);
    };

    return (
        <Section>
            <Header>
                <Title>Workload Cluster</Title>
            </Header>
            <Description>
                The workload cluster is deployed by the management. The workload clusters is used to run your application workloads.
                <br/><br/>
                TODO: refactor header and description elements to reflect mockup layout
            </Description>
            <SubTitle>Select a management cluster</SubTitle>
            <div cds-layout="vertical">
                {
                    clusters.map((cluster) => { return ManagementClusterInList(cluster, onSelectManagementCluster); })
                }
            </div>
            <div>You have selected cluster: {selectedManagementCluster}</div>
        </Section>
    );
};

function ManagementClusterInList(cluster: any, setter: any) {
    return <CdsCard aria-labelledby="containerOfCards1" key={`management-cluster-${cluster.name}`}>
        <h2 id="containerOfCards1" cds-text="section" cds-layout="horizontal align:vertical-center">
            <img src={ManagementClusterLogo} className="logo logo-26"/> &nbsp;
            {cluster.name}
            <div cds-layout="align:right">
                <CdsButton status="primary" onClick={()=> setter(cluster.name)}>
                    SELECT
                </CdsButton>
            </div>
        </h2>
    </CdsCard>;
}
export default WorkloadClusterLanding;
