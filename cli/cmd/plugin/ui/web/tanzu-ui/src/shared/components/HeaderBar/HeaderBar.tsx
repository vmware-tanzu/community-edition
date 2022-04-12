// React imports
import React from 'react';
import { Link, useNavigate } from 'react-router-dom';

// App imports
import { Routes } from '../../constants/Routes.constants';
import VMWLogo from '../../../assets/vmw-logo.svg';
import './HeaderBar.scss';

function HeaderBar() {
    const navigate = useNavigate();

    const navigateHome = (): void => {
        navigate(Routes.WELCOME);
    };

    return (
        <div className="header" >
            <div className="branding" aria-label="navigate-click" onClick={() => {
                navigateHome();
            }}>
                <img src={VMWLogo} className="logo logo-26" alt="vmware logo home" aria-label="header-logo" onClick={() => {
                    navigateHome();
                }}/>
                <span className="title" aria-label="header-title">
                    Tanzu
                </span>
            </div>
        </div>
    );
}

export default HeaderBar;