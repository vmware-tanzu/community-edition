// React imports
import React, { useContext, useEffect } from 'react';
import { Routes, Route } from 'react-router-dom';

// App imports
import { APP_ENV_CHANGE, AppActionNames } from './state-management/actions/App.actions';
import { NavRoutes } from './shared/constants/NavRoutes.constants';
import { Store } from './state-management/stores/Store';
import AwsManagementClusterSimple from './views/management-cluster/AwsManagementClusterSimple';
import DeployProgress from './shared/components/DeployProgress/DeployProgress';
import GettingStarted from './views/getting-started/GettingStarted';
import HeaderBar from './shared/components/HeaderBar/HeaderBar';
import ManagementClusterLanding from './views/management-cluster/ManagementClusterLanding';
import SideNavigation from './shared/components/SideNavigation/SideNavigation';
import UnmanagedClusterLanding from './views/unmanaged-cluster/UnmanagedClusterLanding';
import VSphere from './components/VSphere';
import Welcome from './views/welcome/Welcome';
import WorkloadClusterWorkflow from './views/workload-cluster/WorkloadClusterWorkflow';
import DockerManagementClusterSimple from './views/management-cluster/DockerManagementClusterSimple';

function App() {
    const { dispatch } = useContext(Store);

    // TODO: this is for testing/setup of dark mode; sets body theme to dark
    // Will be refactored
    document.body.setAttribute('cds-theme', 'dark');
    document.body.setAttribute('class', 'dark');

    useEffect(() => {
        if (process.env.NODE_ENV) {
            dispatch({
                type: APP_ENV_CHANGE,
                payload: {
                    name: AppActionNames.appEnv,
                    value: process.env.NODE_ENV
                }
            });
        }
    }, []); // eslint-disable-line react-hooks/exhaustive-deps

    return (
        <main cds-layout="vertical align:horizontal-stretch" cds-text="body">
            <HeaderBar/>
            <section cds-layout="horizontal align:vertical-stretch wrap:none">
                <SideNavigation/>
                <div cds-layout="vertical align:stretch">
                    <div cds-layout="grid gap:md gap@md:lg p:lg p@sm:lg p-y@lg:lg container:fill container:left cols:12">
                        <Routes>
                            <Route path={NavRoutes.WELCOME} element={<Welcome />}></Route>
                            <Route path={NavRoutes.GETTING_STARTED} element={<GettingStarted />}></Route>
                            <Route path={NavRoutes.MANAGEMENT_CLUSTER_LANDING} element={<ManagementClusterLanding />}></Route>
                            <Route path={NavRoutes.WORKLOAD_CLUSTER_WIZARD} element={<WorkloadClusterWorkflow />}></Route>
                            <Route path={NavRoutes.UNMANAGED_CLUSTER_LANDING} element={<UnmanagedClusterLanding />}></Route>
                            <Route path={NavRoutes.VSPHERE} element={<VSphere />}></Route>
                            <Route path={NavRoutes.DOCKER} element={<DockerManagementClusterSimple />}></Route>
                            <Route path={NavRoutes.AWS} element={<AwsManagementClusterSimple />}></Route>
                            <Route path={NavRoutes.DEPLOY_PROGRESS} element={<DeployProgress />}></Route>
                        </Routes>
                    </div>
                </div>
            </section>
        </main>
    );
}

export default App;
