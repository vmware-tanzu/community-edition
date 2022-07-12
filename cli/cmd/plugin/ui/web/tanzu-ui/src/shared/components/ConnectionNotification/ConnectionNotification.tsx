// React imports
import React from 'react';
// Library imports
import { CdsAlert, CdsAlertGroup } from '@cds/react/alert';
// App imports
import './ConnectionNotification.scss';

export function ConnectionNotification(
    connected: boolean,
    connectedMessage: string,
    connecting: boolean,
    connectingMessage: string,
    errorMessage?: string
) {
    return (
        <div>
            {connected && (
                <CdsAlertGroup status="success" aria-label={connectedMessage}>
                    <CdsAlert>{connectedMessage}</CdsAlert>
                </CdsAlertGroup>
            )}
            {connecting && (
                <CdsAlertGroup status="loading" aria-label={connectedMessage}>
                    <CdsAlert>{connectingMessage}</CdsAlert>
                </CdsAlertGroup>
            )}
            {!connected && errorMessage && (
                <CdsAlertGroup status="danger">
                    <CdsAlert>{errorMessage}</CdsAlert>
                </CdsAlertGroup>
            )}
            {!connected && !connecting && !errorMessage && (
                <div className="hide-me">
                    <CdsAlertGroup status="neutral">
                        <CdsAlert>&nbsp;</CdsAlert>
                    </CdsAlertGroup>
                </div>
            )}
        </div>
    );
}
