export const enum DeploymentTypes {
    MANAGEMENT_CLUSTER = 'management-cluster',
    WORKLOAD_CLUSTER = 'workload-cluster',
    UNMANAGED_CLUSTER = 'unmanaged-cluster',
}

export const enum DeploymentStates {
    FAILED = 'failed',
    RUNNING = 'running',
    SUCCESSFUL = 'successful',
}
