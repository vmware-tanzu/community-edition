// React imports
import React, { ChangeEvent, useContext, useEffect, useState } from 'react';

// Library imports
import { CdsFormGroup } from '@cds/react/forms';
import { useFormContext } from 'react-hook-form';

// App imports
import './ManagementCredentialOneTime.scss';
import { AwsStore } from '../../../../../state-management/stores/Store.aws';
import TextInputWithError from '../../../../../shared/components/Input/TextInputWithError';
import { STORE_SECTION_FORM } from '../../../../../state-management/reducers/Form.reducer';
import { AWS_FIELDS } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
import SpinnerSelect from '../../../../../shared/components/Select/SpinnerSelect';
import { AwsService } from '../../../../../swagger-api';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import { FormAction } from '../../../../../shared/types/types';

interface Props {
    selectCallback: () => void;
}

function ManagementCredentialOneTime(props: Props) {
    const { selectCallback } = props;
    const { awsState, awsDispatch } = useContext(AwsStore);
    const {
        register,
        formState: { errors },
    } = useFormContext();

    const [regions, setRegions] = useState<string[]>([]);
    const [regionLoading, setRegionLoading] = useState(false);

    useEffect(() => {
        // fetch regions
        const fetchRegions = async () => {
            try {
                setRegionLoading(true);
                const data = await AwsService.getAwsRegions();
                setRegions(data);
            } catch (e: any) {
                console.log(`Unabled to get regions: ${e}`);
            } finally {
                setRegionLoading(false);
            }
        };
        fetchRegions();
    }, []);

    const onInputValueChange = (field: string, value: string) => {
        awsDispatch({
            type: INPUT_CHANGE,
            field,
            payload: value,
        } as FormAction);
    };
    return (
        <div className="credential-one-time-container">
            <p cds-layout="m-y:lg">
                Enter AWS account credentials directly in the Access Key ID and Secret Access Key fields for your Amazon Web Services
                account. Optionally specify an AWS session token in Session Token if your AWS account is configured to require temporary
                credentials.
            </p>
            <CdsFormGroup layout="vertical-inline" control-width="shrink">
                <div cds-layout="horizontal gap:lg align:vertical-top">
                    <div cds-layout="vertical gap:lg align:vertical-center">
                        <TextInputWithError
                            defaultValue={awsState[STORE_SECTION_FORM][AWS_FIELDS.SECRET_ACCESS_KEY]}
                            label="Secret access key"
                            name={AWS_FIELDS.SECRET_ACCESS_KEY}
                            handleInputChange={onInputValueChange}
                            maskText
                        />
                        <TextInputWithError
                            defaultValue={awsState[STORE_SECTION_FORM][AWS_FIELDS.ACCESS_KEY_ID]}
                            label="Access key ID"
                            name={AWS_FIELDS.ACCESS_KEY_ID}
                            handleInputChange={onInputValueChange}
                            maskText
                        />
                        <TextInputWithError
                            defaultValue={awsState[STORE_SECTION_FORM][AWS_FIELDS.SESSION_TOKEN]}
                            label="Session token"
                            name={AWS_FIELDS.SESSION_TOKEN}
                            handleInputChange={onInputValueChange}
                            maskText
                        />
                    </div>
                    <div cds-layout="vertical gap:lg align:vertical-top">
                        <SpinnerSelect
                            label="AWS region"
                            className="select-sm-width"
                            handleSelect={(e: ChangeEvent<HTMLSelectElement>) => {
                                onInputValueChange(AWS_FIELDS.REGION, e.target.value);
                                selectCallback();
                            }}
                            name={AWS_FIELDS.REGION}
                            isLoading={regionLoading}
                            register={register}
                            error={errors[AWS_FIELDS.REGION]?.message}
                        >
                            <option></option>
                            {regions.map((region) => (
                                <option key={region} value={region}>
                                    {region}
                                </option>
                            ))}
                        </SpinnerSelect>
                    </div>
                </div>
            </CdsFormGroup>
        </div>
    );
}

export default ManagementCredentialOneTime;
