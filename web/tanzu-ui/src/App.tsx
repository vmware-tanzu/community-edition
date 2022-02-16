// React imports
import * as React from "react";
import { Routes, Route, Link } from "react-router-dom";

// Component imports
import HeaderComponent from "./shared/components/header/header.component";
import {CdsAlert, CdsAlertGroup} from '@cds/react/alert';
import {CdsCard} from "@cds/react/card";
import {CdsDivider} from "@cds/react/divider";
import {CdsButton} from "@cds/react/button";

function App() {
    // sets body theme to dark
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
                That feels like an existential question, don't you
                think?
            </p>
            <nav>
                <Link to="/">Home</Link>
            </nav>
        </>
    );
}

export default App;