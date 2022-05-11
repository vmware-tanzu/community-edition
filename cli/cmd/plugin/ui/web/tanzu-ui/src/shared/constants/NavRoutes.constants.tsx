export const NavRoutes = {
    // app general
    WELCOME: '/',
    GETTING_STARTED: 'getting-started',
    MANAGEMENT_CLUSTER_LANDING: '/management-cluster-landing',
    WORKLOAD_CLUSTER_WIZARD: '/workload-cluster-wizard',
    UNMANAGED_CLUSTER_LANDING: '/unmanaged-cluster-landing',

    // provider workflows
    VSPHERE: 'vsphere', // TODO: refactor to management/workload specific route
    AWS: 'aws',
    DOCKER: 'docker',

    // temp routes to be refactored out
    DEPLOY_PROGRESS: 'progress',
};

export const TAB_NAMES = {
    awsManagementClusterCreateSimple: [
        'AWS Credentials',
        'Cluster settings'
    ],
    dockerManagementClusterCreateSimple: [
        'Prerequisites',
        'Cluster settings'
    ]
};
