import React from 'react';
import './TooltipContent.scss';

export type Position = 'left' | 'right' | 'top-right' | 'top-left' | 'bottom-left' | 'bottom-right';
export type Size = 'small' | 'medium' | 'large';

interface TooltipProps {
    position: Position;
    size?: Size;
    children?: React.ReactNode;
}

function TooltipContent(props: TooltipProps) {
    return (
        <div className={'tooltip ' + props.position + ' ' + props.size}>
            <div cds-text="body left" className="tooltip-content">
                {props.children}
            </div>
        </div>
    );
}

TooltipContent.defaultProps = {
    size: 'large',
};

export default TooltipContent;
