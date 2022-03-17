// React imports
import React from 'react';
import { Routes, Route } from 'react-router-dom';

// App imports
import HeaderComponent from './shared/components/header/header.component';
import SideNavigationComponent from './shared/components/side-navigation/side-navigation.component';
import GetStartedComponent from './views/getting-started/getting-started.component';
import VSphere from './components/VSphere';
import WelcomeComponent from './views/welcome/welcome.component';

function App(this: any) {
    // Note: this is for testing/setup of dark mode; sets body theme to dark
    // Will be refactored
    document.body.setAttribute('cds-theme', 'dark');
    document.body.setAttribute('class', 'dark');

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
                            </Routes>
                        </div>
                    </div>
                </div>
            </section>
        </main>
    );
}

export default App;