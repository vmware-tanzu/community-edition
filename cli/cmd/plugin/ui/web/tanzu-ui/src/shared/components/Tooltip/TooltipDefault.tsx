import React from 'react';
import { CdsIcon } from '@cds/react/icon';
import Tooltip from './Tooltip';
import TooltipContent, { Position, Size } from './TooltipContent/TooltipContent';

export interface TooltipDefaultProps {
    tooltipMessage: string;
    positionClass?: Position;
    sizeClass?: Size;
    iconSize?: string;
}

const DEFAULT_POSITION = 'top-right';
const DEFAULT_SIZE = undefined;
const MAX_LINE_LENGTH = 70;

export function TooltipDefault(props: Partial<TooltipDefaultProps>) {
    const ttMessage = breakupLongString(props.tooltipMessage, ' ');
    const ttPositionClass = props.positionClass ?? DEFAULT_POSITION;
    const ttSizeClass = props.sizeClass ?? DEFAULT_SIZE;
    const iconSize = props.iconSize ?? 'md';
    return (
        <>
            {ttMessage && (
                <Tooltip>
                    <CdsIcon shape="info-circle" size={iconSize}></CdsIcon>
                    <TooltipContent position={ttPositionClass} size={ttSizeClass}>
                        <div>{ttMessage}</div>
                    </TooltipContent>
                </Tooltip>
            )}
        </>
    );
}

function breakupLongString(source: string | undefined, sep: string): string {
    if (!source) {
        return '';
    }
    const lines = source.match(new RegExp('.{1,' + MAX_LINE_LENGTH + '}', 'g'));
    if (!lines) {
        return source;
    }
    return lines.join(sep);
}
