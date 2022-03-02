// React imports
import React from 'react';
import { Routes, Route, Link } from 'react-router-dom';

// Library imports
import { CdsIcon } from "@cds/react/icon";
import { ClarityIcons, homeIcon, compassIcon } from '@cds/core/icon'
import { CdsNavigation, CdsNavigationItem, CdsNavigationStart } from '@cds/react/navigation';

// App imports
import HeaderComponent from './shared/components/header/header.component';
import { handleNavToggle } from './shared/utilities/event-handlers.utility';

ClarityIcons.addIcons(homeIcon, compassIcon);

function App(this: any) {
    // Note: this is for testing/setup of dark mode; sets body theme to dark
    // Will be refactored
    document.body.setAttribute("cds-theme", "dark");
    document.body.setAttribute("class", "dark");

    return (
        <main cds-layout="vertical align:horizontal-stretch" cds-text="body">
            <HeaderComponent/>
            <section cds-layout="horizontal align:vertical-stretch wrap:none">
                <CdsNavigation expanded onExpandedChange={handleNavToggle}>
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
                                <Route path="/" element={<Home />} />
                                <Route path="/about" element={<About />} />
                            </Routes>
                        </div>
                    </div>
                </div>
            </section>
        </main>
    );
}

function Home() {
    return (
        <>
            <h2>Welcome to the Home page</h2>
            <p>Home page content</p>
            <nav>
                <Link to="/about">About</Link>
            </nav>
        </>
    );
}

function About() {
    return (
        <>
            <h2>Who are we?</h2>
            <p>
                We are Tanzu
            </p>
            <nav>
                <Link to="/">Home</Link>
            </nav>
        </>
    );
}

export default App;