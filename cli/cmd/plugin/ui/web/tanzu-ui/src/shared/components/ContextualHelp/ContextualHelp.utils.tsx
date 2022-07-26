import { ContextualHelpState } from './ContextualHelp.store';
import { NavRoutes } from '../../constants/NavRoutes.constants';

const enum CtxTitle {
    TanzuCommunityEdition = 'Tanzu Community Edition',
    TanzuClusters = 'Tanzu Clusters',
    GettingStarted = 'Getting Started',
    ManagementClusters = 'Management Clusters',
    WorkloadClusters = 'Workload Clusters',
    UnmanagedClusters = 'Unmanaged Clusters',
    CreateManagementCluster = 'Create Management Cluster',
    CreateWorkloadCluster = 'Create Workload Cluster',
    CreateUnmanagedCluster = 'Create Unmanaged Cluster',
    DeployProgress = 'Deploy Progress',
}

const enum CtxKeys {
    Welcome = 'ctx-welcome',
    GettingStarted = 'ctx-getting-started',
    ManagementClusters = 'ctx-management-clusters',
    WorkloadClusters = 'ctx-workload-clusters',
    UnmanagedClusters = 'ctx-unmanaged-clusters',
}

const welcomePage: ContextualHelpState = {
    externalLink: 'https://tanzucommunityedition.io/resources',
    keywords: [CtxKeys.Welcome],
    title: {
        contextTitle: CtxTitle.TanzuCommunityEdition,
        pageTitle: CtxTitle.TanzuCommunityEdition,
    },
};

const gettingStartedPage: ContextualHelpState = {
    externalLink: 'https://tanzucommunityedition.io/docs/',
    keywords: [CtxKeys.GettingStarted],
    title: {
        contextTitle: CtxTitle.GettingStarted,
        pageTitle: CtxTitle.TanzuClusters,
    },
};

const managementClustersPage: ContextualHelpState = {
    externalLink: 'http://tanzucommunityedition.io',
    keywords: [CtxKeys.ManagementClusters],
    title: {
        contextTitle: CtxTitle.ManagementClusters,
        pageTitle: CtxTitle.ManagementClusters,
    },
};

const workloadClustersPage: ContextualHelpState = {
    externalLink: 'http://tanzucommunityedition.io',
    keywords: [CtxKeys.WorkloadClusters],
    title: {
        contextTitle: CtxTitle.WorkloadClusters,
        pageTitle: CtxTitle.WorkloadClusters,
    },
};

const unmanagedClustersPage: ContextualHelpState = {
    externalLink: 'http://tanzucommunityedition.io',
    keywords: [CtxKeys.UnmanagedClusters],
    title: {
        contextTitle: CtxTitle.UnmanagedClusters,
        pageTitle: CtxTitle.UnmanagedClusters,
    },
};

const createManagementCluster: ContextualHelpState = {
    externalLink: 'http://tanzucommunityedition.io',
    keywords: [CtxKeys.ManagementClusters],
    title: {
        contextTitle: CtxTitle.CreateManagementCluster,
        pageTitle: CtxTitle.CreateManagementCluster,
    },
};

const createWorkloadCluster: ContextualHelpState = {
    externalLink: 'http://tanzucommunityedition.io',
    keywords: [CtxKeys.WorkloadClusters],
    title: {
        contextTitle: CtxTitle.CreateWorkloadCluster,
        pageTitle: CtxTitle.CreateWorkloadCluster,
    },
};

const createUnmanagedCluster: ContextualHelpState = {
    externalLink: 'http://tanzucommunityedition.io',
    keywords: [CtxKeys.UnmanagedClusters],
    title: {
        contextTitle: CtxTitle.CreateUnmanagedCluster,
        pageTitle: CtxTitle.CreateUnmanagedCluster,
    },
};

const defaultContextualHelp: ContextualHelpState = {
    title: {
        contextTitle: CtxTitle.TanzuCommunityEdition,
        pageTitle: CtxTitle.TanzuCommunityEdition,
    },
    keywords: [CtxKeys.Welcome, CtxKeys.GettingStarted],
    externalLink: 'http://tanzucommunityedition.io',
};

export const determineContextualHelpContent = (pathname: NavRoutes): ContextualHelpState => {
    const mapper: Partial<Record<NavRoutes, ContextualHelpState>> = {
        [NavRoutes.WELCOME]: welcomePage,
        [NavRoutes.GETTING_STARTED]: gettingStartedPage,

        [NavRoutes.MANAGEMENT_CLUSTER_INVENTORY]: managementClustersPage,
        [NavRoutes.WORKLOAD_CLUSTER_INVENTORY]: workloadClustersPage,
        [NavRoutes.UNMANAGED_CLUSTER_INVENTORY]: unmanagedClustersPage,

        [NavRoutes.MANAGEMENT_CLUSTER_SELECT_PROVIDER]: createManagementCluster,
        [NavRoutes.WORKLOAD_CLUSTER_LANDING]: createWorkloadCluster,
        [NavRoutes.UNMANAGED_CLUSTER_WIZARD]: createUnmanagedCluster,
    };
    return mapper[pathname] ?? defaultContextualHelp;
};
