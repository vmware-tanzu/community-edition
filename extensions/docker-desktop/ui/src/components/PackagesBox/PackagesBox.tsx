import { Grid } from '@mui/material';
import React from 'react';
import { Actions, ClusterStates, IClusterStatus } from 'providers/globalAppState/reducer';
import { GlobalAppState } from 'providers/globalAppState/GlobalAppStateProvider';
import PackageCard from './PackageCard';
import { icons } from "shared/images";
import { TCE_CERTMANAGER_README, TCE_CONTOUR_README, TCE_KUBEAPPS_README } from 'shared/constants';


export default function PackagesBox() {
  const { state: appState, dispatch } = React.useContext(GlobalAppState);

  function provisionIngress() {
    dispatch({ type: "CLUSTER_STATUS", payload: { status: ClusterStates.RUNNING, description: "Cluster is being provisioned" } } as Actions);
    window.ddClient.extension.vm.cli.exec(`/backend/clustermgr`, [`provision-ingress`])
      .then((cmdResult: any) => {
        let status: IClusterStatus = cmdResult.parseJsonObject();
        dispatch({ type: "PROVISION_INGRESS", payload: true } as Actions);
        // TODO: Really we should make a real status provisioning checker
        dispatch({ type: "CLUSTER_STATUS", payload: status } as Actions);
      });
  }

  return (
  <Grid container spacing={2}>
    <Grid item xs={12}>
      <PackageCard
        title='Ingress'
        description='Use Contour as Ingress controller. This will expose port 80 and 443 on your local laptop (127.0.0.1) and have Ingress controller ready to be used'
        icon={icons.contour}
        createAction={provisionIngress}
        showCreate={(appState.isIngressProvisioned == false)}
        helpLink={TCE_CONTOUR_README}
      />
      <PackageCard
        title='Cert Manager'
        description='Use Cert-manager bla bla bla'
        icon={icons.certmanager}
        createAction={provisionIngress}
        showCreate={false}
        deleteAction={provisionIngress}
        showDelete={false}
        openAction={provisionIngress}
        showOpen={false}
        helpLink={TCE_CERTMANAGER_README}
      />
      <PackageCard
        title='Kubeapps'
        description='Install your favourite application/software using Kubeapps.'
        icon={icons.kubeapps}
        showCreate={true}
        helpLink={TCE_KUBEAPPS_README}
      />
    </Grid>
  </Grid>);
}