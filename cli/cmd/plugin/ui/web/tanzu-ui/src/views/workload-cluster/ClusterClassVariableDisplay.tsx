// React imports
import React, { ChangeEvent } from 'react';
// Library imports
import * as yup from 'yup';
import { CdsAccordion, CdsAccordionContent, CdsAccordionHeader, CdsAccordionPanel } from '@cds/react/accordion';
import { CdsCard } from '@cds/react/card';
import { CdsCheckbox } from '@cds/react/checkbox';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';
import { CdsSelect } from '@cds/react/select';
import { CdsTextarea } from '@cds/react/textarea';
import { CdsToggle } from '@cds/react/toggle';
// App imports
import { CCCategory, CCDefinition, CCVariable, ClusterClassVariableType } from '../../shared/models/ClusterClass';
import {
    isK8sCompliantString,
    isValidCidr,
    isValidCommaSeparatedIpOrFqdn,
    isValidFqdn,
    isValidIp,
} from '../../shared/validations/Validation.service';
import { FIELD_PATH_SEPARATOR } from '../../state-management/reducers/Form.reducer';

const NCOL_DESCRIPTION = 'col:3';
const NCOL_INPUT_CONTROL = 'col:9';
export interface ClusterClassVariableDisplayOptions {
    register: any;
    errors: any;
    expanded: boolean;
    toggleCategoryExpanded: () => void;
    onValueChangeFactory: (ccVar: CCVariable) => (evt: ChangeEvent<HTMLSelectElement>) => void;
    getFieldValue: (fieldName: string) => any;
    path?: string;
}

function CCVariableInput(ccVar: CCVariable, options: ClusterClassVariableDisplayOptions) {
    switch (ccVar.taxonomy) {
        case ClusterClassVariableType.BOOLEAN:
            return CCVariableInputBoolean(ccVar, options);
        case ClusterClassVariableType.INTEGER:
            return CCVariableInputInteger(ccVar, options);
        case ClusterClassVariableType.STRING:
        case ClusterClassVariableType.IP:
        case ClusterClassVariableType.IP_LIST:
        case ClusterClassVariableType.CIDR:
        case ClusterClassVariableType.NUMBER:
            return CCVariableInputString(ccVar, options);
        case ClusterClassVariableType.STRING_PARAGRAPH:
            return CCVariableInputStringParagraph(ccVar, options);
        default:
            if (ccVar.taxonomy) {
                console.warn(`Encountered unsupported ClusterClassVariableType: ${ccVar.taxonomy}`);
                return (
                    <div cds-layout={NCOL_INPUT_CONTROL} className="error-text">
                        {ccVar.name}: ClusterClassVariableInput unsupported value type: {ccVar.taxonomy}{' '}
                    </div>
                );
            } else {
                console.warn(`ccVar with no taxonomy: ${JSON.stringify(ccVar)}`);
            }
            return <></>;
    }
}

function CCVariableInputInteger(ccVar: CCVariable, options: ClusterClassVariableDisplayOptions) {
    const ccVarFieldName = genCCVarFieldName(ccVar.name, options.path);
    return (
        <div cds-layout={NCOL_INPUT_CONTROL}>
            <CdsFormGroup layout="vertical">
                <CdsInput layout="vertical">
                    <label>{ccVar.name}</label>
                    <input
                        type="number"
                        placeholder={ccVar.default}
                        {...options.register(ccVarFieldName)}
                        onChange={options.onValueChangeFactory(ccVar)}
                    />
                    {options.errors[ccVarFieldName] && (
                        <CdsControlMessage status="error">{options.errors[ccVarFieldName].message}</CdsControlMessage>
                    )}
                </CdsInput>
            </CdsFormGroup>
        </div>
    );
}

function CCVariableInputString(ccVar: CCVariable, options: ClusterClassVariableDisplayOptions) {
    if (ccVar.possibleValues && ccVar.possibleValues.length > 0) {
        return CCVariableInputListbox(ccVar, options);
    }
    const ccVarFieldName = genCCVarFieldName(ccVar.name, options.path);
    return (
        <div cds-layout={NCOL_INPUT_CONTROL}>
            <CdsFormGroup layout="vertical">
                <CdsInput layout="vertical" control-width="shrink">
                    <label></label>
                    <input
                        placeholder={ccVar.default}
                        {...options.register(ccVarFieldName)}
                        onChange={options.onValueChangeFactory(ccVar)}
                    />
                    {options.errors[ccVarFieldName] && (
                        <CdsControlMessage status="error">{options.errors[ccVarFieldName].message}</CdsControlMessage>
                    )}
                </CdsInput>
            </CdsFormGroup>
        </div>
    );
}

