export interface DependencyMap {
    [key: string]: string[];
}
export const dependencyMap: DependencyMap = {
    VCENTER_PASSWORD: ['CLUSTER_NAME'],
};
