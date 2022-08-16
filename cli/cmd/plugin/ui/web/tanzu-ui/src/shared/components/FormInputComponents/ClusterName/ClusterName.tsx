// React imports
import React from 'react';
// Library imports
import { CdsFormGroup } from '@cds/react/forms';
import * as yup from 'yup';
// App imports
import { isK8sCompliantString } from '../../../validations/Validation.service';
import './ClusterName.scss';
import TextInputWithError from '../../Input/TextInputWithError';
import { AZURE_FIELDS } from '../../../../views/management-cluster/azure/azure-mc-basic/AzureManagementClusterBasic.constants';

interface ClusterNameProps {
    field: string;
    clusterNameChange: (clusterNameValue: string) => void;
    placeholderClusterName?: string;
    defaultClusterName?: string;
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
 * @param register - a callback used for tracking the field
 * @param clusterNameChange - function called whenever the cluster name changes
 * @param defaultClusterName - optional value of the field when user first comes to the page
 * @param placeholderClusterName -optional value of the placeholder
 */
export function ClusterName(props: ClusterNameProps) {
    const { clusterNameChange, placeholderClusterName, defaultClusterName, field } = props;
    const onClusterNameChange = (field: string, value: string) => {
        clusterNameChange(value);
    };
    return (
        <div className="cluster-name-container">
            <CdsFormGroup layout="vertical-inline">
                <div cds-layout="horizontal gap:lg align:vertical-center p-b:sm">
                    <TextInputWithError
                        label="Cluster name"
                        name={field}
                        handleInputChange={onClusterNameChange}
                        placeholder={placeholderClusterName}
                        defaultValue={defaultClusterName}
                    />
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
