// React imports
import React from 'react';
import { Link } from 'react-router-dom';

// Library imports
import { CdsButton } from '@cds/react/button';

// App imports
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';
import RolloverBanner from './RolloverBanner/RolloverBanner';
import './Welcome.scss';

const Welcome: React.FC = () => {
    return (
        <>
            <div cds-layout="vertical col:8 ">
                <img className="title-image tce-log-and-title-img" alt="tce logo" />
            </div>
            <div cds-layout="vertical gap:lg col:8">
                <div cds-text="subsection">
                    Tanzu Community Edition is VMware&apos;s Open Source Kubernetes distribution. VMware Tanzu Community Edition is a
                    full-featured, easy-to-manage Kubernetes platform for learners and users, especially those working in small-scale or
                    pre-production environments.
                </div>
                <div className="section-raised getting-started-container" cds-layout="grid vertical gap:lg p:lg">
                    <div cds-text="title" cds-layout="col:12">
                        Ready to dive in?
                    </div>
                    <div cds-text="body" cds-layout="col:7">
                        Get started with creating a local development environment or a production-ready environment on a cloud provider.
                    </div>
                    <nav cds-layout="col:12">
                        <Link to={NavRoutes.GETTING_STARTED}>
                            <CdsButton>Let&apos;s Get Started</CdsButton>
                        </Link>
                    </nav>
                </div>
            </div>
            <div cds-layout="grid col:12 gap:md align:stretch">
                <div cds-text="title" cds-layout="col:12 p-y:md">
                    Explore VMware Tanzu
                </div>
                <RolloverBanner />
            </div>
        </>
    );
};

export default Welcome;
