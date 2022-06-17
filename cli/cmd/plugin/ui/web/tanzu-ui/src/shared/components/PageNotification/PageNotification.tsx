// React imports
import React from 'react';

// Library imports
import { CdsAlert, CdsAlertGroup } from '@cds/react/alert';
import { StatusTypes } from '@cds/core/internal';

export const NotificationStatus = {
    DANGER: 'danger',
    INFO: 'info',
    SUCCESS: 'success',
    WARNING: 'warning',
};

export interface Notification {
    status: string;
    message: string;
}

interface PageNotificationProps {
    notification: Notification | null;
    closeChangeCallback?: () => void;
}

function PageNotification(props: PageNotificationProps) {
    const { notification, closeChangeCallback } = props;

    function renderAlert() {
        if (notification?.status && notification?.message) {
            return (
                <CdsAlertGroup cds-layout="col:12" status={notification.status as StatusTypes}>
                    <CdsAlert closable aria-label="page notification" onCloseChange={closeChangeCallback}>
                        {notification.message}
                    </CdsAlert>
                </CdsAlertGroup>
            );
        }
    }

    return <>{renderAlert()}</>;
}

export default PageNotification;
