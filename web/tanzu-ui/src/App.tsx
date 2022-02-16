// React imports
import * as React from "react";
import { Routes, Route, Link } from "react-router-dom";

// App imports
import HeaderComponent from "./shared/components/header/header.component";
import { CdsDivider } from "@cds/react/divider";
import { CdsButton } from "@cds/react/button";

function App() {
    // Note: this is for testing/setup of dark mode; sets body theme to dark
    // Will be refactored
    document.body.setAttribute("cds-theme", "dark");
    document.body.setAttribute("class", "dark");

    return (
        <main cds-layout="vertical align:stretch" cds-text="body">
            <HeaderComponent/>
            <section cds-layout="horizontal align:vertical-stretch wrap:none">
                <nav cds-text="demo-sidenav" cds-layout="p:md p@md:lg">sidebar
                    <Routes>
                        <Route path="/" element={<Home />} />
                        <Route path="about" element={<About />} />
                    </Routes>
                </nav>
                <CdsDivider cds-text="demo-divider" orientation="vertical"></CdsDivider>
                <div cds-layout="vertical align:stretch">
                    <div cds-text="demo-content demo-scrollable-content">
                        <div cds-layout="vertical gap:md p:lg">
                            <Routes>
                                <Route path="/" element={<Home />} />
                                <Route path="about" element={<About />} />
                            </Routes>
                        </div>
                    </div>
                </div>
            </section>

            {/* Testing CDS Core buttons. TODO: wire up some events */}
            <section cds-layout="horizontal gap:sm">
                <CdsButton status="primary">primary</CdsButton>
                <CdsButton status="success">success</CdsButton>
                <CdsButton status="danger">danger</CdsButton>
                <CdsButton status="danger" disabled>
                    disabled
                </CdsButton>
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