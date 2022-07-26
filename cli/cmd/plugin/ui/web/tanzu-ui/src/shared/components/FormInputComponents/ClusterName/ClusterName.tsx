// React imports
import React, { ChangeEvent } from 'react';
import { FieldError, RegisterOptions, UseFormRegisterReturn } from 'react-hook-form';
// Library imports
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';
import * as yup from 'yup';
// App imports
import { isK8sCompliantString } from '../../../validations/Validation.service';
import './ClusterName.scss';

interface ClusterNameProps {
    field: string;
    errors: { [key: string]: FieldError | undefined };
    register: (name: any, options?: RegisterOptions<any, any>) => UseFormRegisterReturn;
    clusterNameChange: (clusterName: string, fieldName?: string) => void;
    placeholderClusterName?: string;
}

/**
 * addClusterNameValidation returns a "yup" validation to be used in the yup schema object
 * and associated with the field that is being used for the cluster name.
 */
export function clusterNameValidation() {
    return yup
        .string()
        .nullable()
        .required('Please enter a name for your cluster')
        .test(
            '',
            'The cluster name should contain only lower case alphanumeric characters and hyphens, beginning and ending with a non-hyphen',
            (value) => value !== null && isK8sCompliantString(value)
        );
}

// NOTE: design choice is for the clusterNameChange callback to have two parameters: the changed cluster name AND the field name.
// This is to make it easier for the parent component to use the same callback for a number of different fields; they all take the data AND
// the field name and it's possible the parent component will process all fields in a like manner.
/**
 * ClusterNameSection provides a DOM representation that allows the user to enter a cluster name
 * @param field - name of the field that holds the user input for the cluster name
 * @param errors - an errors object that indicates whether there is an error associated with this field
 * @param register - a callback used for tracking the field
 * @param clusterNameChange - function called whenever the cluster name changes
 * @param defaultClusterName - optional value of the field when user first comes to the page
 */
export function ClusterName(props: ClusterNameProps) {
    const { field, errors, register, clusterNameChange, placeholderClusterName } = props;

    const onClusterNameChange = (event: ChangeEvent<HTMLInputElement>) => {
        clusterNameChange(event.target.value || '', field);
    };

    const fieldError = errors[field];
    return (
        <div className="cluster-name-container" cds-layout="m:lg">
            <CdsFormGroup layout="vertical-inline">
                <div cds-layout="horizontal gap:lg align:vertical-center p-b:sm">
                    <CdsInput layout="compact">
                        <label cds-layout="p-b:xs">Cluster name</label>
                        <input type="text" {...register(field)} onChange={onClusterNameChange} placeholder={placeholderClusterName} />
                        {fieldError && <CdsControlMessage status="error">{fieldError.message}</CdsControlMessage>}
                    </CdsInput>
                </div>
            </CdsFormGroup>
            <div className="description" cds-layout="m-t:sm">
                Cluster name can only contain lowercase alphanumeric characters and dashes.
            </div>
            <div className="description" cds-layout="m-t:sm">
                You will reference this cluster name when using the Tanzu CLI and kubectl utilities.
            </div>
        </div>
    );
}
