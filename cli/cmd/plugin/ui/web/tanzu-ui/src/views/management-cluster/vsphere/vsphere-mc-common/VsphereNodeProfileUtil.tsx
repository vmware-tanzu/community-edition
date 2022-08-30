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

const workerNodeTypes = {
    [VSPHERE_NODE_PROFILES.SINGLE_NODE]: 'small',
    [VSPHERE_NODE_PROFILES.HIGH_AVAILABILITY]: 'large',
    [VSPHERE_NODE_PROFILES.PRODUCTION_READY]: 'large',
};

export function controlPlaneFlavorFromNodeProfile(nodeProfileId: VSPHERE_NODE_PROFILES): string {
    return controlPlaneFlavors[nodeProfileId];
}

export function workerNodeTypeFromNodeProfile(nodeProfileId: VSPHERE_NODE_PROFILES): string {
    return workerNodeTypes[nodeProfileId];
}
