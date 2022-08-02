// React imports
import React, { ChangeEvent } from 'react';

// Library imports
import { CdsSelect } from '@cds/react/select';
import { CdsControlMessage } from '@cds/react/forms';
import { CdsIcon } from '@cds/react/icon';
import { UseFormRegister } from 'react-hook-form';
import { CdsProgressCircle } from '@cds/react/progress-circle';

// App imports
import './SpinnerSelect.scss';
import Tooltip from '../Tooltip/Tooltip';
import TooltipContent from '../Tooltip/TooltipContent/TooltipContent';

interface Props {
    label?: string;
    className?: string;
    defaultValue?: string;
    error?: string;
    disabled?: boolean;
    controlMessage?: string;
    isLoading?: boolean;
    name: string;
    tooltipMessage?: string;
    handleSelect: (event: ChangeEvent<HTMLSelectElement>) => void;
    children: React.ReactNode;
    register: UseFormRegister<any>;
}
function SpinnerSelect(props: Props) {
    const {
        children,
        controlMessage,
        className,
        defaultValue,
        disabled,
        error,
        handleSelect,
        label,
        name,
        register,
        isLoading,
        tooltipMessage,
    } = props;
    return (
        <div className="select-container">
            <CdsSelect layout="compact">
                <label>
                    {label}
                    {tooltipMessage && (
                        <Tooltip>
                            <CdsIcon shape="info-circle" size="md"></CdsIcon>
                            <TooltipContent position={'top-right'}>
                                <div>{tooltipMessage}</div>
                            </TooltipContent>
                        </Tooltip>
                    )}
                </label>
                <select
                    className={className}
                    {...register(name, {
                        onChange: handleSelect,
                    })}
                    defaultValue={defaultValue}
                    data-testid={name.toLowerCase() + '-select'}
                    disabled={isLoading || disabled}
                >
                    {children}
                </select>
                {error && (
                    <CdsControlMessage status="error" className="error-height">
                        {error}
                    </CdsControlMessage>
                )}
                <CdsControlMessage className="control-message-width">{controlMessage}</CdsControlMessage>
            </CdsSelect>
            {isLoading && (
                <div className="select-spinner-container">
                    <CdsProgressCircle size="sm" aria-label={`loading ${label} options`} status="info"></CdsProgressCircle>
                </div>
            )}
        </div>
    );
}

export default SpinnerSelect;
