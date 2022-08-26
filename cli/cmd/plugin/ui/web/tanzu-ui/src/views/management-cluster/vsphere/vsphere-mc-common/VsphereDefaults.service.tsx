import { VSphereDatastore, VSphereFolder, VSphereNetwork, VSphereVirtualMachine } from '../../../../swagger-api';
import { filterTemplates } from './VsphereOsImageUtil';
import { first } from '../../../../shared/utilities/Array.util';
import { VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';

const staticDefaults = {
    [VSPHERE_FIELDS.CNI_TYPE]: 'antrea',
    [VSPHERE_FIELDS.MACHINE_HEALTH_CHECK_ACTIVATED]: false,
    [VSPHERE_FIELDS.CEIP_OPT_IN]: false,
    [VSPHERE_FIELDS.ENABLE_AUDIT_LOGGING]: false,
    [VSPHERE_FIELDS.CLUSTER_NODE_CIDR]: '',
    [VSPHERE_FIELDS.CLUSTER_SERVICE_CIDR]: '100.64.0.0/13',
    [VSPHERE_FIELDS.CLUSTER_POD_CIDR]: '100.96.0.0/11',
} as { [key: string]: any };

export class VsphereDefaultsService {
    static selectDefaultOsImage = (osImages: VSphereVirtualMachine[]) => {
        return first<VSphereVirtualMachine>(filterTemplates(osImages));
    };

    static selectDefaultNetwork = (networks: VSphereNetwork[]) => {
        return first<VSphereNetwork>(networks);
    };

    static selectDefaultDatastore = (datastores: VSphereDatastore[]) => {
        return first<VSphereDatastore>(datastores);
    };

    static selectDefaultFolder = (folders: VSphereFolder[]) => {
        return first<VSphereFolder>(folders);
    };

    static getStaticDefaults = () => {
        return { ...staticDefaults };
    };
}
