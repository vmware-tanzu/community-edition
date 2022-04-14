// React imports
import React from 'react';
import { Link } from 'react-router-dom';

// Library imports
import { CdsButton } from '@cds/react/button';

// App imports
import { NavRoutes } from '../../shared/constants/NavRoutes.constants';
import './Welcome.scss';
import TceLogo from '../../assets/tce-logo.svg';


function Welcome() {
    return (
        <>
            <p cds-text="title">
                <img src={TceLogo} className="logo logo-42" alt="tce logo"/>
                Welcome to Tanzu Community Edition
            </p>
            <p>
                Tanzu Community Edition is VMware&apos;s Open Source Kubernetes distribution. VMware Tanzu Community Edition
                is a full-featured, easy-to-manage Kubernetes platform for learners and users, especially those working
                in small-scale or pre-production environments.
                <br/><br/>
                To get started you can create a <b>managed</b> or <b>unmanaged</b> cluster. An unmanaged cluster is good
                for quickly spinning up a lightweight cluster for a development/experimental environment. A managed cluster
                provides a full-featured, scalable Kubernetes implementation suitable for a development or production
                environment.
            </p>
            <nav>
                <Link to={NavRoutes.GETTING_STARTED}>
                    <CdsButton>Get Started</CdsButton>
                </Link>
            </nav>
        </>
    );
}

export default Welcome;