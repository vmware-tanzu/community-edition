// React imports
import React, { useContext, useEffect } from 'react';
import { Routes, Route, useLocation } from 'react-router-dom';

// App imports
import { APP_ENV_CHANGE, APP_ROUTE_CHANGE } from './state-management/actions/App.actions';
import { NavRoutes } from './shared/constants/NavRoutes.constants';
import { Store } from './state-management/stores/Store';
import AwsManagementCluster from './views/management-cluster/aws/AwsManagementCluster';
import DeployProgress from './shared/components/DeployProgress/DeployProgress';
import GettingStarted from './views/getting-started/GettingStarted';
import HeaderBar from './shared/components/HeaderBar/HeaderBar';
import ManagementClusterSelectProvider from './views/management-cluster/ManagementClusterSelectProvider';
import SideNavigation from './shared/components/SideNavigation/SideNavigation';
import UnmanagedClusterInventory from './views/unmanaged-cluster/UnmanagedClusterInventory';
import VSphere from './components/VSphere';
import Welcome from './views/welcome/Welcome';
import WorkloadClusterWorkflow from './views/workload-cluster/WorkloadClusterWorkflow';
import DockerManagementCluster from './views/management-cluster/docker/DockerManagementCluster';
import ManagementClusterInventory from './views/management-cluster/ManagementClusterInventory';

function App() {
    const { dispatch } = useContext(Store);
    const location = useLocation();
    const currentPath = location.pathname;

    // TODO: this is for testing/setup of dark mode; sets body theme to dark
    // Will be refactored
    document.body.setAttribute('cds-theme', 'dark');
    document.body.setAttribute('class', 'dark');

    // set router path in store
    useEffect(() => {
        if (currentPath) {
            dispatch({
                type: APP_ROUTE_CHANGE,
                payload: {
                    value: currentPath,
                },
            });
        }
    }, [currentPath]); // eslint-disable-line react-hooks/exhaustive-deps

    // set app environment in store (dev/prod)
    useEffect(() => {
        if (process.env.NODE_ENV) {
            dispatch({
                type: APP_ENV_CHANGE,
                payload: {
                    value: process.env.NODE_ENV,
                },
            });
        }
    }, []); // eslint-disable-line react-hooks/exhaustive-deps

    return (
        <main cds-layout="vertical align:horizontal-stretch" cds-text="body">
            <HeaderBar />
            <section cds-layout="horizontal align:vertical-stretch wrap:none">
                <SideNavigation />
                <div cds-layout="vertical align:stretch">
                    <div cds-layout="grid gap:md gap@md:lg p:lg p@sm:lg p-y@lg:lg container:fill container:left cols:12">
                        <Routes>
                            <Route path={NavRoutes.WELCOME} element={<Welcome />}></Route>
                            <Route path={NavRoutes.GETTING_STARTED} element={<GettingStarted />}></Route>
                            <Route path={NavRoutes.MANAGEMENT_CLUSTER_INVENTORY} element={<ManagementClusterInventory />}></Route>
                            <Route path={NavRoutes.UNMANAGED_CLUSTER_INVENTORY} element={<UnmanagedClusterInventory />}></Route>
                            <Route
                                path={NavRoutes.MANAGEMENT_CLUSTER_SELECT_PROVIDER}
                                element={<ManagementClusterSelectProvider />}
                            ></Route>
                            <Route path={NavRoutes.WORKLOAD_CLUSTER_WIZARD} element={<WorkloadClusterWorkflow />}></Route>
                            <Route path={NavRoutes.VSPHERE} element={<VSphere />}></Route>
                            <Route path={NavRoutes.DOCKER} element={<DockerManagementCluster />}></Route>
                            <Route path={NavRoutes.AWS} element={<AwsManagementCluster />}></Route>
                            <Route path={NavRoutes.DEPLOY_PROGRESS} element={<DeployProgress />}></Route>
                        </Routes>
                    </div>
                </div>
            </section>
        </main>
    );
}

export default App;
