// React imports
import React from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import styled from 'styled-components';

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

const WorkloadClusterLanding: React.FC = () => {
    const navigate = useNavigate();
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
            <div>What? You&apos;re not seeing a list yet?!</div>
        </Section>
    );
};

export default WorkloadClusterLanding;
