import React from 'react';
import VMWLogo from '../../../assets/vmw-logo.svg';
import './HeaderBar.scss';

function HeaderBarComponent() {
    return (
        <div className="header" >
            {/* TODO: Click navigate home */}
            <div className="branding">
                <img src={VMWLogo} className="logo logo-26" alt="vmware logo home"/>
                <span className="title">
                    Tanzu
                </span>
            </div>
        </div>
    );
}

export default HeaderBarComponent;