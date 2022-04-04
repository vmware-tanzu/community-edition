import * as React from 'react';

import '../../assets/css/Concierge.css'

function Concierge() {
  const [existingTanzu] = React.useState({})
  return (
    <div>
      <h1>Welcome to the Tanzu Concierge!</h1>

      <p>We're here to help you install the Tanzu binary</p>
        <div id="existingTanzuInfo"></div>
      <button type="button" id="buttonPreInstall">
        Re-do the pre-install check
      </button>
        <p>&nbsp;</p>
        <button type="button" id="buttonInstall">
            Install Tanzu, baby!
        </button>
        <div id="stepName"></div><div id="percentComplete"></div>

        <p>&nbsp;</p>
        <div id="installProgressDisplay"></div>
    </div>
  )
}

export default Concierge
