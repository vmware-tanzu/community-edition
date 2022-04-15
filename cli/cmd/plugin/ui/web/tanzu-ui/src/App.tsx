// React imports
import React, { useContext, useEffect } from 'react';
import { Routes, Route } from 'react-router-dom';

// App imports
import HeaderBar from './shared/components/HeaderBar/HeaderBar';
import SideNavigation from './shared/components/SideNavigation/SideNavigation';
import GettingStarted from './views/getting-started/GettingStarted';
import VSphere from './components/VSphere';
import Welcome from './views/welcome/Welcome';
import DeployProgress from './views/temp/DeployProgress';
import { APP_ENV_CHANGE, AppActionNames } from './state-management/actions/App.actions';
import { Store } from './state-management/stores/Store';
import { NavRoutes } from './shared/constants/NavRoutes.constants';
import ManagementClusterLanding from './views/management-cluster/ManagementClusterLanding';
import WorkloadClusterLanding from './views/workload-cluster/WorkloadClusterLanding';

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
                    <div cds-text="demo-content demo-scrollable-content">
                        <div cds-layout="vertical gap:md p:lg">
                            <Routes>
                                <Route path={NavRoutes.WELCOME} element={<Welcome />}></Route>
                                <Route path={NavRoutes.GETTING_STARTED} element={<GettingStarted />}></Route>
                                <Route path={NavRoutes.MANAGEMENT_CLUSTER_LANDING} element={<ManagementClusterLanding />}></Route>
                                <Route path={NavRoutes.WORKLOAD_CLUSTER_LANDING} element={<WorkloadClusterLanding />}></Route>
                                <Route path={NavRoutes.VSPHERE} element={<VSphere />}></Route>
                                <Route path={NavRoutes.DEPLOY_PROGRESS} element={<DeployProgress />}></Route>
                            </Routes>
                        </div>
                    </div>
                </div>
            </section>
        </main>
    );
}

export default App;
