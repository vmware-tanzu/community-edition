// React imports
import React, { ChangeEvent } from 'react';
import { FieldErrors, RegisterOptions, UseFormRegisterReturn } from 'react-hook-form';
// Library imports
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';

export function ClusterNameSection(
    field: string,
    errors: FieldErrors,
    register: (name: any, options?: RegisterOptions<any, any>) => UseFormRegisterReturn,
    onEnterClusterName: (event: ChangeEvent<HTMLInputElement>) => void
) {
    return (
        <div cds-layout="vertical gap:lg gap@md:lg col@sm:6 col:6">
            <CdsFormGroup layout="vertical">
                <CdsInput layout="vertical">
                    <label>Cluster Name</label>
                    <input placeholder="cluster-name" {...register(field)} onChange={onEnterClusterName} />
                    {errors[field] && <CdsControlMessage status="error">{errors[field].message}</CdsControlMessage>}
                </CdsInput>
            </CdsFormGroup>
            <div>Can only contain lowercase alphanumeric characters and dashes. </div>
            <div>You will use this cluster name when using the Tanzu CLI and kubectl utilities.</div>
        </div>
    );
}
