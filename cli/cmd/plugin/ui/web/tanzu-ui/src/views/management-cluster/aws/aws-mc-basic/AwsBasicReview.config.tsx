// App imports
import { AWS_FIELDS } from './AwsManagementClusterBasic.constants';
import { CommonConfigTransformationFunctions, ConfigGroup } from '../../../../shared/components/ConfigReview/ConfigGrid';
import { ConfigDisplayData } from '../../../../shared/components/ConfigReview/ConfigDisplay';

const configGroupOneTimeCredentials: ConfigGroup = {
    label: 'One-time credentials',
    pairs: [
        { label: 'Access Key', field: AWS_FIELDS.SECRET_ACCESS_KEY, transform: CommonConfigTransformationFunctions.MASK },
        { label: 'Key Id', field: AWS_FIELDS.ACCESS_KEY_ID, transform: CommonConfigTransformationFunctions.MASK },
        { label: 'Session token', field: AWS_FIELDS.SESSION_TOKEN, transform: CommonConfigTransformationFunctions.MASK },
        { label: 'Region', field: AWS_FIELDS.REGION },
        { label: 'EC2 key pair', field: AWS_FIELDS.EC2_KEY_PAIR },
    ],
};

const configGroupCredentialProfile: ConfigGroup = {
    label: 'AWS credential profile',
    pairs: [
        { label: 'Profile', field: AWS_FIELDS.PROFILE },
        { label: 'Region', field: AWS_FIELDS.REGION },
        { label: 'EC2 key pair', field: AWS_FIELDS.EC2_KEY_PAIR },
    ],
};

const configGroupClusterSettings: ConfigGroup = {
    label: 'Cluster settings',
    pairs: [
        { label: 'Name', field: AWS_FIELDS.CLUSTER_NAME },
        { label: 'Node Type', field: AWS_FIELDS.NODE_PROFILE },
        {
            label: 'OS Image',
            field: AWS_FIELDS.OS_IMAGE,
            longValue: true,
            transform: (pair) => {
                return { ...pair, value: pair.value?.name ?? '' };
            },
        },
    ],
};

export const AwsConfigDisplayBasicOneTimeCredentials: ConfigDisplayData = {
    label: 'Basic Configuration',
    groups: [configGroupOneTimeCredentials, configGroupClusterSettings],
    about: 'This is the basic configuration display',
};

export const AwsConfigDisplayBasicProfileCredentials: ConfigDisplayData = {
    label: 'Basic Configuration',
    groups: [configGroupCredentialProfile, configGroupClusterSettings],
    about: 'This is the basic configuration display',
};

const configGroupsDefault: ConfigGroup[] = [
    {
        label: 'Some kinda default group',
        pairs: [
            { label: 'FIELD 1', value: 'something' },
            { label: 'FIELD 2', value: 'something else' },
        ],
    },
];

export const AwsConfigDisplayDefaults: ConfigDisplayData = {
    label: 'Configuration Defaults',
    groups: configGroupsDefault,
    about: 'These are default values that are common to most AWS clusters. If you need to change any of them, use an advanced configuration option.',
};
