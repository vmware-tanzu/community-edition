// React imports
import React from 'react';
import { useNavigate } from 'react-router-dom';

// App imports
import VMWLogo from '../../../assets/vmw-logo.svg';
import { NavRoutes } from '../../constants/NavRoutes.constants';
import ContextualHelp from '../ContextualHelp/ContextualHelp';
import './HeaderBar.scss';
import ContextualHelp from '../ContextualHelp/ContextualHelp';

function HeaderBar() {
    const navigate = useNavigate();

    const navigateToWelcome = (): void => {
        navigate(NavRoutes.WELCOME);
    };

    return (
        <div className="header" cds-layout="horizontal">
            <div
                className="branding"
                aria-label="navigate-to-welcome"
                onClick={() => {
                    navigateToWelcome();
                }}
            >
                <img src={VMWLogo} className="logo logo-26" alt="vmware logo home" aria-label="header-logo" />
                <span className="title" aria-label="header-title">
                    Tanzu Community Edition
                </span>
            </div>
            <ContextualHelp keywords={['tce-docker']} title={'Management cluster'} />
        </div>
    );
}

export default HeaderBar;
