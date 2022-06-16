import { ClarityIcons, pinIcon, timesIcon } from '@cds/core/icon';
import { CdsIcon } from '@cds/react/icon';
import React, { MouseEvent } from 'react';
import { Direction } from './Drawer.enum';
import './Drawer.scss';

ClarityIcons.addIcons(timesIcon, pinIcon);

const Drawer = ({
    direction,
    open,
    pinned,
    onClose,
    togglePin,
    children,
}: {
    direction: Direction;
    open: boolean;
    pinned: boolean;
    onClose: (event: MouseEvent<HTMLElement>) => void;
    togglePin: (event: MouseEvent<HTMLElement>) => void;
    children: any;
}) => {
    const drawerClassNames = (open: boolean, pinned: boolean) => {
        if (open && pinned) {
            return `drawer-container open pinned ${direction}`;
        } else if (open) {
            return `drawer-container open ${direction}`;
        } else if (pinned) {
            return `drawer-container pinned ${direction}`;
        } else {
            return `drawer-container ${direction}`;
        }
    };

    const classNames = drawerClassNames(open, pinned);
    return (
        <div className={classNames}>
            {open && !pinned && <div className="mask h-full w-full" id="mask" onClick={onClose}></div>}

            <div className="drawer-content h-full" cds-layout="vertical">
                <header className="drawer-header" cds-layout="horizontal p:md">
                    <CdsIcon
                        aria-label="drawer-toggle-pin"
                        className="icon drawer-pin"
                        shape="pin"
                        size="md"
                        solid={pinned}
                        onClick={togglePin}
                    ></CdsIcon>
                    <CdsIcon
                        aria-label="drawer-close"
                        cds-layout="align:right"
                        className="icon drawer-close"
                        shape="times"
                        size="md"
                        onClick={onClose}
                    ></CdsIcon>
                </header>

                <div aria-label="drawer-content-body" className="drawer-content-body h-full" cds-layout="vertical align:stretch">
                    {children}
                </div>
            </div>
        </div>
    );
};

export default Drawer;
