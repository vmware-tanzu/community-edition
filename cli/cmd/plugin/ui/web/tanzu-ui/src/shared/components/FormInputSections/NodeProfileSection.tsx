// React imports
import React, { ChangeEvent } from 'react';
// Library imports
import { CdsControlMessage } from '@cds/react/forms';
import { CdsIcon } from '@cds/react/icon';
import * as yup from 'yup';
import { FieldError } from 'react-hook-form';
import { CdsRadio, CdsRadioGroup } from '@cds/react/radio';
// App imports
import './NodeProfileSection.scss';

export interface NodeInstanceType {
    id: string;
    icon: string;
    isSolidIcon?: boolean;
    label: string;
    description: string;
}

/**
 * addNodeInstanceTypeValidation takes a "yup" schema object and returns a new "yup" schema object with the field added,
 * and associated with a yup validation
 * @param field - the name of the field for the node instance id
 * @param yupObject - the yup schema object
 */
export function addNodeInstanceTypeValidation(field: string, yupObject: any): any {
    return { ...yupObject, [field]: yup.string().nullable().required('Please select an instance type for your cluster nodes') };
}

function InstanceTypeInList(field: string, instance: NodeInstanceType, register: any, selectedInstanceId?: string) {
    return (
        <CdsRadio cds-layout="m:lg m-l:xl p-b:sm" key={instance.id + '-cds-radio'} data-testid="cds-radio">
            <label>
                {instance.label}
                <CdsIcon
                    shape={instance.icon}
                    size="md"
                    className={selectedInstanceId === instance.id ? 'node-icon selected' : 'node-icon'}
                    solid={instance.isSolidIcon}
                ></CdsIcon>
                <div className="radio-message">{instance.description}</div>
            </label>
            <input type="radio" key={instance.id} value={instance.id} checked={selectedInstanceId === instance.id} readOnly />
        </CdsRadio>
    );
}

const DEFAULT_PROMPT = 'Select a node profile';

export function NodeProfileSection(
    field: string,
    nodeInstanceTypes: NodeInstanceType[],
    errors: { [key: string]: FieldError | undefined },
    register: any,
    nodeInstanceTypeChange: (nodeInstanceTypeId: string, fieldName?: string) => void,
    selectedInstanceId?: string,
    prompt?: string
) {
    const onSelectNodeInstanceType = (event: ChangeEvent<HTMLSelectElement>) => {
        nodeInstanceTypeChange(event.target.value || '', field);
    };

    const fieldError = errors[field];
    const displayPrompt = prompt ? prompt : DEFAULT_PROMPT;
    return (
        <div className="cluster-node-instance-container" cds-layout="m:lg">
            <div cds-layout="col@sm:8 p-l:xl">
                <CdsRadioGroup layout="vertical" onChange={onSelectNodeInstanceType}>
                    <label>{displayPrompt}</label>
                    {nodeInstanceTypes.map((instanceType) => {
                        return InstanceTypeInList(field, instanceType, register, selectedInstanceId);
                    })}
                </CdsRadioGroup>
                {fieldError && <CdsControlMessage status="error">{fieldError.message}</CdsControlMessage>}
            </div>
        </div>
    );
}
