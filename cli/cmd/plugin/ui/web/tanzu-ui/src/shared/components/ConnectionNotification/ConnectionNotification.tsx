// React imports
import React from 'react';

// Library imports
import { AlertStatusTypes } from '@cds/core/alert';
import { CdsAlert, CdsAlertGroup } from '@cds/react/alert';

// App imports
import './ConnectionNotification.scss';
import '../../../../src/scss/utils.scss';

export enum CONNECTION_STATUS {
    DISCONNECTED,
    CONNECTED,
    CONNECTING,
    ERROR,
}

const connectionStatusToAlertStatus = {
    [CONNECTION_STATUS.CONNECTED]: 'success',
    [CONNECTION_STATUS.CONNECTING]: 'loading',
    [CONNECTION_STATUS.DISCONNECTED]: 'neutral',
    [CONNECTION_STATUS.ERROR]: 'danger',
};

// NOTE: there SHOULD NOT be a need to have CONNECTING as a separate case below, but the Clarity CdsAlertGroup does not repaint
// correctly without it.
interface Props {
    status: CONNECTION_STATUS;
    message: string;
}
function ConnectionNotification(props: Props) {
    const { status, message } = props;
    const alertStatus = (connectionStatusToAlertStatus[status] || 'neutral') as AlertStatusTypes;
    return (
        <div className="connection-notification-container">
            {status !== CONNECTION_STATUS.DISCONNECTED && status !== CONNECTION_STATUS.CONNECTING && (
                <CdsAlertGroup status={alertStatus} aria-label={message}>
                    <CdsAlert>{message}</CdsAlert>
                </CdsAlertGroup>
            )}
            {status === CONNECTION_STATUS.CONNECTING && (
                <CdsAlertGroup status={alertStatus} aria-label={message}>
                    <CdsAlert>{message}</CdsAlert>
                </CdsAlertGroup>
            )}
            {status === CONNECTION_STATUS.DISCONNECTED && (
                <div className="hide-me">
                    <CdsAlertGroup status="neutral">
                        <CdsAlert>&nbsp;</CdsAlert>
                    </CdsAlertGroup>
                </div>
            )}
        </div>
    );
}

export default ConnectionNotification;
