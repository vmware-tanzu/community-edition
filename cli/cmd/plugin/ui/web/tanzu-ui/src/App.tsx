// React imports
import React, {useContext, useEffect} from 'react';
import { Routes, Route } from 'react-router-dom';

// App imports
import HeaderComponent from './shared/components/header/header.component';
import SideNavigationComponent from './shared/components/side-navigation/SideNavigation.component';
import GetStartedComponent from './views/getting-started/getting-started.component';
import VSphere from './components/VSphere';
import WelcomeComponent from './views/welcome/welcome.component';
import ProgressComponent from './views/temp/progress.component';
import { APP_ENV_CHANGE } from './state-management/actions/app.actions';
import { Store } from './state-management/stores/store';

function App(this: any) {
    const { dispatch } = useContext(Store);

    // TODO: this is for testing/setup of dark mode; sets body theme to dark
    // Will be refactored
    document.body.setAttribute('cds-theme', 'dark');
    document.body.setAttribute('class', 'dark');

    useEffect(() => {
        if (process.env.NODE_ENV) {
            // TODO: why can't I seem to apply a type to the generic for 'dispatch' <Action>
            dispatch({
                type: APP_ENV_CHANGE,
                payload: {
                    name: 'appEnv',
                    value: process.env.NODE_ENV
                }
            });
        }
    }, [],);

    return (
        <main cds-layout="vertical align:horizontal-stretch" cds-text="body">
            <HeaderComponent/>
            <section cds-layout="horizontal align:vertical-stretch wrap:none">
                <SideNavigationComponent/>
                <div cds-layout="vertical align:stretch">
                    <div cds-text="demo-content demo-scrollable-content">
                        <div cds-layout="vertical gap:md p:lg">
                            <Routes>
                                <Route path="/" element={<WelcomeComponent />}></Route>
                                <Route path="/getting-started" element={<GetStartedComponent />}></Route>
                                <Route path="/vsphere" element={<VSphere />}></Route>
                                <Route path="/progress" element={<ProgressComponent />}></Route>
                            </Routes>
                        </div>
                    </div>
                </div>
            </section>
        </main>
    );
}

export default App;