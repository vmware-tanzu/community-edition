import { VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';

export enum VSPHERE_NODE_PROFILES {
    SINGLE_NODE = 'single-node',
    HIGH_AVAILABILITY = 'high-availability',
    PRODUCTION_READY = 'compute-optimized',
}

const controlPlaneFlavors = {
    [VSPHERE_NODE_PROFILES.SINGLE_NODE]: 'dev',
    [VSPHERE_NODE_PROFILES.HIGH_AVAILABILITY]: 'prod',
    [VSPHERE_NODE_PROFILES.PRODUCTION_READY]: 'prod',
};

const controlPlaneNodeTypes = {
    [VSPHERE_NODE_PROFILES.SINGLE_NODE]: 'small',
    [VSPHERE_NODE_PROFILES.HIGH_AVAILABILITY]: 'large',
    [VSPHERE_NODE_PROFILES.PRODUCTION_READY]: 'large',
};

const workerNodeTypes = {
    [VSPHERE_NODE_PROFILES.SINGLE_NODE]: 'small',
    [VSPHERE_NODE_PROFILES.HIGH_AVAILABILITY]: 'large',
    [VSPHERE_NODE_PROFILES.PRODUCTION_READY]: 'large',
};

export function createNodeProfileFieldUpdateObject(nodeProfileId: VSPHERE_NODE_PROFILES) {
    return {
        [VSPHERE_FIELDS.CONTROL_PLANE_FLAVOR]: controlPlaneFlavors[nodeProfileId],
        [VSPHERE_FIELDS.CONTROL_PLANE_INSTANCE_TYPE]: controlPlaneNodeTypes[nodeProfileId],
        [VSPHERE_FIELDS.WORKER_INSTANCE_TYPE]: workerNodeTypes[nodeProfileId],
    };
}
