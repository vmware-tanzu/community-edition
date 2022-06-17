// React imports
import React from 'react';
import { StatusTypes } from '@cds/core/internal';

// Library imports
import { CdsAlert, CdsAlertGroup } from '@cds/react/alert';

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
        if (!notification?.status && !notification?.message) {
            return;
        } else {
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
