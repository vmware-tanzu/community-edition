// React imports
import React, { useContext } from 'react';
import { Routes, Route, Link } from 'react-router-dom';

// Library imports
import { CdsIcon } from '@cds/react/icon';
import { ClarityIcons, homeIcon, compassIcon } from '@cds/core/icon';
import { CdsNavigation, CdsNavigationItem, CdsNavigationStart } from '@cds/react/navigation';

// App imports
import HeaderComponent from './shared/components/header/header.component';
import LandingPage from './components/LandingPage';
import VSphere from './components/VSphere';
import { Store } from './stores/store';
import { TOGGLE_NAV } from './constants/actionTypes';

ClarityIcons.addIcons(homeIcon, compassIcon);

function App(this: any) {
    // Note: this is for testing/setup of dark mode; sets body theme to dark
    // Will be refactored
    document.body.setAttribute('cds-theme', 'dark');
    document.body.setAttribute('class', 'dark');

    const { state, dispatch } = useContext(Store);

    const onNavExpandedChange = () => {
        dispatch({
            type: TOGGLE_NAV
        });
    };

    return (
        <main cds-layout="vertical align:horizontal-stretch" cds-text="body">
            <HeaderComponent/>
            <section cds-layout="horizontal align:vertical-stretch wrap:none">
                <CdsNavigation expanded={state.ui.navExpanded} onExpandedChange={onNavExpandedChange}>
                    <CdsNavigationStart></CdsNavigationStart>
                    <CdsNavigationItem>
                        <Link to="/">
                            <CdsIcon shape="home" size="sm"></CdsIcon>
                            Home
                        </Link>
                    </CdsNavigationItem>
                    <CdsNavigationItem>
                        <Link to="/about">
                            <CdsIcon shape="compass" size="sm"></CdsIcon>
                            Welcome
                        </Link>
                    </CdsNavigationItem>
                </CdsNavigation>
                <div cds-layout="vertical align:stretch">
                    <div cds-text="demo-content demo-scrollable-content">
                        <div cds-layout="vertical gap:md p:lg">
                            <Routes>
                                <Route path="/" element={<LandingPage />}></Route>
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