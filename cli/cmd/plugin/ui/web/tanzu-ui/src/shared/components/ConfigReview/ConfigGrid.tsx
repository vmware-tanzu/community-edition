// React imports
import React from 'react';
// App imports
import './ConfigReview.scss';

const GRIDCELLS_PER_LABEL = 2;
const GRIDCELLS_PER_VALUE = 4;
const GRIDCELLS_PER_PAIR = GRIDCELLS_PER_LABEL + GRIDCELLS_PER_VALUE;

export interface ConfigGridProps {
    groups: ConfigGroup[];
}
export interface ConfigPair {
    field: string; // used for retrieving the data from a store object
    label: string;
    value: string;
    createValueDisplay?: (value: string) => string; // used to modify the value into a more display-friendly string
}

export const CommonValueDisplayFunctions = {
    MASK: (value: string) => '*'.repeat(value ? value.length : 0),
    TRUNCATE: (nChars: number) => (value: string) => value && value.length > nChars ? `${value.substring(0, nChars - 3)}...` : value,
};

export interface ConfigGroup {
    label: string;
    pairs: ConfigPair[];
    pairsPerLine: number;
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
    const totalCellsPerLine = group.pairsPerLine * GRIDCELLS_PER_PAIR;
    const cdsLayout = `col:${totalCellsPerLine}`;
    return (
        <>
            <div cds-layout={cdsLayout} className="config-group">
                {group?.label}
            </div>
            {ConfigLinesDisplay(group.pairs, group.pairsPerLine)}
        </>
    );
}

function ConfigLinesDisplay(pairs: ConfigPair[], pairsPerLine: number) {
    const result: JSX.Element[] = [];
    const nLines = pairs.length / pairsPerLine; // NOTE: any remainder will cause an extra line in the loop below (as desired)

    for (let xLine = 0; xLine < nLines; xLine++) {
        result.push(...ConfigSingleLineDisplay(pairs.slice(xLine * pairsPerLine), pairsPerLine));
    }
    return <>{result}</>;
}

function ConfigSingleLineDisplay(pairs: ConfigPair[], pairsPerLine: number): JSX.Element[] {
    const result: JSX.Element[] = [];
    // We want to return the given number of pairs for the line. If we run out of actual pairs, we return blank elements
    for (let x = 0; x < pairsPerLine; x++) {
        result.push(ConfigPairDisplay(x < pairs.length ? pairs[x] : undefined));
    }
    return result;
}

function ConfigPairDisplay(pair: ConfigPair | undefined) {
    const cdsLayoutLabel = `col:${GRIDCELLS_PER_LABEL}`;
    const cdsLayoutValue = `col:${GRIDCELLS_PER_VALUE}`;
    const displayValue = pair?.value && pair?.createValueDisplay ? pair.createValueDisplay(pair.value) : pair?.value;
    // NOTE: the &nbsp; are there to ensure the background is painted when there is no data; a simple space does not work
    return (
        <>
            <div cds-layout={cdsLayoutLabel} className="config-label">
                {pair?.label}&nbsp;
            </div>
            <div cds-layout={cdsLayoutValue} className="config-value">
                {displayValue}&nbsp;
            </div>
        </>
    );
}
