// React imports
import React from 'react';
// App imports
import './ConfigReview.scss';

const GRIDCELLS_PER_LABEL = 2;
const GRIDCELLS_PER_VALUE = 4;
const GRIDCELLS_PER_PAIR = GRIDCELLS_PER_LABEL + GRIDCELLS_PER_VALUE;
const PAIRS_PER_LINE = 2;

export interface ConfigGridProps {
    groups: ConfigGroup[];
}
export interface ConfigPair {
    label: string;
    field?: string; // used for retrieving the data from a store object
    value?: string; // may be hard-coded or dynamically retrieved from store object using field attribute
    createValueDisplay?: (value: any) => string; // used to modify the value into a more display-friendly string
    longValue?: boolean; // if a long value is expected, we use 3 columns to display it
}

export const CommonValueDisplayFunctions = {
    MASK: (value: string) => '*'.repeat(value ? value.length : 0),
    TRUNCATE: (nChars: number) => (value: string) => value && value.length > nChars ? `${value.substring(0, nChars - 3)}...` : value,
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
    const totalCellsPerLine = PAIRS_PER_LINE * GRIDCELLS_PER_PAIR;
    const cdsLayout = `col:${totalCellsPerLine}`;
    return (
        <>
            <div cds-layout={cdsLayout} className="config-group">
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
        if (pairs[0].longValue) {
            result.push(ConfigLongValueDisplay(pairs[0]));
        } else {
            result.push(ConfigPairDisplay(pairs[0]));
            result.push(ConfigPairDisplay(undefined));
        }
    }

    return result;
}

function ConfigPairDisplay(pair: ConfigPair | undefined) {
    const displayValue = pair?.value && pair?.createValueDisplay ? pair.createValueDisplay(pair.value) : pair?.value;
    // NOTE: the &nbsp; are there to ensure the background is painted when there is no data; a simple space does not work
    return (
        <div cds-layout="col:6" className="config-pair">
            <div cds-layout="grid align:horizontal-stretch">
                <div cds-layout="col:4" className="config-label">
                    {pair?.label ? pair.label + ':' : ''}&nbsp;
                </div>
                <div cds-layout="col:8" className="config-value">
                    {displayValue}&nbsp;
                </div>
            </div>
        </div>
    );
}

function ConfigLongValueDisplay(pair: ConfigPair | undefined) {
    const displayValue = pair?.value && pair?.createValueDisplay ? pair.createValueDisplay(pair.value) : pair?.value;
    // NOTE: the &nbsp; are there to ensure the background is painted when there is no data; a simple space does not work
    return (
        <div cds-layout="col:12" className="config-pair">
            <div cds-layout="grid align:horizontal-stretch">
                <div cds-layout="col:2" className="config-label">
                    {pair?.label ? pair.label + ':' : ''}&nbsp;
                </div>
                <div cds-layout="col:10" className="config-value">
                    {displayValue}&nbsp;
                </div>
            </div>
        </div>
    );
}
