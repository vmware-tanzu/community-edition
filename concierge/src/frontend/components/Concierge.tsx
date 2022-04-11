import * as React from 'react';
import { Routes, Route } from 'react-router-dom';

import '../../assets/css/Concierge.css'
import Welcome from './Welcome'
import Install from './Install'

function Concierge(props) {
  return (
    <main cds-layout="vertical align:horizontal-stretch" cds-text="body">
        <section cds-layout="horizontal align:vertical-stretch wrap:none">
                <div cds-layout="vertical align:stretch">
                    <div cds-text="demo-content demo-scrollable-content">
                        <div cds-layout="vertical gap:md p:lg">
                            <Routes>
                                <Route path="/" element={<Welcome refreshInstallations={props.refreshInstallations} />}></Route>
                                <Route path="/install" element={<Install />}></Route>
                            </Routes>
                        </div>
                    </div>
                </div>
        </section>
    </main>
  )
}

export default Concierge
