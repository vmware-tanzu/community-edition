// React imports
import React from 'react';

// Library imports
import { CdsProgressCircle } from '@cds/react/progress-circle';
import { StatusTypes } from '@cds/core/internal';

export const LoadingSpinnerStatus = {
    DANGER: 'danger',
    INFO: 'info',
    SUCCESS: 'success',
    WARNING: 'warning',
};

interface PageLoadingProps {
    message?: string;
    status?: string;
}

function PageLoading(props: PageLoadingProps) {
    const { message, status } = props;

    const statusType = status ? status : LoadingSpinnerStatus.INFO;

    return (
        <>
            <div cds-layout="horizontal align:center">
                <CdsProgressCircle status={statusType as StatusTypes} size="xl"></CdsProgressCircle>
                {message && (
                    <span cds-text="section" cds-layout="p-l:sm">
                        {message}
                    </span>
                )}
            </div>
        </>
    );
}

export default PageLoading;
