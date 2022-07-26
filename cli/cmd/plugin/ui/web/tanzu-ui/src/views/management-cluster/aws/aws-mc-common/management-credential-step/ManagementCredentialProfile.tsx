// React imports
import React, { ChangeEvent, useEffect, useState } from 'react';

// Library imports
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { CdsSelect } from '@cds/react/select';
import { UseFormReturn } from 'react-hook-form';

// App imports
import { AwsService } from '../../../../../swagger-api/services/AwsService';
import { FormInputs } from './ManagementCredentials';
import './ManagementCredentialProfile.scss';

interface Props {
    initialProfile: string;
    initialRegion: string;
    regions: string[];
    handleSelectProfile: (profile: string) => void;
    handleSelectRegion: (region: string) => void;
    methods: UseFormReturn<FormInputs, any>;
}

function ManagementCredentialProfile(props: Props) {
    const {
        methods: {
            formState: { errors },
            register,
            setValue,
        },
        regions,
        initialProfile,
        initialRegion,
        handleSelectProfile,
        handleSelectRegion,
    } = props;
    const [profiles, setProfiles] = useState<string[]>([]);
    useEffect(() => {
        // fetch profiles
        AwsService.getAwsCredentialProfiles().then((data: string[]) => setProfiles(data));
        setValue('REGION', initialRegion);
    }, []); // eslint-disable-line react-hooks/exhaustive-deps

    const handleProfileChange = (event: ChangeEvent<HTMLSelectElement>) => {
        handleSelectProfile(event.target.value);
    };

    const handleRegionChange = (event: ChangeEvent<HTMLSelectElement>) => {
        setValue('REGION', event.target.value, { shouldValidate: true });
        handleSelectRegion(event.target.value);
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
                    <CdsSelect layout="compact">
                        <label>AWS credential profile</label>
                        <select
                            className="select-sm-width"
                            onChange={handleProfileChange}
                            value={initialProfile}
                            data-testid="profile-select"
                        >
                            <option></option>
                            {profiles.map((profile) => (
                                <option key={profile} value={profile}>
                                    {profile}
                                </option>
                            ))}
                        </select>
                    </CdsSelect>

                    <CdsSelect layout="compact">
                        <label>AWS region </label>
                        <select
                            className="select-sm-width"
                            {...register('REGION')}
                            onChange={handleRegionChange}
                            value={initialRegion}
                            data-testid="region-select"
                        >
                            <option></option>
                            {regions.map((region) => (
                                <option key={region} value={region}>
                                    {region}
                                </option>
                            ))}
                        </select>
                        {errors['REGION'] && <CdsControlMessage status="error">{errors['REGION'].message}</CdsControlMessage>}
                    </CdsSelect>
                </div>
            </CdsFormGroup>
        </>
    );
}

export default ManagementCredentialProfile;
