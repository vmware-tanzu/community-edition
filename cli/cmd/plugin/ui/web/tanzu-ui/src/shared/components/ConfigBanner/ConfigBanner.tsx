// React imports
import React from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import './ConfigBanner.scss';

function ConfigBanner() {
    return (
        <div className="banner-container">
            <p>Basic configuration</p>
            <CdsButton action="outline" className="btn">
                Use advanced configuration
            </CdsButton>
        </div>
    );
}

export default ConfigBanner;
