// React imports
import React, { ChangeEvent } from 'react';

// Library imports
import { CdsControlMessage } from '@cds/react/forms';
import { CdsIcon } from '@cds/react/icon';
import { CdsRadio, CdsRadioGroup } from '@cds/react/radio';
import { useFormContext } from 'react-hook-form';
import * as yup from 'yup';
import './NodeProfile.scss';

export interface NodeInstanceType {
    id: string;
    icon: string;
    isSolidIcon?: boolean;
    label: string;
    description: string;
}

interface NodeProfileProps {
    field: string;
    nodeInstanceTypes: NodeInstanceType[];
    nodeInstanceTypeChange: (nodeInstanceTypeId: string, fieldName?: string) => void;
    selectedInstanceId?: string;
    prompt?: string;
}

const DEFAULT_PROMPT = 'Select a node profile';

/**
 * addNodeInstanceTypeValidation returns a "yup" validation to be used in the yup schema object
 * and associated with the field that is being used for the node instance type.
 */
export function nodeInstanceTypeValidation() {
    return yup.string().nullable().required('Please select an instance type for your cluster nodes');
}

function instanceTypeInList(field: string, instance: NodeInstanceType, register: any, selectedInstanceId?: string) {
    return (
        <CdsRadio cds-layout="m:lg m-l:xl p-b:sm" key={instance.id + '-cds-radio'} data-testid="cds-radio">
            <label>
                {instance.label} selectedInstanceId = {selectedInstanceId}
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

export function NodeProfile(props: NodeProfileProps) {
    const { field, nodeInstanceTypes, nodeInstanceTypeChange, selectedInstanceId, prompt } = props;
    const {
        register,
        formState: { errors },
    } = useFormContext();

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
                        return instanceTypeInList(field, instanceType, register, selectedInstanceId);
                    })}
                </CdsRadioGroup>
                {fieldError && <CdsControlMessage status="error">{fieldError.message}</CdsControlMessage>}
            </div>
        </div>
    );
}
