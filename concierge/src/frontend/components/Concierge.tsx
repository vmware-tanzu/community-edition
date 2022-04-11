import * as React from 'react';
import { Routes, Route } from 'react-router-dom';

import '../../assets/css/Concierge.css'
import Welcome from './Welcome'
import Install from './Install'
import { useState } from 'react';
import { AvailableInstallation, PreInstallation } from '_models/installation';
import { ProgressMessage } from '_models/progressMessage';

function findSelectedVersion(version: string, installations: AvailableInstallation[]) : AvailableInstallation {
    return installations.find(installation => installation.version === version)
}

// Expects:
// props.runTanzu                               a (callback) function that allows this component to kick off running the Tanzu UI
// props.install                                a (callback) function that allows this component to kick off installing Tanzu
// props.refreshInstallations                   a (callback) function that allows this component to kick off checking for available
//                                              installations
// props.registerSetPreinstallationFunction     a (callback) function that allows this component to register a setPreinstallation function
//                                              which our parent can call to set the preinstallation data into our state. (Used when our
//                                              parent receives a PreInstallation data object from the main process.)
// props.registerProgressReceiver               a (callback) function that allows this component to register a progressReceiver function
//                                              which our parent can call to give progress updates during installation. (Used when our
//                                              parent receives a ProgressMessage data object from the main process.)
function Concierge(props) {
    const [preinstallationData, setPreinstallationData] = useState<PreInstallation>()
    const [selectedInstallationVersion, setInstallationVersion] = useState<string>()
    const [selectedInstallation, setSelectedInstallation] = useState<AvailableInstallation>()
    const [progressMessages, setProgressMessages] = useState<ProgressMessage[]>([])

    validateProps(props)

    const onSelectInstallationVersion = (version) => {
        console.log(`CONCIERGE: user selected version ${version}`)
        setInstallationVersion(version)
        const selectedVersion = findSelectedVersion(version, preinstallationData.availableInstallations)
        setSelectedInstallation(selectedVersion)
        if (!selectedVersion) {
            console.log(`PROGRAMMER ERROR: (CONCIERGE) user selected version ${version}, but unable to find corresponding object`)
        }
    }

    // We register a method that allows our enclosing component to send us preinstallation data
    props.registerSetPreinstallationFunction((preinstallationData) => {
        console.log(`CONCIERGE: just received preinstallation data ${JSON.stringify(preinstallationData)}`)
        setPreinstallationData(preinstallationData)
        // by default, select the first version/installation in the list
        if (preinstallationData && preinstallationData.availableInstallations && preinstallationData.availableInstallations.length > 0) {
            const defaultInstallation = preinstallationData.availableInstallations[0]
            setSelectedInstallation(defaultInstallation)
            setInstallationVersion(defaultInstallation.version)
        }
    })

    // We register a method that allows our enclosing component to send us new progress messages during the installation
    props.registerProgressReceiver((pm: ProgressMessage) => {
        console.log(`RECEIVED PROGRESS MESSAGE: ${JSON.stringify(pm)}`)
        setProgressMessages([...progressMessages, pm])
    })

  return (
    <main cds-layout="vertical align:horizontal-stretch" cds-text="body">
        <section cds-layout="horizontal align:vertical-stretch wrap:none">
                <div cds-layout="vertical align:stretch">
                    <div cds-text="demo-content demo-scrollable-content">
                        <div cds-layout="vertical gap:md p:lg">
                            <Routes>
                                <Route path="/" element={<Welcome
                                    install={props.install}
                                    preinstallationData={preinstallationData}
                                    refreshInstallations={props.refreshInstallations}
                                    runTanzu={props.runTanzu}
                                    selectedInstallationVersion={selectedInstallationVersion}
                                    selectedInstallation={selectedInstallation}
                                    setInstallationVersion={onSelectInstallationVersion}
                                />}></Route>
                                <Route path="/install" element={<Install
                                    progressMessages={progressMessages}
                                />}></Route>
                            </Routes>
                        </div>
                    </div>
                </div>
        </section>
    </main>
  )
}

function validateProps(props) {
    if (!props.registerSetPreinstallationFunction) {
        console.log('PROGRAMMER ERROR: no props.registerSetPreinstallationFunction sent to the Concierge')
    }
    if (!props.registerProgressReceiver) {
        console.log('PROGRAMMER ERROR: no props.registerProgressReceiver sent to the Concierge')
    }
}

export default Concierge
