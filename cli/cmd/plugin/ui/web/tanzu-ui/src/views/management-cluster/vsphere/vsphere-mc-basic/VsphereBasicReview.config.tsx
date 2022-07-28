import { CommonValueDisplayFunctions, ConfigGroup } from '../../../../shared/components/ConfigReview/ConfigGrid';
import { ConfigDisplayData } from '../../../../shared/components/ConfigReview/ConfigDisplay';
import { VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';

const configGroupsBasic: ConfigGroup[] = [
    {
        label: 'Credentials',
        pairsPerLine: 2,
        pairs: [
            { label: 'Server', value: '', field: VSPHERE_FIELDS.SERVERNAME },
            { label: 'Username', value: '', field: VSPHERE_FIELDS.USERNAME },
            { label: 'Datacenter', value: '', field: VSPHERE_FIELDS.DATACENTER },
            {
                label: 'Password',
                value: '',
                field: VSPHERE_FIELDS.PASSWORD,
                createValueDisplay: CommonValueDisplayFunctions.MASK,
            },
        ],
    },
    {
        label: 'Cluster settings',
        pairsPerLine: 2,
        pairs: [
            { label: 'Name', value: '', field: VSPHERE_FIELDS.CLUSTERNAME },
            { label: 'Node Type', value: '', field: VSPHERE_FIELDS.INSTANCETYPE },
            { label: 'OS Image', value: '', field: VSPHERE_FIELDS.VMTEMPLATE },
            {
                label: 'SSH key',
                value: '',
                field: VSPHERE_FIELDS.SSHKEY,
                createValueDisplay: CommonValueDisplayFunctions.TRUNCATE(24),
            },
        ],
    },
    {
        label: 'Load Balancer',
        pairsPerLine: 2,
        pairs: [
            { label: 'Provider', value: 'Kube-vip', field: '' },
            { label: 'Endpoint', value: '', field: VSPHERE_FIELDS.CLUSTER_ENDPOINT },
        ],
    },
    {
        label: 'Resources',
        pairsPerLine: 2,
        pairs: [
            { label: 'VM Folder', value: '', field: VSPHERE_FIELDS.VMFolder },
            { label: 'Data Store', value: '', field: VSPHERE_FIELDS.DataStore },
            { label: 'Network', value: '', field: VSPHERE_FIELDS.NetworkName },
            { label: 'Pool', value: '', field: VSPHERE_FIELDS.Pool },
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
        pairsPerLine: 2,
        pairs: [
            { label: 'Server', value: '', field: VSPHERE_FIELDS.SERVERNAME },
            { label: 'Username', value: '', field: VSPHERE_FIELDS.USERNAME },
            { label: 'Datacenter', value: '', field: VSPHERE_FIELDS.DATACENTER },
            { label: 'Password', value: '********', field: '' },
        ],
    },
];

export const configDisplayDefaults: ConfigDisplayData = {
    label: 'Configuration Defaults',
    groups: configGroupsDefault,
    about: 'These are values set "behind the scenes". If you want to change any of them, use an advanced configuration option.',
};
