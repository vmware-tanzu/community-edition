// React imports
import React from 'react';
import { useNavigate } from 'react-router-dom';

// Library imports
import styled from 'styled-components';
import { CdsButton } from '@cds/react/button';
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';

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

const GettingStarted: React.FC = () => {
    const navigate = useNavigate();
    return (
        <div cds-layout="vertical gap:lg gap@md:xl col@sm:12">
            <Header>
                <Title>Getting Started</Title>
            </Header>
            <Description>
                TODO: This page should reflect the Getting Started mockups when ready
                <br/><br/>
                TODO: refactor header and description elements to reflect mockup layout
            </Description>
            <SubTitle>What do you want to do?</SubTitle>
            <ButtonContainer>
                <CdsButton status="primary" onClick={()=> navigate(NavRoutes.MANAGEMENT_CLUSTER_LANDING)}>
                    MC
                </CdsButton>
                &nbsp;
                <CdsButton status="neutral" onClick={()=> navigate(NavRoutes.WORKLOAD_CLUSTER_LANDING)}>
                    WC
                </CdsButton>
            </ButtonContainer>
        </div>
    );
};

export default GettingStarted;
