// React imports
import React from 'react';
// Library imports
import { CdsIcon } from '@cds/react/icon';
// App imports
import './ConfigReview.scss';
import '../../../scss/utils.scss';
import { TooltipDefault } from '../Tooltip/TooltipDefault';

export interface ConfigGridProps {
    groups: ConfigGroup[];
}

export type ConfigDisplayFunction = (pair: ConfigPair) => ConfigPair;

export interface ConfigPair {
    label: string;
    field?: string; // used for retrieving the data from a store object
    value?: any; // may be hard-coded or dynamically retrieved from store object using field attribute
    transform?: ConfigDisplayFunction; // used to modify the pair data into a more display-friendly data
    longValue?: boolean; // if a long value is expected, we use 3 columns to display it
    tooltip?: string;
}

export const CommonConfigTransformationFunctions = {
    MASK: (pair: ConfigPair) => {
        const value = '*'.repeat(pair?.value ? pair.value.length : 0);
        return { ...pair, value };
    },
    TRUNCATE: (nChars: number) => (pair: ConfigPair) => {
        const result = { ...pair };
        if (!pair.value || pair.value.length <= nChars) {
            return result;
        }
        if (!pair.tooltip) {
            result.tooltip = pair.value;
        }
        result.value = `${pair.value.substring(0, nChars - 3)}...`;
        return result;
    },
    // The NAME transformer is used when the value is an object that has an attribute "name"
    NAME: (pair: ConfigPair) => {
        const value = pair?.value?.name || '';
        return { ...pair, value };
    },
};

export interface ConfigGroup {
    label: string;
    pairs: ConfigPair[];
}

export function ConfigGrid(props: ConfigGridProps) {
    const groups = props.groups;
    return (
        <div cds-layout="grid cols:12" className="config-grid">
            {groups.map((group) => ConfigGroupDisplay(group))}
        </div>
    );
}

function ConfigGroupDisplay(group: ConfigGroup) {
    return (
        <>
            <div cds-layout="col:12" className="config-group">
                {group?.label}
            </div>
            {ConfigLinesDisplay(group.pairs)}
        </>
    );
}

function ConfigLinesDisplay(pairs: ConfigPair[]) {
    const result: JSX.Element[] = [];

    for (let xPair = 0; xPair < pairs.length; xPair++) {
        const thisPair = pairs[xPair];
        if (xPair === pairs.length - 1) {
            // Last pair of values, so it should be on its own line
            result.push(...ConfigSingleLineDisplay([thisPair]));
        } else {
            const nextPair = pairs[xPair + 1];
            // if this pair has a long value, it should be on its own line;
            // if the NEXT pair has a long value, this pair should also be on its own line (to allow next pair its own line)
            if (thisPair.longValue || nextPair.longValue) {
                result.push(...ConfigSingleLineDisplay([thisPair]));
            } else {
                result.push(...ConfigSingleLineDisplay([thisPair, nextPair]));
                xPair++;
            }
        }
    }
    return <>{result}</>;
}

function ConfigSingleLineDisplay(pairs: ConfigPair[]): JSX.Element[] {
    const result: JSX.Element[] = [];
    if (pairs.length === 2) {
        result.push(ConfigPairDisplay(pairs[0]));
        result.push(ConfigPairDisplay(pairs[1]));
    } else if (pairs.length === 1) {
        result.push(ConfigPairDisplay(pairs[0]));
        if (!pairs[0].longValue) {
            result.push(ConfigPairDisplay(undefined));
        }
    }

    return result;
}

function ConfigPairDisplay(pair: ConfigPair | undefined) {
    const displayValue = pair?.value ?? '';
    const displayLabel = pair?.label ?? '';
    let outerDivCols = 6;
    let labelDivCols = 6;
    let valueDivCols = 6;
    if (pair?.longValue) {
        outerDivCols = 12;
        labelDivCols = 3;
        valueDivCols = 9;
    }

    if (typeof displayValue !== 'string') {
        console.error(`ConfigPairDisplay detects non-string display value for ${pair?.label}: ${JSON.stringify(displayValue)}`);
    }

    // NOTE: the blank() ensure the background is painted when there is no data; a simple space does not work
    return (
        <div cds-layout={`col:${outerDivCols}`} className="config-pair">
            <div cds-layout="grid align:horizontal-stretch">
                <div cds-layout={`col:${labelDivCols}`} className="config-label">
                    {displayLabel ? displayLabel + ':' : blank()}
                </div>
                <div cds-layout={`col:${valueDivCols}`} className="config-value">
                    {displayValue ?? blank()}
                    {pair?.tooltip && <TooltipDefault tooltipMessage={pair.tooltip} positionClass="top-left" iconSize="sm" />}
                    {!pair?.tooltip && <CdsIcon className="hide-me" shape="info-circle" size="sm"></CdsIcon>}
                </div>
            </div>
        </div>
    );
}

function blank() {
    return <div>&nbsp;</div>;
}
