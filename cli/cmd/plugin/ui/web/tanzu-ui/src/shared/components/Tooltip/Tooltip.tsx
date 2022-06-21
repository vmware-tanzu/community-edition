import React from 'react';
import './Tooltip.scss';

interface TooltipProps {
    children?: React.ReactNode;
}

function Tooltip(props: TooltipProps) {
    return <div className="tooltip">{props.children}</div>;
}

export default Tooltip;
