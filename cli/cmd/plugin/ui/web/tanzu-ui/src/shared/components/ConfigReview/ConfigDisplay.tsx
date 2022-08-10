// React imports
import React, { useState } from 'react';
// Library imports
import { CdsAccordion, CdsAccordionContent, CdsAccordionHeader, CdsAccordionPanel } from '@cds/react/accordion';
// App imports
import { ConfigGrid, ConfigGroup, ConfigPair } from './ConfigGrid';
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

    const groups = transformConfigData(props.data.groups, props.store);
    const header = props.data.label;
    const about = props.data.about;

    return (
        <div className="config-section" cds-layout="m-y:sm">
            <CdsAccordion className="accordion-normal">
                <CdsAccordionPanel expanded={open} cds-motion="off" onExpandedChange={toggleOpen}>
                    <CdsAccordionHeader>{header}</CdsAccordionHeader>
                    <CdsAccordionContent>
                        <div className="config-about">{about}</div>
                        <ConfigGrid groups={groups} />
                    </CdsAccordionContent>
                </CdsAccordionPanel>
            </CdsAccordion>
        </div>
    );
}

function transformConfigData(groups: ConfigGroup[], dataObject: any): ConfigGroup[] {
    return groups.map((group) => transformConfigGroup(group, dataObject));
}

function transformConfigGroup(group: ConfigGroup, dataObject: any): ConfigGroup {
    return { ...group, pairs: group.pairs?.map((pair) => transformConfigPair(pair, dataObject)) };
}

function transformConfigPair(pair: ConfigPair, dataObject: any): ConfigPair {
    const result = { ...pair };
    // if there is no given value, but there is a field and a dataObject, use the dataObject to get a value
    if (!pair.value && pair.field && dataObject) {
        result.value = dataObject[pair.field];
    }
    return result.transform ? result.transform(result) : result;
}
