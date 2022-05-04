import CircleRoundedIcon from '@mui/icons-material/CircleRounded';
import {
  Box,
  Button,
  Container,
  Stack,
  Tab,
  Tabs,
  Typography,
} from '@mui/material';

import React, { useEffect } from 'react';
import {
  KubeconfigBox,
  FeedbackLink,
  LogTailerBox,
  PackagesBox,
} from 'components';
import { GlobalAppState } from 'providers';
import {
  Actions,
  ClusterStates,
  emptyClusterStats,
  IAppState,
  IClusterResourceStats,
  IClusterStatus,
} from 'providers/globalAppState/reducer';
import StatsPane from 'components/StatsPane/StatsPane';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <Box
      sx={{ height: '100%' }}
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3, height: '100%' }}>{children}</Box>}
    </Box>
  );
}

function a11yProps(index: number) {
  return {
    id: `simple-tab-${index}`,
    'aria-controls': `simple-tabpanel-${index}`,
  };
}

export interface AppProps {
  createCluster?: boolean;
}

export default function App(props: AppProps) {
  const TABS = {
    LOGS: 0,
    KUBECONFIG: 1,
    PACKAGES: 2,
    STATS: 3,
  };

  // TODO: Call to create a cluster if App is invoked from Intro page
  const { state: appState, dispatch } = React.useContext(GlobalAppState);

  // First load of the page. We refresh the cluster status from backend to verify it matches localStorage
  // TODO: See if we can move this to the intro screen or somewhere that can execute faster
  // TODO: Also, have a periodic sync of state
  useEffect(() => {
    if (props.createCluster) createCluster();
    window.ddClient.extension.vm.cli
      .exec(`/backend/clustermgr`, [`status`])
      .then((cmdResult: any) => {
        let status: IClusterStatus = cmdResult.parseJsonObject();
        if (
          status.status !== appState.clusterStatus.status ||
          status.description !== appState.clusterStatus.description
        ) {
          // Depends on what's the new state, do some other actions
          if (status.status === ClusterStates.DELETED) {
            dispatch({ type: 'CLUSTER_LOGS', payload: [] } as Actions);
            dispatch({ type: 'FETCH_KUBECONFIG', payload: '' } as Actions);
            dispatch({ type: 'CLUSTER_STARTED', payload: false } as Actions);
            dispatch({ type: 'PROVISION_INGRESS', payload: false } as Actions);
            dispatch({ type: 'ERROR', payload: '' } as Actions);
          }
          dispatch({ type: 'CLUSTER_STATUS', payload: status } as Actions);
        }
      });
  }, []);

  // Update other states based on cluster status
  useEffect(() => {
    let interval: any = null;
    if (appState.clusterStatus.status === ClusterStates.RUNNING) {
      // TODO: Move this into a dispatch for when Created and deleted
      window.ddClient.desktopUI.toast.success(
        appState.clusterStatus.description,
      );
      dispatch({ type: 'CLUSTER_STARTED', payload: true } as Actions);
      if (interval != undefined) clearInterval(interval);
    } else {
      if (
        appState.clusterStatus.status === ClusterStates.CREATING ||
        appState.clusterStatus.status === ClusterStates.INITIALIZING ||
        appState.clusterStatus.status === ClusterStates.DELETING
      ) {
        dispatch({ type: 'CLUSTER_STARTED', payload: true } as Actions);
        interval = setInterval(() => {
          executeWhileCreating();
        }, 2000);
        return () => clearInterval(interval);
      } else {
        if (
          appState.clusterStatus.status === ClusterStates.DELETED ||
          appState.clusterStatus.status === ClusterStates.NOT_EXISTS
        ) {
          dispatch({ type: 'CLUSTER_STARTED', payload: false } as Actions);
        }
        if (interval != undefined) clearInterval(interval);
      }
    }
  }, [appState.clusterStatus.status]);

  useEffect(() => {
    const textarea = document.getElementById('log');
    if (textarea != null) textarea.scrollTop = textarea.scrollHeight;
  }, [appState.logs]);

  function executeWhileCreating() {
    getClusterStatus();
    getLogs();
  }

  async function getClusterStatus() {
    await window.ddClient.extension.vm.cli
      .exec(`/backend/clustermgr`, [`status`])
      .then((cmdResult: any) => {
        let status: IClusterStatus = cmdResult.parseJsonObject();
        dispatch({ type: 'CLUSTER_STATUS', payload: status } as Actions);
        // TODO: Update other CLUSTER_STARTED based on states
      });
  }

  function createCluster() {
    dispatch({ type: 'CLUSTER_LOGS', payload: [] } as Actions);
    dispatch({ type: 'FETCH_KUBECONFIG', payload: '' } as Actions);
    dispatch({ type: 'CLUSTER_STARTED', payload: true } as Actions);
    dispatch({
      type: 'CLUSTER_STATUS',
      payload: {
        status: ClusterStates.CREATING,
        description: 'Cluster is being created',
      },
    } as Actions);
    window.ddClient.extension.vm.cli
      .exec(`/backend/clustermgr`, [`create`]) // cluster-create
      .then((cmdResult: any) => {
        let status: IClusterStatus = cmdResult.parseJsonObject();
      });
    setSelectedTab(TABS.LOGS);
  }

  function deleteCluster() {
    dispatch({ type: 'CLUSTER_LOGS', payload: [] } as Actions);
    dispatch({ type: 'FETCH_KUBECONFIG', payload: '' } as Actions);
    dispatch({ type: 'CLUSTER_STARTED', payload: false } as Actions);
    dispatch({ type: 'PROVISION_INGRESS', payload: false } as Actions);
    dispatch({ type: 'CLUSTER_STATS', payload: emptyClusterStats } as Actions);
    dispatch({
      type: 'CLUSTER_STATUS',
      payload: {
        status: ClusterStates.DELETING,
        description: 'Cluster is being deleted',
      },
    } as Actions);
    window.ddClient.extension.vm.cli
      .exec(`/backend/clustermgr`, [`delete`]) // /backend/cluster-reset
      .then((cmdResult: any) => {
        let status: IClusterStatus = cmdResult.parseJsonObject();
      });
    setSelectedTab(TABS.LOGS);
  }

  function getClusterKubeconfig() {
    if (appState.clusterStatus.status === ClusterStates.RUNNING) {
      window.ddClient.extension.vm.cli
        .exec(`/backend/clustermgr`, [`kubeconfig`])
        .then((cmdResult: any) => {
          let status: IClusterStatus = cmdResult.parseJsonObject();
          // let kubeconfig: string = cmdResult.lines().map(
          //   (line: any) =>
          //     `${line}`
          // )
          //   .join("\n");
          dispatch({
            type: 'FETCH_KUBECONFIG',
            payload: status.output,
          } as Actions);
        });
      // setSelectedTab(TABS.KUBECONFIG);
    }
  }

  function getLogs() {
    if (appState.clusterStatus.status != ClusterStates.DELETED) {
      window.ddClient.extension.vm.cli
        .exec(`/backend/clustermgr`, [`logs`]) // cluster-create-logs
        .then((cmdResult: any) => {
          let status: IClusterStatus = cmdResult.parseJsonObject();
          let logs: string[] = [];
          logs[0] = status.output || '';
          // TODO: Convert string to string[]
          dispatch({ type: 'CLUSTER_LOGS', payload: logs } as Actions);
        });
    }
  }

  async function getClusterStats() {
    if (appState.clusterStatus.status === ClusterStates.RUNNING) {
      window.ddClient.extension.vm.cli
        .exec(`/backend/clustermgr`, [`stats`])
        .then((cmdResult: any) => {
          // TODO: Add error handler
          let stats: IClusterResourceStats = cmdResult.parseJsonObject().stats;
          dispatch({ type: 'CLUSTER_STATS', payload: stats } as Actions);
        });
    }
  }

  const [selectedTab, setSelectedTab] = React.useState(0);

  const selectTab = (event: React.SyntheticEvent, newValue: number) => {
    if (newValue === TABS.LOGS) {
      getLogs();
    }
    // Start a trigger to refresh stats every 2 seconds while the Stats tab is selected
    setSelectedTab(newValue);
  };

  function clusterStateIcon(state: IAppState) {
    let showIcon: boolean = true;
    let iconColor: string = '';
    switch (state.clusterStatus.status) {
      case ClusterStates.UNKNOWN:
      case ClusterStates.NOT_EXISTS:
        showIcon = false;
        break;
      case ClusterStates.DELETED:
      case ClusterStates.DELETING:
      case ClusterStates.ERROR:
        iconColor = 'red';
        break;
      case ClusterStates.CREATING:
      case ClusterStates.INITIALIZING:
        iconColor = 'orange';
        break;
      case ClusterStates.RUNNING:
        iconColor = 'green';
        break;
    }
    return (
      <CircleRoundedIcon sx={{ fontSize: 12 }} style={{ color: iconColor }} />
    );
  }

  return (
    <Container>
      <FeedbackLink />
      <Stack spacing={2}>
        <Box py={5} textAlign="left" alignItems="left" sx={{ height: '25vh' }}>
          <Stack
            direction="row"
            sx={{ display: 'flex', justifyContent: 'space-between' }}
          >
            <Stack sx={{ display: 'flex', justifyContent: 'space-between' }}>
              <Box sx={{ mt: 2 }}>
                <img
                  src="TCE-logo.svg"
                  width="300"
                  className="Header-logo"
                  alt="logo"
                />
              </Box>
              <Box sx={{ ml: 11 }}>
                <Typography style={{ fontWeight: 600 }}>
                  &nbsp;{clusterStateIcon(appState)}&nbsp;&nbsp;
                  {appState.clusterStatus?.description}&nbsp;&nbsp;
                  {appState.clusterStatus?.isError}&nbsp;&nbsp;
                  {appState.clusterStatus?.errorMessage}
                </Typography>
              </Box>
            </Stack>
            <Stack
              sx={{
                mt: 4,
                display: 'flex',
                alignItems: 'center',
                width: '50%',
                justifyContent: 'right',
              }}
              direction="row"
              spacing={2}
              justifyContent="center"
            >
              {(appState.clusterStatus.status === ClusterStates.DELETED ||
                appState.clusterStatus.status === ClusterStates.NOT_EXISTS) && (
                <Button variant="contained" onClick={createCluster}>
                  Create
                </Button>
              )}
              {(appState.clusterStatus.status === ClusterStates.RUNNING ||
                appState.clusterStatus.status === ClusterStates.ERROR) && (
                <Button variant="contained" onClick={deleteCluster}>
                  Delete
                </Button>
              )}
            </Stack>
          </Stack>
        </Box>
        <Box sx={{ width: '100%', height: '65vh' }}>
          <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
            <Tabs value={selectedTab} onChange={selectTab}>
              <Tab label="Logs" {...a11yProps(TABS.LOGS)} />
              {appState.clusterStatus.status === ClusterStates.RUNNING && (
                <Tab
                  label="KubeConfig"
                  {...a11yProps(TABS.KUBECONFIG)}
                  onClick={getClusterKubeconfig}
                />
              )}
              {/* {appState.clusterStatus.status === ClusterStates.RUNNING && (
                <Tab label="Packages" {...a11yProps(TABS.PACKAGES)} />
              )}
              {appState.clusterStatus.status === ClusterStates.RUNNING && (
                <Tab
                  label="Stats"
                  {...a11yProps(TABS.STATS)}
                  onClick={getClusterStats}
                />
              )} */}
            </Tabs>
          </Box>
          <TabPanel value={selectedTab} index={TABS.LOGS}>
            <LogTailerBox log={appState.logs} />
          </TabPanel>
          <TabPanel value={selectedTab} index={TABS.KUBECONFIG}>
            <KubeconfigBox kubeconfig={appState.kubeconfig} />
          </TabPanel>
          {/* <TabPanel value={selectedTab} index={TABS.PACKAGES}>
            <PackagesBox/>
          </TabPanel>
          <TabPanel value={selectedTab} index={TABS.STATS}>
            <StatsPane/>
          </TabPanel> */}
        </Box>
      </Stack>
    </Container>
  );
}
