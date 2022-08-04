// React imports
import React, { useState } from 'react';
// Library imports
import { CdsAccordion, CdsAccordionContent, CdsAccordionHeader, CdsAccordionPanel } from '@cds/react/accordion';
// App imports
import { ConfigGrid, ConfigGroup } from './ConfigGrid';
import './ConfigReview.scss';

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

    if (props.store) {
        populateConfigData(props.data.groups, props.store);
    }

    return (
        <div className="config-section" cds-layout="m-y:sm">
            <CdsAccordion className="accordion-normal">
                <CdsAccordionPanel expanded={open} cds-motion="off" onExpandedChange={toggleOpen}>
                    <CdsAccordionHeader>{props.data.label}</CdsAccordionHeader>
                    <CdsAccordionContent>
                        <div className="config-about">{props.data.about}</div>
                        <ConfigGrid groups={props.data.groups} />
                    </CdsAccordionContent>
                </CdsAccordionPanel>
            </CdsAccordion>
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
