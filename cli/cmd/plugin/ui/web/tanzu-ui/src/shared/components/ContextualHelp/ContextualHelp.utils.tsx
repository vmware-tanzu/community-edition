import { NavRoutes } from '../../constants/NavRoutes.constants';
import { ContextualHelpState } from './ContextualHelp.store';

const enum ContextTitle {
    TanzuCommunityEdition = 'Tanzu Community Edition',
    GettingStarted = 'Getting Started',
    ManagementClusters = 'Management Clusters',
    WorkloadClusters = 'Workload Clusters',
    UnmanagedClusters = 'Unmanaged Clusters',
    CreateManagementCluster = 'Create Management Cluster',
    CreateWorkloadCluster = 'Create Workload Cluster',
    CreateUnmanagedCluster = 'Create Unmanaged Cluster',
    DeployProgress = 'Deploy Progress',
}

const welcomePage: ContextualHelpState = {
    externalLink: 'http://tanzucommunityedition.io',
    keywords: ['tce-welcome'],
    title: {
        contextTitle: ContextTitle.TanzuCommunityEdition,
        pageTitle: 'Welcome to Tanzu',
    },
};

const gettingStartedPage: ContextualHelpState = {
    externalLink: 'http://tanzucommunityedition.io',
    keywords: ['tce-getting-started'],
    title: {
        contextTitle: ContextTitle.GettingStarted,
        pageTitle: ContextTitle.GettingStarted,
    },
};

const contextualHelpForManagementClustersPage: ContextualHelpState = {
    externalLink: 'http://tanzucommunityedition.io',
    keywords: ['tce-management-clusters'],
    title: {
        contextTitle: ContextTitle.ManagementClusters,
        pageTitle: ContextTitle.ManagementClusters,
    },
};

const contextualHelpForWorkloadClustersPage: ContextualHelpState = {
    externalLink: 'http://tanzucommunityedition.io',
    keywords: ['tce-workload-clusters'],
    title: {
        contextTitle: ContextTitle.WorkloadClusters,
        pageTitle: ContextTitle.WorkloadClusters,
    },
};

const contextualHelpForUnmanagedClustersPage: ContextualHelpState = {
    externalLink: 'http://tanzucommunityedition.io',
    keywords: ['tce-unmanaged-clusters'],
    title: {
        contextTitle: ContextTitle.UnmanagedClusters,
        pageTitle: ContextTitle.UnmanagedClusters,
    },
};

const contextualHelpForCreateManagementCluster: ContextualHelpState = {
    externalLink: 'http://tanzucommunityedition.io',
    keywords: ['tce-create-management-cluster'],
    title: {
        contextTitle: ContextTitle.CreateManagementCluster,
        pageTitle: ContextTitle.CreateManagementCluster,
    },
};

const contextualHelpForCreateWorkloadCluster: ContextualHelpState = {
    externalLink: 'http://tanzucommunityedition.io',
    keywords: ['tce-docker'],
    title: {
        contextTitle: ContextTitle.CreateWorkloadCluster,
        pageTitle: ContextTitle.CreateWorkloadCluster,
    },
};

const contextualHelpForCreateUnmanagedCluster: ContextualHelpState = {
    externalLink: 'http://tanzucommunityedition.io',
    keywords: ['tce-docker'],
    title: {
        contextTitle: ContextTitle.CreateUnmanagedCluster,
        pageTitle: ContextTitle.CreateUnmanagedCluster,
    },
};

// const contextualHelpForVSphere: ContextualHelpState = {
//     externalLink: 'http://tanzucommunityedition.io',
//     keywords: ['tce-docker'],
//     title: {
//         contextTitle: 'Unmanaged Clusters',
//         pageTitle: 'Getting Started',
//     },
// };
// const contextualHelpForAWS: ContextualHelpState = {
//     externalLink: 'http://tanzucommunityedition.io',
//     keywords: ['tce-docker'],
//     title: {
//         contextTitle: 'Unmanaged Clusters',
//         pageTitle: 'Getting Started',
//     },
// };

// const contextualHelpForDocker: ContextualHelpState = {
//     externalLink: 'http://tanzucommunityedition.io',
//     keywords: ['tce-docker'],
//     title: {
//         contextTitle: 'Unmanaged Clusters',
//         pageTitle: 'Getting Started',
//     },
// };

// const contextualHelpForDeployProgress: ContextualHelpState = {
//     externalLink: 'http://tanzucommunityedition.io',
//     keywords: ['tce-docker'],
//     title: {
//         contextTitle: ContextTitle.DeployProgress,
//         pageTitle: ContextTitle.DeployProgress,
//     },
// };
const defaultContextualHelp: ContextualHelpState = {
    title: {
        contextTitle: ContextTitle.TanzuCommunityEdition,
        pageTitle: 'Welcome to Tanzu',
    },
    keywords: ['tce-welcome', 'default'],
    externalLink: 'http://tanzucommunityedition.io',
};

export const determineContextualHelpContent = (pathname: NavRoutes): ContextualHelpState => {
    const mapper: Partial<Record<NavRoutes, ContextualHelpState>> = {
        [NavRoutes.WELCOME]: welcomePage,
        [NavRoutes.GETTING_STARTED]: gettingStartedPage,

        [NavRoutes.MANAGEMENT_CLUSTER_INVENTORY]: contextualHelpForManagementClustersPage,
        [NavRoutes.WORKLOAD_CLUSTER_INVENTORY]: contextualHelpForWorkloadClustersPage,
        [NavRoutes.UNMANAGED_CLUSTER_INVENTORY]: contextualHelpForUnmanagedClustersPage,

        [NavRoutes.MANAGEMENT_CLUSTER_SELECT_PROVIDER]: contextualHelpForCreateManagementCluster,
        [NavRoutes.WORKLOAD_CLUSTER_LANDING]: contextualHelpForCreateWorkloadCluster,
        [NavRoutes.UNMANAGED_CLUSTER_WIZARD]: contextualHelpForCreateUnmanagedCluster,
        // [NavRoutes.WORKLOAD_CLUSTER_WIZARD]: contextualHelpForCreateWorkloadClusterWizard,

        // [NavRoutes.VSPHERE]: contextualHelpForVSphere,
        // [NavRoutes.AWS]: contextualHelpForAWS,
        // [NavRoutes.DOCKER]: contextualHelpForDocker,
        // [NavRoutes.DEPLOY_PROGRESS]: contextualHelpForDeployProgress,
    };
    return mapper[pathname] ?? defaultContextualHelp;
};
