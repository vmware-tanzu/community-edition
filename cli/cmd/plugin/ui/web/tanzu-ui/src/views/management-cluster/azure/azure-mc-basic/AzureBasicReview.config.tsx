import { CommonValueDisplayFunctions, ConfigGroup } from '../../../../shared/components/ConfigReview/ConfigGrid';
import { ConfigDisplayData } from '../../../../shared/components/ConfigReview/ConfigDisplay';
import { OsImage } from '../../../../shared/components/FormInputComponents/OsImageSelect/OsImageSelect';
import { AZURE_FIELDS } from './AzureManagementClusterBasic.constants';

const configGroupCredentials: ConfigGroup = {
    label: 'Azure Credentials',
    pairs: [
        { label: 'Tenant Id', field: AZURE_FIELDS.TENANT_ID },
        { label: 'Client Id', field: AZURE_FIELDS.CLIENT_ID },
        { label: 'Client Secret', field: AZURE_FIELDS.CLIENT_SECRET },
        { label: 'Subscription', field: AZURE_FIELDS.SUBSCRIPTION_ID },
        { label: 'Environment', field: AZURE_FIELDS.AZURE_ENVIRONMENT },
        { label: 'Region', field: AZURE_FIELDS.REGION },
        { label: 'Public SSH Key', field: AZURE_FIELDS.SSH_PUBLIC_KEY, createValueDisplay: CommonValueDisplayFunctions.TRUNCATE(24) },
    ],
};

const configGroupClusterSettings: ConfigGroup = {
    label: 'Cluster settings',
    pairs: [
        { label: 'Name', field: AZURE_FIELDS.CLUSTER_NAME },
        { label: 'Node Type', field: AZURE_FIELDS.NODE_PROFILE },
        {
            label: 'Image (AMI)',
            field: AZURE_FIELDS.OS_IMAGE,
            longValue: true,
            createValueDisplay: (osImageInfo: OsImage) => osImageInfo?.name ?? '',
        },
    ],
};

export const AzureConfigDisplayConfig: ConfigDisplayData = {
    label: 'Basic Configuration',
    groups: [configGroupCredentials, configGroupClusterSettings],
    about: 'This is the basic configuration display',
};

const configGroupsDefault: ConfigGroup[] = [
    {
        label: 'Some kinda default group',
        pairs: [
            { label: 'FIELD A', value: 'something' },
            { label: 'FIELD B', value: 'something else' },
        ],
    },
];

export const AzureConfigDisplayDefaults: ConfigDisplayData = {
    label: 'Configuration Defaults',
    groups: configGroupsDefault,
    about: 'These are default values that are common to most Azure clusters. If you need to change any of them, use an advanced configuration option.',
};
