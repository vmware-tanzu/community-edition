import React, { ChangeEvent } from 'react';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsInput } from '@cds/react/input';

export function ClusterNameSection(
    field: string,
    errors: any,
    register: any,
    onEnterClusterName: (evt: ChangeEvent<HTMLSelectElement>) => void
) {
    return (
        <div cds-layout={`vertical gap:lg gap@md:lg col@sm:6 col:6`}>
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
