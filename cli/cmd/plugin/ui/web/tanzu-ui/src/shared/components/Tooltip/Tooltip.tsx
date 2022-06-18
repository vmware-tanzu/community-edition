import React, { ReactElement } from 'react';
import './Tooltip.scss';

interface TooltipProps {
    children?: React.ReactNode;
}

function Tooltip(props: TooltipProps) {
    return <div className="container">{props.children}</div>;
}

export default Tooltip;