function CCVariableInputStringParagraph(ccVar: CCVariable, options: ClusterClassVariableDisplayOptions) {
    if (ccVar.possibleValues && ccVar.possibleValues.length > 0) {
        return CCVariableInputListbox(ccVar, options);
    }
    const ccVarFieldName = genCCVarFieldName(ccVar.name, options.path);
    return (
        <div cds-layout={NCOL_INPUT_CONTROL}>
            <CdsFormGroup layout="vertical">
                <CdsTextarea layout="vertical">
                    <label></label>
                    <textarea
                        placeholder={ccVar.default}
                        {...options.register(ccVarFieldName)}
                        onChange={options.onValueChangeFactory(ccVar)}
                    ></textarea>
                    {options.errors[ccVarFieldName] && (
                        <CdsControlMessage status="error">{options.errors[ccVarFieldName].message}</CdsControlMessage>
                    )}
                </CdsTextarea>
            </CdsFormGroup>
        </div>
    );
}

function CCVariableInputListbox(ccVar: CCVariable, options: ClusterClassVariableDisplayOptions) {
    const ccVarFieldName = genCCVarFieldName(ccVar.name, options.path || '');
    return (
        <div cds-layout={NCOL_INPUT_CONTROL}>
            <CdsFormGroup layout="vertical">
                <CdsSelect layout="compact" controlWidth="shrink">
                    <label></label>
                    <select
                        className="select-sm-width"
                        {...options.register(ccVarFieldName)}
                        onChange={options.onValueChangeFactory(ccVar)}
                    >
                        <option></option>
                        {ccVar.possibleValues &&
                            ccVar.possibleValues.map((value) => (
                                <option key={value} value={value}>
                                    {displayValue(value, ccVar.default)}
                                </option>
                            ))}
                    </select>
                    {options.errors[ccVarFieldName] && (
                        <CdsControlMessage status="error">{options.errors[ccVarFieldName].message}</CdsControlMessage>
                    )}
                </CdsSelect>
            </CdsFormGroup>
        </div>
    );
}

function CCVariableInputBoolean(ccVar: CCVariable, options: ClusterClassVariableDisplayOptions) {
    const ccVarFieldName = genCCVarFieldName(ccVar.name, options.path);
    const box = ccVar.default ? (
        <input type="checkbox" {...options.register(ccVarFieldName)} onChange={options.onValueChangeFactory(ccVar)} checked />
    ) : (
        <input type="checkbox" {...options.register(ccVarFieldName)} onChange={options.onValueChangeFactory(ccVar)} />
    );
    return (
        <div cds-layout={NCOL_INPUT_CONTROL}>
            <br />
            <CdsFormGroup layout="vertical">
                <CdsCheckbox layout="horizontal">
                    <label></label>
                    {box}
                </CdsCheckbox>
            </CdsFormGroup>
        </div>
    );
}

function displayValue(value: string, defaultValue: string | undefined): string {
    if (value === defaultValue) {
        return value + ' (default)';
    }
    return value;
}

function CCVariableDisplay(ccVar: CCVariable, options: ClusterClassVariableDisplayOptions) {
    return ccVar.children?.length ? CCParentVariableDisplay(ccVar, options) : CCSingleVariableDisplay(ccVar, options);
}

function CCSingleVariableDisplay(ccVar: CCVariable, options: ClusterClassVariableDisplayOptions) {
    return (
        <>
            <div cds-layout={NCOL_DESCRIPTION}>
                <br />
                <span className="cc-description-text">{ccVar.prompt}</span>
                {/*
                        <span className="cc-variable-name-text"><br/>({ccVar.name})</span>
            */}
            </div>
            {CCVariableInput(ccVar, options)}
        </>
    );
}

export function CCMultipleVariablesDisplay(ccVars: CCVariable[], ccCategory: CCCategory, options: ClusterClassVariableDisplayOptions) {
    if (!ccVars || ccVars.length === 0) {
        console.warn(`CCMultipleVariablesDisplay received empty list of vars for label ${ccCategory.label}`);
        return <></>;
    }

    const hasErrors = anyErrors(ccCategory, options.errors);
    return (
        <>
            <CdsAccordion className={hasErrors ? 'accordion-error' : 'accordion-normal'}>
                <CdsAccordionPanel expanded={options.expanded} cds-motion="off" onExpandedChange={options.toggleCategoryExpanded}>
                    {innerAccordionCC(ccCategory, options)}
                </CdsAccordionPanel>
            </CdsAccordion>
        </>
    );
}

