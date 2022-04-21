// React imports
import React from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import styled from 'styled-components';
import { CdsButton } from '@cds/react/button';
import { CdsCard } from '@cds/react/card';

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

const ButtonContainer = styled.div`
    text-align: center;
    padding-top: 50px;
`;

const ManagementClusterLanding: React.FC = () => {
    const navigate = useNavigate();
    const cards = [{
        name: 'Docker',
        path: '/docker'
    }, {
        name: 'VMware vsphere',
        path: '/vsphere'
    }, {
        name: 'Microsoft Azure',
        path: '/azure'
    }, {
        name: 'Amazon EC2',
        path: '/aws'
    }];
    return (
        <div cds-layout="vertical gap:lg gap@md:xl col@sm:12">
            <Header>
                <Title>Management Cluster</Title>
            </Header>
            <Description>
                The management cluster provides management and operations for Tanzu. It runs Cluster-API, which is used to manage
                workload clusters and multi-cluster services. (The workload clusters are where developers&apos; workloads run.)
                <br/><br/>
                TODO: refactor header and description elements to reflect mockup layout
            </Description>
            <SubTitle>Select a supported cloud provider</SubTitle>
            <div cds-layout="grid cols@md:6 cols@lg:3 gap:sm">
                {
                    cards.map((card, index) => {
                        return (
                            <CdsCard aria-labelledby="containerOfCards1" key={index}>
                                <div cds-layout="vertical gap:md">
                                    <h2 id="containerOfCards1" cds-text="section" cds-layout="horizontal align:vertical-center">
                                        {card.name}
                                    </h2>

                                    <div cds-text="body light">
                                        <ButtonContainer>
                                            <CdsButton status="primary" onClick={()=> navigate(card.path)}>
                                                Deploy
                                            </CdsButton>
                                        </ButtonContainer>
                                    </div>
                                </div>
                            </CdsCard>
                        );
                    })
                }
            </div>
        </div>
    );
};

export default ManagementClusterLanding;
