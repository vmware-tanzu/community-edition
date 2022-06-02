export const NavRoutes = {
    // app general
    WELCOME: '/',
    GETTING_STARTED: '/getting-started',

    // cluster inventory landing pages
    MANAGEMENT_CLUSTER_INVENTORY: '/management-clusters',
    WORKLOAD_CLUSTER_INVENTORY: '/workload-clusters', // TODO: refactor w/shimon
    UNMANAGED_CLUSTER_INVENTORY: '/unmanaged-clusters',

    // cluster create pages
    MANAGEMENT_CLUSTER_SELECT_PROVIDER: '/management-cluster-provider',
    WORKLOAD_CLUSTER_LANDING: '/workload-cluster-landing',
    UMANAGED_CLUSTER_LANDING: '/unmanaged-cluster-landing',
    UMANAGED_CLUSTER_WIZARD: '/unmanaged-cluster-wizard',
    WORKLOAD_CLUSTER_WIZARD: '/workload-cluster-wizard',

    // provider workflows
    VSPHERE: 'vsphere', // TODO: refactor to management/workload specific route
    AWS: 'aws',
    DOCKER: 'docker',

    // temp routes to be refactored out
    DEPLOY_PROGRESS: 'progress',
};

export const AWS_MC_BASIC_TAB_NAMES = ['AWS Credentials', 'Cluster settings'];

export const DOCKER_MC_BASIC_TAB_NAMES = ['Prerequisites', 'Cluster settings'];

export const UMC_BASIC_TAB_NAMES = ['Cluster settings', 'Optional settings'];
