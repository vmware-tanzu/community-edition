// React imports
import React from 'react';
import { useNavigate } from 'react-router-dom';

// App imports
import { NavRoutes } from '../../constants/NavRoutes.constants';
import VMWLogo from '../../../assets/vmw-logo.svg';
import './HeaderBar.scss';

function HeaderBar() {
    const navigate = useNavigate();

    const navigateToWelcome = (): void => {
        navigate(NavRoutes.WELCOME);
    };

    return (
        <div className="header">
            <div
                className="branding"
                aria-label="navigate-to-welcome"
                onClick={() => {
                    navigateToWelcome();
                }}
            >
                <img src={VMWLogo} className="logo logo-26" alt="vmware logo home" aria-label="header-logo" />
                <span className="title" aria-label="header-title">
                    Tanzu
                </span>
            </div>
        </div>
    );
}

export default HeaderBar;