function innerAccordionCC(ccCategory: CCCategory, options: ClusterClassVariableDisplayOptions) {
    if (!ccCategory || !ccCategory.variables || ccCategory.variables.length === 0) {
        return <></>;
    }
    // We set the path to be this category in order to namespace the fields in this category.
    // In theory, the same field could also exist in another category.
    const categoryOptions = { ...options, path: ccCategory.name };
    return (
        <>
            <CdsAccordionHeader>{ccCategory.label}</CdsAccordionHeader>
            <CdsAccordionContent>
                <div cds-layout="grid gap:lg cols:12" key="header-mc-grid">
                    {ccCategory.variables.map((ccVar: CCVariable) => CCVariableDisplay(ccVar, categoryOptions))}
                </div>
            </CdsAccordionContent>
        </>
    );
}

export function CCParentVariableDisplay(ccVar: CCVariable, options: ClusterClassVariableDisplayOptions) {
    const newDataPath = options.path ? options.path + FIELD_PATH_SEPARATOR + ccVar.name : ccVar.name;
    const newOptions = { ...options, path: newDataPath };
    const hasErrors = anyErrorsInCCVars(ccVar.children, options.errors, newDataPath);
    return ccVar.taxonomy === ClusterClassVariableType.GROUP_OPTIONAL
        ? CCPanelOptional(ccVar, hasErrors, newOptions)
        : CCPanel(ccVar, hasErrors, newOptions);
}

function CCPanelOptional(ccVar: CCVariable, hasErrors: boolean, options: ClusterClassVariableDisplayOptions) {
    const fieldName = genCCVarFieldName('__activated', options.path);
    return (
        <>
            <CdsCard className="section-raised">
                <div cds-layout="grid cols:12">
                    <div cds-layout={NCOL_DESCRIPTION} className={hasErrors ? 'error-text' : 'text-blue'}>
                        {ccVar.label}
                    </div>
                    <div cds-layout={NCOL_INPUT_CONTROL}>
                        <CdsToggle>
                            <label></label>
                            <input type="checkbox" {...options.register(fieldName)} onChange={options.onValueChangeFactory(ccVar)} />
                        </CdsToggle>
                    </div>
                </div>
                <div cds-layout="col:12 gap:lg" className="text-white">
                    &nbsp;
                </div>
                {options.getFieldValue(fieldName) && CCPanelChildren(ccVar.children || [], options)}
            </CdsCard>
        </>
    );
}

function CCPanelChildren(children: CCVariable[], options: ClusterClassVariableDisplayOptions) {
    return (
        <div cds-layout="grid gap:lg cols:12" key="panel-children-grid">
            {children?.map((ccChildVar: CCVariable) => {
                return CCSingleVariableDisplay(ccChildVar, options);
            })}
        </div>
    );
}

function CCPanel(ccVar: CCVariable, hasErrors: boolean, options: ClusterClassVariableDisplayOptions) {
    return (
        <>
            <CdsCard className="section-raised">
                <div cds-layout="col:12" className={hasErrors ? 'error-text' : 'text-blue'}>
                    {ccVar.label}
                </div>
                <div cds-layout="col:12 gap:lg" className="text-white">
                    &nbsp;
                </div>
                {CCPanelChildren(ccVar.children || [], options)}
            </CdsCard>
        </>
    );
}

export function createFormSchemaCC(cc: CCDefinition | undefined) {
    if (!cc) {
        console.error('createFormSchemaCC received an undefined CCDefinition!');
        return undefined;
    }
    if (!cc.categories?.length) {
        console.error('createFormSchemaCC received a CCDefinition with no categories!');
        return undefined;
    }
    const schemaObject = createFormSchemaFromCCDefinition(cc);
    return yup.object(schemaObject);
}

function createFormSchemaFromCCDefinition(cc: CCDefinition): any {
    return cc.categories?.reduce<any>(
        (accumulator, ccCategory: CCCategory) => addFormSchemaFromCCVars(ccCategory.variables, ccCategory.name, accumulator),
        {}
    );
}

// The form schema adds all the fields and their yup objects to a single object
function addFormSchemaFromCCVars(ccVars: CCVariable[], path: string, accumulator: any): any {
    if (!ccVars) {
        console.warn(`createFormSchemaFromCCVars received undefined ccVars, path=${path}`);
        return accumulator;
    }
    return ccVars?.reduce<any>((acc, ccVar) => addFormSchemaFromSingleVar(ccVar, path, acc), accumulator);
}

