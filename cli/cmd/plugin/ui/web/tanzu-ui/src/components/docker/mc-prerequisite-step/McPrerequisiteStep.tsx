// React imports
import React, { useEffect, useState } from 'react';

// Library imports
import { CdsAlert, CdsAlertGroup } from '@cds/react/alert';
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';

// App imports
import { CriService } from '../../../swagger-api';
import { StepProps } from '../../../shared/components/wizard/Wizard';
import './McPrerequisiteStep.scss';
import { STATUS } from '../../../shared/constants/App.constants';

function McPrerequisiteStep(props: Partial<StepProps>) {
    const { currentStep, goToStep, setTabStatus, tabStatus } = props;
    const [connected, setConnection] = useState(false);
    const [errorMessage, setErrorMessage] = useState('');
    const connect = async () => {
        try {
            await CriService.getContainerRuntimeInfo();
            setConnection(true);
        } catch (err: any) {
            setErrorMessage(err.body.message);
            setConnection(false);
        }
    };
    useEffect(() => {
        connect();
    }, []);

    const handleNext = () => {
        if (connected) {
            if (tabStatus && currentStep && setTabStatus) {
                console.log(tabStatus, currentStep);
                tabStatus[currentStep - 1] = STATUS.VALID;
                setTabStatus(tabStatus);
            }
            if (goToStep && currentStep) {
                goToStep(currentStep + 1);
            }
        }
    };
    return (
        <div className="wizard-content-container">
            <h2 cds-layout="m-t:lg">Docker prerequisite</h2>
            <p cds-layout="m-y:lg" className="description">
                Management cluster with the Docker daemon requires minimum
                allocated 4 CPUs and total memory of 6GB.
            </p>
            {connected && (
                <CdsAlertGroup
                    status="success"
                    aria-label="Management cluster with the Docker daemon requires minimum allocated 4 CPUs and total memory of 6GB."
                >
                    <CdsAlert cds-i18n='{ "closeButtonAriaLabel": "close dark theme success alert"}'>
                        Running Docker daemon
                    </CdsAlert>
                </CdsAlertGroup>
            )}  
            { !connected && errorMessage && <CdsAlertGroup status="danger">
                <CdsAlert cds-i18n='{ "closeButtonAriaLabel": "close dark theme success alert"}'>
                    {errorMessage}
                </CdsAlert>
            </CdsAlertGroup>
            }
            <div cds-layout="p-y:lg">
                <CdsButton onClick={connect} disabled={connected}>
                    <CdsIcon shape="connect" size="md"></CdsIcon>
                    CONNECT DOCKER DAEMON
                </CdsButton>
            </div>
            <CdsButton onClick={handleNext}>NEXT</CdsButton>
        </div>
    );
}

export default McPrerequisiteStep;
