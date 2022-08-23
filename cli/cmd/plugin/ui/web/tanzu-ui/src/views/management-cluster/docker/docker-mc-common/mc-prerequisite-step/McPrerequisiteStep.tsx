// React imports
import React, { useCallback, useEffect, useState } from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';

// App imports
import ConnectionNotification, { CONNECTION_STATUS } from '../../../../../shared/components/ConnectionNotification/ConnectionNotification';
import { CriService } from '../../../../../swagger-api';
import { StepProps } from '../../../../../shared/components/wizard/Wizard';
import './McPrerequisiteStep.scss';

function McPrerequisiteStep(props: Partial<StepProps>) {
    const { currentStep, goToStep, submitForm } = props;
    const [connectionStatus, setConnectionStatus] = useState<CONNECTION_STATUS>(CONNECTION_STATUS.DISCONNECTED);
    const [message, setMessage] = useState('');
    const connect = useCallback(async () => {
        setConnectionStatus(CONNECTION_STATUS.CONNECTING);
        setMessage('Connecting to docker deamon');
        try {
            await CriService.getContainerRuntimeInfo();
            setConnectionStatus(CONNECTION_STATUS.CONNECTED);
            setMessage('Running Docker daemon');
        } catch (err: any) {
            setConnectionStatus(CONNECTION_STATUS.ERROR);
            setMessage(`Unable to connect to Docker: ${err.body.message}`);
        }
    }, []);

    useEffect(() => {
        connect();
    }, [connect]);

    const handleNext = () => {
        if (connectionStatus === CONNECTION_STATUS.CONNECTED) {
            if (goToStep && submitForm && currentStep) {
                goToStep(currentStep + 1);
                submitForm(currentStep);
            }
        }
    };
    return (
        <div className="wizard-content-container">
            <h2 cds-layout="m-t:md m-b:xl" cds-text="title">
                Docker prerequisite
            </h2>
            <p cds-layout="m-y:lg" className="description">
                Management cluster with the Docker daemon requires minimum allocated 4 CPUs and total memory of 6GB.
            </p>
            <ConnectionNotification status={connectionStatus} message={message}></ConnectionNotification>
            <div cds-layout="p-t:lg" className={connectionStatus === CONNECTION_STATUS.ERROR ? '' : 'hidden'}>
                <CdsButton onClick={connect}>
                    <CdsIcon shape="connect" size="md"></CdsIcon>
                    CONNECT DOCKER DAEMON
                </CdsButton>
            </div>
            <div cds-layout="p-y:lg">
                <CdsButton onClick={handleNext}>NEXT</CdsButton>
            </div>
        </div>
    );
}

export default McPrerequisiteStep;
