import React from 'react';
import { useNavigate } from 'react-router-dom';
import styled from 'styled-components';
 
import { CdsButton } from '@cds/react/button';
import { CdsCard } from '@cds/react/card';

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
const ButtonContainer = styled.div`
    text-align: center;
    padding-top: 50px;
`;
const LandingPage: React.FC = () => {
    const navigate = useNavigate();
    const cards = ['Docker', 'VMware vsphere', 'Microsoft Azure', 'Amazon EC2'];
    return (
        <Section>
            <Header> 
                <img src="/ui/tce-logo.png" alt="Logo"/>
                <Title>Welcome to the Tanzu Community Edition Launcher</Title>
            </Header>
            <Description>
                Tanzu Community Edition (TCE) is VMware&apos;s Open Source Kubernetes distribution.
                This installer will guide you through the installation of clusters necessary to get started with Kubernetes and TCE.
                For more details see the getting started guide.
            </Description>
            <SubTitle>What do you want to do?</SubTitle>
            <div cds-layout="grid cols@md:6 cols@lg:3 gap:sm">
                {
                    cards.map((card, index) => {
                        return (
                            <CdsCard aria-labelledby="containerOfCards1" key={index}>
                                <div cds-layout="vertical gap:md">
                                    <h2 id="containerOfCards1" cds-text="section" cds-layout="horizontal align:vertical-center">
                                        {card}
                                    </h2>

                                    <div cds-text="body light">
                                        <ButtonContainer>
                                            <CdsButton status="primary" onClick={()=> navigate('/vsphere')}>
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
        </Section>
        
    );
};

export default LandingPage;