// React imports
import React from 'react';
import { Link } from 'react-router-dom';

// Library imports
import { CdsButton } from '@cds/react/button';

// App imports
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';
import './Welcome.scss';
import TceLogo from '../../assets/tce-logo.svg';
import RolloverBanner  from './RolloverBanner/RolloverBanner';

const Welcome: React.FC = () => {
    return (
        <>
            <div cds-layout="vertical gap:md gap@md:lg col:8">
                <p cds-text="title">
                    <img src={TceLogo} className="logo logo-42" alt="tce logo"/>
                    Welcome to Tanzu Community Edition
                </p>
                <p cds-text="subsection">
                    Tanzu Community Edition is VMware&apos;s Open Source Kubernetes distribution. VMware Tanzu Community Edition
                    is a full-featured, easy-to-manage Kubernetes platform for learners and users, especially those working
                    in small-scale or pre-production environments.
                </p>
                <p cds-text="section" className="text-blue">Ready to dive in?</p>
                <p cds-text="subsection">
                    Get started with creating a local development environment or a production-ready environment on a cloud provider.
                </p>
                <nav>
                    <Link to={NavRoutes.GETTING_STARTED}>
                        <CdsButton>Let&apos;s Get Started</CdsButton>
                    </Link>
                </nav>
            </div>
            <div cds-layout="col:4 container:fill">
                <CdsButton action="flat" onClick={() => {
                    window.open('http://tanzucommunityedition.io', '_blank');
                }}>Learn more at tanzucommunityedition.io</CdsButton>
            </div>
            <div cds-layout="grid col:12 gap:md align:stretch">
                <p cds-text="title" cds-layout="col:12 p-y:md">
                    Tanzu Community Edition Overview
                </p>
                <RolloverBanner/>
            </div>
        </>
    );
};

export default Welcome;
