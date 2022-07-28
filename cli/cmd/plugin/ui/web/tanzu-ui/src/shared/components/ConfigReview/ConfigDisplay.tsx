// React imports
import React, { useState } from 'react';
// App imports
import { ConfigGrid, ConfigGroup } from './ConfigGrid';
import './ConfigReview.scss';
import { CdsIcon } from '@cds/react/icon';
import { ClarityIcons, resizeIcon, shrinkIcon } from '@cds/core/icon';

ClarityIcons.addIcons(shrinkIcon, resizeIcon);

export interface ConfigDisplayData {
    label: string;
    groups: ConfigGroup[];
    about?: string;
}

export interface ConfigDisplayProps {
    data: ConfigDisplayData;
    startsOpen?: boolean;
    store?: any;
}

export function ConfigDisplay(props: ConfigDisplayProps) {
    const [open, setOpen] = useState(props.startsOpen);

    const toggleOpen = () => {
        setOpen((prevState) => !prevState);
    };

    const iconName = open ? 'shrink' : 'resize';
    const iconClass = open ? 'config-icon-collapse' : 'config-icon-expand';

    if (props.store) {
        populateConfigData(props.data.groups, props.store);
    }

    return (
        <div className="config-section" cds-layout="m-y:sm">
            <CdsIcon shape={iconName} size="sm" className={iconClass} onClick={toggleOpen}></CdsIcon>
            &nbsp;
            {props.data.label}
            {open && (
                <>
                    <div className="config-about">{props.data.about}</div>
                    <ConfigGrid groups={props.data.groups} />
                </>
            )}
        </div>
    );
}

function populateConfigData(groups: ConfigGroup[], dataObject: any) {
    groups.forEach((group) => {
        group?.pairs?.forEach((pair) => {
            if (pair.field) {
                pair.value = dataObject[pair.field];
            }
        });
    });
}
