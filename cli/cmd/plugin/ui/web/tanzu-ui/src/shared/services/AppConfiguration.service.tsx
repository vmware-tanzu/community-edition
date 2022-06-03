export enum AppFeature {
    WORKLOAD_CLUSTER_SUPPORT,
    UNMANAGED_CLUSTER_SUPPORT,
}

export function featureAvailable(feature: AppFeature): boolean {
    switch (feature) {
        case AppFeature.UNMANAGED_CLUSTER_SUPPORT:
            return true;
        case AppFeature.WORKLOAD_CLUSTER_SUPPORT:
            return true;
        default:
            return false;
    }
}
