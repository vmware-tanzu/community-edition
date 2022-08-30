// App imports
import { CommonConfigTransformationFunctions, ConfigGroup } from '../../../../shared/components/ConfigReview/ConfigGrid';
import { ConfigDisplayData } from '../../../../shared/components/ConfigReview/ConfigDisplay';
import { VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';

const configGroupsBasic: ConfigGroup[] = [
    {
        label: 'Credentials',
        pairs: [
            { label: 'Server', field: VSPHERE_FIELDS.SERVERNAME },
            { label: 'Username', field: VSPHERE_FIELDS.USERNAME },
            { label: 'Datacenter', field: VSPHERE_FIELDS.DATACENTER },
            {
                label: 'Password',
                field: VSPHERE_FIELDS.PASSWORD,
                transform: CommonConfigTransformationFunctions.MASK,
            },
        ],
    },
    {
        label: 'Cluster settings',
        pairs: [
            { label: 'Name', field: VSPHERE_FIELDS.CLUSTERNAME },
            { label: 'Node Type', field: VSPHERE_FIELDS.NODE_PROFILE_TYPE },
            { label: 'OS Image', field: VSPHERE_FIELDS.VMTEMPLATE, transform: CommonConfigTransformationFunctions.NAME },
            {
                label: 'SSH key',
                field: VSPHERE_FIELDS.SSHKEY,
                transform: CommonConfigTransformationFunctions.TRUNCATE(24),
            },
        ],
    },
    {
        label: 'Load Balancer',
        pairs: [
            { label: 'Provider', value: 'Kube-vip' },
            { label: 'Endpoint', field: VSPHERE_FIELDS.CLUSTER_ENDPOINT },
        ],
    },
    {
        label: 'Resources',
        pairs: [
            { label: 'VM Folder', field: VSPHERE_FIELDS.VMFolder, transform: CommonConfigTransformationFunctions.NAME },
            { label: 'Data Store', field: VSPHERE_FIELDS.DataStore, transform: CommonConfigTransformationFunctions.NAME },
            { label: 'Network', field: VSPHERE_FIELDS.Network, transform: CommonConfigTransformationFunctions.NAME },
            { label: 'Pool', field: VSPHERE_FIELDS.Pool },
        ],
    },
];
export const configDisplayBasic: ConfigDisplayData = {
    label: 'Basic Configuration',
    groups: configGroupsBasic,
    about: 'This is the basic configuration display...',
};
const configGroupsDefault: ConfigGroup[] = [
    {
        label: 'Some Kinda Default Group...',
        pairs: [
            { label: 'Server', field: VSPHERE_FIELDS.SERVERNAME },
            { label: 'Username', field: VSPHERE_FIELDS.USERNAME },
            { label: 'Datacenter', field: VSPHERE_FIELDS.DATACENTER },
        ],
    },
];

export const configDisplayDefaults: ConfigDisplayData = {
    label: 'Configuration Defaults',
    groups: configGroupsDefault,
    about: 'These are default values that are common to most clusters. If you need to change any of them, use an advanced configuration option.',
};
