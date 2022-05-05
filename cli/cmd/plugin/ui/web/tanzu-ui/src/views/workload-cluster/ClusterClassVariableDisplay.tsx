// React imports
import React, { ChangeEvent } from 'react';
// Library imports
import * as yup from 'yup';
import { CdsAccordion, CdsAccordionContent, CdsAccordionHeader, CdsAccordionPanel } from '@cds/react/accordion';
import { CdsCheckbox } from '@cds/react/checkbox';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';
import { CdsSelect } from '@cds/react/select';
import { CdsTextarea } from '@cds/react/textarea';
// App imports
import { ClusterClassDefinition, ClusterClassVariable, ClusterClassVariableType } from '../../shared/models/ClusterClass';
import { isValidCidr, isValidCommaSeparatedIpOrFqdn, isValidFqdn, isValidIp } from '../../shared/validations/Validation.service';

const NCOL_DESCRIPTION = 'col:3'
const NCOL_INPUT_CONTROL = 'col:9'
export interface ClusterClassVariableDisplayOptions {
    register: any,
    errors: any,
    expanded: boolean,
    toggleExpanded: () => void,
    onValueChange: (evt: ChangeEvent<HTMLSelectElement>) => void,
}

function ClusterClassVariableInput(ccVar: ClusterClassVariable, options: ClusterClassVariableDisplayOptions) {
    switch (ccVar.valueType) {
    case ClusterClassVariableType.BOOLEAN:
        return ClusterClassVariableInputBoolean(ccVar, options)
    case ClusterClassVariableType.INTEGER:
        return ClusterClassVariableInputInteger(ccVar, options)
    case ClusterClassVariableType.STRING:
    case ClusterClassVariableType.IP:
    case ClusterClassVariableType.IP_LIST:
    case ClusterClassVariableType.CIDR:
    case ClusterClassVariableType.NUMBER:
        return ClusterClassVariableInputString(ccVar, options)
    case ClusterClassVariableType.STRING_PARAGRAPH:
        return ClusterClassVariableInputStringParagraph(ccVar, options)
    default:
        if (ccVar.valueType) {
            console.warn(`Encountered unsupported ClusterClassVariableType: ${ClusterClassVariableType[ccVar.valueType]}`)
            return <div cds-layout={NCOL_INPUT_CONTROL} className="error-text">{ccVar.name}: ClusterClassVariableInput unsupported value
                type: {ClusterClassVariableType[ccVar.valueType]} </div>
        }
        return <></>
    }
}

function ClusterClassVariableInputInteger(ccVar: ClusterClassVariable, options: ClusterClassVariableDisplayOptions) {
    return <div cds-layout={NCOL_INPUT_CONTROL}>
        <CdsFormGroup layout="vertical">
            <CdsInput layout="vertical">
                <label>{ccVar.name}</label>
                <input type="number" placeholder={ccVar.defaultValue} {...options.register(ccVar.name)} onChange={options.onValueChange} />
                { options.errors[ccVar.name] &&
                    <CdsControlMessage status="error">{options.errors[ccVar.name].message}</CdsControlMessage>
                }
            </CdsInput>
        </CdsFormGroup>
    </div>
}

function ClusterClassVariableInputString(ccVar: ClusterClassVariable, options: ClusterClassVariableDisplayOptions) {
    if (ccVar.possibleValues && ccVar.possibleValues.length > 0) {
        return ClusterClassVariableInputListbox(ccVar, options)
    }
    return <div cds-layout={NCOL_INPUT_CONTROL}>
        <CdsFormGroup layout="vertical">
            <CdsInput layout="vertical">
                <label>{ccVar.name}</label>
                <input placeholder={ccVar.defaultValue} {...options.register(ccVar.name)} onChange={options.onValueChange} />
                { options.errors[ccVar.name] &&
                    <CdsControlMessage status="error">{options.errors[ccVar.name].message}</CdsControlMessage>
                }
            </CdsInput>
        </CdsFormGroup>
    </div>
}

