// React imports
import React, { useEffect, useState } from 'react';

// App imports
import { CdsAccordion, CdsAccordionContent, CdsAccordionHeader, CdsAccordionPanel } from '@cds/react/accordion';
import { CdsCheckbox } from '@cds/react/checkbox';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';
import { CdsSelect } from '@cds/react/select';
import { ClusterClassVariable, ClusterClassVariableType } from '../../shared/models/ClusterClass';

export interface ClusterClassVariableDisplayOptions {
    register: any,
    errors: any,
    expanded: boolean,
    toggleExpanded: () => void,
}

function ClusterClassVariableInput(ccVar: ClusterClassVariable, options: ClusterClassVariableDisplayOptions) {
    switch (ccVar.valueType) {
        case ClusterClassVariableType.BOOLEAN:
            return ClusterClassVariableInputBoolean(ccVar, options)
        case ClusterClassVariableType.INTEGER:
            return <div>not supporting INTEGER</div>
        case ClusterClassVariableType.NUMBER:
            return <div>not supporting NUMBER</div>
        case ClusterClassVariableType.STRING:
            return ClusterClassVariableInputString(ccVar, options)
        default:
            if (ccVar.valueType) {
                return <div>unsupported value type: {ccVar.valueType}</div>
            }
            return <></>
    }
}

function ClusterClassVariableInputString(ccVar: ClusterClassVariable, options: ClusterClassVariableDisplayOptions) {
    if (ccVar.possibleValues && ccVar.possibleValues.length > 0) {
        return ClusterClassVariableInputListbox(ccVar, options)
    }
    return <div cds-layout="col:6">
        <CdsFormGroup layout="vertical">
            <CdsInput layout="vertical">
                <label>{ccVar.name}</label>
                <input placeholder={ccVar.defaultValue} {...options.register(ccVar.name)} />
                { options.errors[ccVar.name] &&
                    <CdsControlMessage status="error">{options.errors[ccVar.name].message}</CdsControlMessage>
                }
            </CdsInput>
        </CdsFormGroup>
    </div>
}

function displayValue(value: string, defaultValue: string | undefined): string {
    if (value === defaultValue) {
        return value + ' (default)'
    }
    return value
}

function ClusterClassVariableInputListbox(ccVar: ClusterClassVariable, options: ClusterClassVariableDisplayOptions) {
    return <div cds-layout="col:6">
    <CdsSelect layout="compact">
        <label>{ccVar.name}</label>
        <select
            className="select-sm-width"
            {...options.register(ccVar.name)}
/*
            value={ccVar.defaultValue}
*/
        >
            <option></option>
            { ccVar.possibleValues && ccVar.possibleValues.map((value) => (
                <option key={value} value={value}> {displayValue(value, ccVar.defaultValue)} </option>
            ))}
        </select>
        { options.errors[ccVar.name] &&
            <CdsControlMessage status="error">{options.errors[ccVar.name].message}</CdsControlMessage>
        }
    </CdsSelect>
    </div>
}

function ClusterClassVariableInputBoolean(ccVar: ClusterClassVariable, options: ClusterClassVariableDisplayOptions) {
    const box = ccVar.defaultValue ?
        <input type="checkbox" {...options.register(ccVar.name)} checked /> :
        <input type="checkbox" {...options.register(ccVar.name)} />
    return <div cds-layout="col:6">
        <CdsFormGroup layout="vertical">
            <CdsCheckbox layout="horizontal" >
                <label>{ccVar.name}</label>
                { box }
            </CdsCheckbox>
        </CdsFormGroup>
    </div>
}

function ClusterClassSingleVariableDisplay(ccVar: ClusterClassVariable, options: ClusterClassVariableDisplayOptions) {
    return <>
        {ClusterClassVariableInput(ccVar, options)}
        <div cds-layout="col:6" className="text-white">{ccVar.description}</div>
    </>
}

export function ClusterClassMultipleVariablesDisplay(ccVars: ClusterClassVariable[], label: string, options: ClusterClassVariableDisplayOptions) {
    if (!ccVars || ccVars.length === 0) {
        return <></>
    }
    return <>
        <CdsAccordion>
            <CdsAccordionPanel expanded={options.expanded} cds-motion="off" onExpandedChange={ options.toggleExpanded }>
                {innerAccordion(ccVars, label, options) }
            </CdsAccordionPanel>
        </CdsAccordion>
    </>
}

function innerAccordion(ccVars: ClusterClassVariable[], label: string, options: ClusterClassVariableDisplayOptions) {
    if (!ccVars || ccVars.length === 0) {
        return <></>
    }
    return  <>
                <CdsAccordionHeader>{label}</CdsAccordionHeader>
                <CdsAccordionContent>
                        <div cds-layout="grid gap:lg cols:12" key="header-mc-grid">
                            { ccVars.map((ccVar: ClusterClassVariable) => {
                                return ClusterClassSingleVariableDisplay(ccVar, options)
                            })
                            }
                        </div>
                </CdsAccordionContent>
    </>
}