function addFormSchemaFromSingleVar(ccVar: CCVariable, path: string, accumulator: any): any {
    return ccVar.children?.length
        ? // for parent objects, we add all the children objects (but not the parent itself)
          addFormSchemaFromCCVars(ccVar.children, addToPath(path, ccVar.name), accumulator)
        : // for simple variables, we just create a yup object to associate with the variable name
          { ...accumulator, [genCCVarFieldName(ccVar.name, path)]: createYupObjectForCCVariable(ccVar) };
}

function createYupObjectForCCVariable(ccVar: CCVariable) {
    let yuppy;
    switch (ccVar.taxonomy) {
        case ClusterClassVariableType.STRING:
            yuppy = yup.string().nullable();
            break;
        case ClusterClassVariableType.STRING_K8S_COMPLIANT:
            yuppy = yup
                .string()
                .test(
                    '',
                    'Please enter a string containing only lower-case letters and hyphens',
                    (value) => (!ccVar.required && !value) || isK8sCompliantString(value)
                );
            break;
        case ClusterClassVariableType.BOOLEAN:
            yuppy = yup.boolean().nullable();
            break;
        case ClusterClassVariableType.CIDR:
            yuppy = yup.string().test('', 'Please enter a CIDR', (value) => (!ccVar.required && !value) || isValidCidr(value));
            break;
        case ClusterClassVariableType.IP:
            yuppy = yup
                .string()
                .test(
                    '',
                    'Please enter a valid ip or fqdn',
                    (value) => (!ccVar.required && !value) || isValidFqdn(value) || isValidIp(value)
                );
            break;
        case ClusterClassVariableType.IP_LIST:
            yuppy = yup
                .string()
                .test(
                    '',
                    'Please enter a comma-separated list of valid ip or fqdn values',
                    (value) => (!ccVar.required && !value) || isValidCommaSeparatedIpOrFqdn(value)
                );
            break;
        default:
            yuppy = yup.string().nullable();
    }

    if (ccVar.required) {
        const prompt = errorPromptFromCCType(ccVar);
        yuppy.required(prompt);
    }
    return yuppy;
}

function errorPromptFromCCType(ccVar: CCVariable): string {
    // NOTE: we have no need for an error prompt for BOOLEAN, because we never require a value
    switch (ccVar.taxonomy) {
        case ClusterClassVariableType.STRING:
            if (ccVar.possibleValues && ccVar.possibleValues.length > 0) {
                return 'Please select a value';
            }
            return 'Please enter a value';
        case ClusterClassVariableType.STRING_K8S_COMPLIANT:
            return 'Please enter a string containing only lower-case letters and hyphens';
        case ClusterClassVariableType.CIDR:
            return 'Please enter a CIDR value';
        case ClusterClassVariableType.IP:
            return 'Please enter an IP address';
        case ClusterClassVariableType.IP_LIST:
            return 'Please enter a comma-separated list of IP addresses';
        case ClusterClassVariableType.STRING_PARAGRAPH:
            return 'Please enter the required text';
        case ClusterClassVariableType.INTEGER:
            return 'Please enter a number (or use the arrows to select)';
        case ClusterClassVariableType.NUMBER:
            return 'Please enter a value';
    }
    return 'Value required';
}

function genCCVarFieldName(ccVarName: string, path: string | undefined): string {
    return addToPath(path, ccVarName);
}

function addToPath(oldpath: string | undefined, parentFieldName: string): string {
    return oldpath ? oldpath + FIELD_PATH_SEPARATOR + parentFieldName : parentFieldName;
}

function anyErrors(ccCategory: CCCategory, errors: any): boolean {
    return anyErrorsInCCVars(ccCategory.variables, errors, ccCategory.name);
}

function anyErrorsInCCVars(ccVars: CCVariable[] | undefined, errors: any, path: string): boolean {
    // if a ccVar has no children, then a simple check of the errors object indicates if there is an error
    // if a ccVar has children, then a check of all the children indicates if there is an error
    return (
        ccVars !== undefined &&
        ccVars.reduce<boolean>((accum, ccVar) => {
            if (ccVar.children?.length) {
                return accum || anyErrorsInCCVars(ccVar.children, errors, addToPath(path, ccVar.name));
            }
            return accum || errors[genCCVarFieldName(ccVar.name, path)];
        }, false)
    );
}
