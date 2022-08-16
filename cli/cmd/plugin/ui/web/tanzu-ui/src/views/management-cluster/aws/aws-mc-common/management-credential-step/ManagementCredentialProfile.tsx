// React imports
import React, { ChangeEvent, useContext, useEffect, useState } from 'react';

// Library imports
import { CdsFormGroup } from '@cds/react/forms';
import { useFormContext } from 'react-hook-form';

// App imports
import { AwsService } from '../../../../../swagger-api/services/AwsService';
import { AwsStore } from '../../store/Aws.store.mc';
import { AWS_FIELDS } from '../../aws-mc-basic/AwsManagementClusterBasic.constants';
import { FormAction } from '../../../../../shared/types/types';
import { INPUT_CHANGE } from '../../../../../state-management/actions/Form.actions';
import SpinnerSelect from '../../../../../shared/components/Select/SpinnerSelect';
import './ManagementCredentialProfile.scss';

interface Props {
    selectCallback: () => void;
}

function ManagementCredentialProfile(props: Props) {
    const { selectCallback } = props;
    const { awsDispatch } = useContext(AwsStore);
    const {
        register,
        formState: { errors },
    } = useFormContext();

    const [profiles, setProfiles] = useState<string[]>([]);
    const [profilesLoading, setProfilesLoading] = useState(false);
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
    useEffect(() => {
        // fetch profiles
        const fetchProfiles = async () => {
            try {
                setProfilesLoading(true);
                const data = await AwsService.getAwsCredentialProfiles();
                setProfiles(data);
            } catch (e: any) {
                console.log(`Unabled to get profiles: ${e}`);
            } finally {
                setProfilesLoading(false);
            }
        };
        fetchProfiles();
    }, []);

    const onSelectChange = (field: string, value: string) => {
        awsDispatch({
            type: INPUT_CHANGE,
            field,
            payload: value,
        } as FormAction);
    };

    return (
        <>
            <p cds-layout="m-y:lg" className="description">
                Select an already existing AWS credential profile. The access keys and session token information configured for your profile
                will be temporarily passed to the installer.
            </p>
            <p cds-layout="m-y:lg" className="description">
                Don&apos;t have an AWS Credential profile? Credential profiles can be configured using the{' '}
                <a
                    target="_blank"
                    rel="noreferrer"
                    href="https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials"
                    className="text-blue"
                >
                    AWS CLI
                </a>
                . More on{' '}
                <a
                    target="_blank"
                    rel="noreferrer"
                    href="https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html#cli-configure-quickstart-precedence"
                    className="text-blue"
                >
                    AWS Credential profile.
                </a>
            </p>
            <CdsFormGroup layout="vertical-inline" control-width="shrink">
                <div cds-layout="horizontal gap:lg align:vertical-center">
                    <SpinnerSelect
                        label="AWS credential profile"
                        className="select-sm-width"
                        handleSelect={(e: ChangeEvent<HTMLSelectElement>) => {
                            onSelectChange(AWS_FIELDS.PROFILE, e.target.value);
                            selectCallback();
                        }}
                        name={AWS_FIELDS.PROFILE}
                        isLoading={profilesLoading}
                        register={register}
                        error={errors[AWS_FIELDS.PROFILE]?.message}
                    >
                        <option></option>
                        {profiles.map((profile) => (
                            <option key={profile} value={profile}>
                                {profile}
                            </option>
                        ))}
                    </SpinnerSelect>
                    <SpinnerSelect
                        label="AWS region"
                        className="select-sm-width"
                        handleSelect={(e: ChangeEvent<HTMLSelectElement>) => {
                            onSelectChange(AWS_FIELDS.REGION, e.target.value);
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
            </CdsFormGroup>
        </>
    );
}

export default ManagementCredentialProfile;