function ClusterClassVariableInputStringParagraph(ccVar: ClusterClassVariable, options: ClusterClassVariableDisplayOptions) {
    if (ccVar.possibleValues && ccVar.possibleValues.length > 0) {
        return ClusterClassVariableInputListbox(ccVar, options)
    }
    return <div cds-layout={NCOL_INPUT_CONTROL}>
        <CdsFormGroup layout="vertical">
            <CdsTextarea layout="vertical">
                <label>{ccVar.name}</label>
                <textarea placeholder={ccVar.defaultValue} {...options.register(ccVar.name)} onChange={options.onValueChange} ></textarea>
                { options.errors[ccVar.name] &&
                    <CdsControlMessage status="error">{options.errors[ccVar.name].message}</CdsControlMessage>
                }
            </CdsTextarea>
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
    return <div cds-layout={NCOL_INPUT_CONTROL}>
        <CdsSelect layout="compact">
            <label>{ccVar.name}</label>
            <select
                className="select-sm-width"
                {...options.register(ccVar.name)}
                onChange={options.onValueChange}
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
        <input type="checkbox" {...options.register(ccVar.name)} onChange={options.onValueChange} checked /> :
        <input type="checkbox" {...options.register(ccVar.name)} onChange={options.onValueChange} />
    return <div cds-layout={NCOL_INPUT_CONTROL}>
        <CdsFormGroup layout="vertical">
            <CdsCheckbox layout="horizontal">
                <label>{ccVar.name}</label>
                { box }
            </CdsCheckbox>
        </CdsFormGroup>
    </div>
}

function ClusterClassSingleVariableDisplay(ccVar: ClusterClassVariable, options: ClusterClassVariableDisplayOptions) {
    return <>
        <div cds-layout={NCOL_DESCRIPTION} className="text-white">{ccVar.description}</div>
        { ClusterClassVariableInput(ccVar, options) }
    </>
}

export function ClusterClassMultipleVariablesDisplay(ccVars: ClusterClassVariable[], label: string,
    options: ClusterClassVariableDisplayOptions) {
    if (!ccVars || ccVars.length === 0) {
        return <></>
    }
    return <>
        <CdsAccordion>
            <CdsAccordionPanel expanded={options.expanded} cds-motion="off" onExpandedChange={ options.toggleExpanded }>
                { innerAccordion(ccVars, label, options) }
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

function allClusterClassVariables(cc: ClusterClassDefinition): ClusterClassVariable[] {
    const a = cc.requiredVariables || []
    const b = cc.optionalVariables || []
    const c = cc.advancedVariables || []
    return a.concat(b, c)
}

export function createFormSchema(cc: ClusterClassDefinition | undefined) {
    if (!cc) {
        return undefined
    }
    const schemaObject = allClusterClassVariables(cc).reduce<any>((accumulator, ccVar) => (
        {...accumulator, [ccVar.name]: createYupObjectForClusterClassVariable(ccVar)}
    ), {})
    return yup.object(schemaObject);
}

function createYupObjectForClusterClassVariable(ccVar: ClusterClassVariable) {
    let yuppy
    switch (ccVar.valueType) {
        case ClusterClassVariableType.STRING:
            yuppy = yup.string().nullable()
            break
        case ClusterClassVariableType.BOOLEAN:
            yuppy = yup.boolean().nullable()
            break
        case ClusterClassVariableType.CIDR:
            yuppy = yup.string().test('', 'Please enter a CIDR', value => (!ccVar.required && !value) || isValidCidr(value) )
            break
        case ClusterClassVariableType.IP:
            yuppy = yup.string().test('', 'Please enter a valid ip or fqdn',
                    value => (!ccVar.required && !value) || isValidFqdn(value) || isValidIp(value))
            break
        case ClusterClassVariableType.IP_LIST:
            yuppy = yup.string().test('', 'Please enter a comma-separated list of valid ip or fqdn values',
                value => isValidCommaSeparatedIpOrFqdn(value))
            break
        default:
            yuppy = yup.string().nullable()
    }

    if (ccVar.required) {
        const prompt = errorPromptFromClusterClassType(ccVar)
        yuppy.required(prompt)
    }
    return yuppy
}

function errorPromptFromClusterClassType(ccVar: ClusterClassVariable): string {
    // NOTE: we have no need for an error prompt for BOOLEAN, because we never require a value
    switch (ccVar.valueType) {
    case ClusterClassVariableType.STRING:
        if (ccVar.possibleValues && ccVar.possibleValues.length > 0) {
            return 'Please select a value'
        }
        return 'Please enter a value'
    case ClusterClassVariableType.CIDR:
        return 'Please enter a CIDR value'
    case ClusterClassVariableType.IP:
        return 'Please enter an IP address'
    case ClusterClassVariableType.IP_LIST:
        return 'Please enter a comma-separated list of IP addresses'
    case ClusterClassVariableType.STRING_PARAGRAPH:
        return 'Please enter the required text'
    case ClusterClassVariableType.INTEGER:
        return 'Please enter a number (or use the arrows to select)'
    case ClusterClassVariableType.NUMBER:
        return 'Please enter a value'
    }
    return 'Value required'
}
