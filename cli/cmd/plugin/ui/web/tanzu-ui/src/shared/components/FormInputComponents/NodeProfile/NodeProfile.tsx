// React imports
import React, { ChangeEvent } from 'react';

// Library imports
import { CdsControlMessage } from '@cds/react/forms';
import { CdsIcon } from '@cds/react/icon';
import { CdsRadio, CdsRadioGroup } from '@cds/react/radio';
import { useFormContext } from 'react-hook-form';
import * as yup from 'yup';
import './NodeProfile.scss';

export interface NodeProfileType {
    id: string;
    icon: string;
    isSolidIcon?: boolean;
    label: string;
    description: string;
}

interface NodeProfileProps {
    field: string;
    nodeProfileTypes: NodeProfileType[];
    nodeProfileTypeChange: (nodeProfileId: string, fieldName?: string) => void;
    selectedProfileId?: string;
    prompt?: string;
}

const DEFAULT_PROMPT = 'Select a node profile';

/**
 * nodeProfileValidation returns a "yup" validation to be used in the yup schema object
 * and associated with the field that is being used for the node profile type.
 */
export function nodeProfileValidation() {
    return yup.string().nullable().required('Please select an node profile for your cluster nodes');
}

function profileTypeInList(field: string, nodeProfile: NodeProfileType, register: any, selectedProfileId?: string) {
    return (
        <CdsRadio
            cds-layout="horizontal align:stretch m-t:sm m-b:xl m-l:sm p-b:sm"
            className="cds-radio-item"
            key={nodeProfile.id + '-cds-radio'}
            data-testid="cds-radio"
        >
            <label cds-layout="horizontal align:horizontal-stretch">
                {nodeProfile.label}
                <CdsIcon
                    shape={nodeProfile.icon}
                    size="md"
                    className={selectedProfileId === nodeProfile.id ? 'node-icon selected' : 'node-icon'}
                    solid={nodeProfile.isSolidIcon}
                ></CdsIcon>
                <div className="radio-message">{nodeProfile.description}</div>
            </label>
            <input type="radio" key={nodeProfile.id} value={nodeProfile.id} checked={selectedProfileId === nodeProfile.id} readOnly />
        </CdsRadio>
    );
}

export function NodeProfile(props: NodeProfileProps) {
    const { field, nodeProfileTypes, nodeProfileTypeChange, selectedProfileId, prompt } = props;
    const {
        register,
        formState: { errors },
    } = useFormContext();

    const onSelectNodeProfileType = (event: ChangeEvent<HTMLSelectElement>) => {
        nodeProfileTypeChange(event.target.value || '', field);
    };

    const fieldError = errors[field];
    const displayPrompt = prompt ? prompt : DEFAULT_PROMPT;
    return (
        <div className="cluster-node-profile-container" cds-layout="horizontal align:stretch">
            <div>
                <CdsRadioGroup layout="vertical" onChange={onSelectNodeProfileType}>
                    <label cds-layout="">{displayPrompt}</label>
                    {nodeProfileTypes.map((profileType) => {
                        return profileTypeInList(field, profileType, register, selectedProfileId);
                    })}
                </CdsRadioGroup>
                {fieldError && <CdsControlMessage status="error">{fieldError.message}</CdsControlMessage>}
            </div>
        </div>
    );
}
