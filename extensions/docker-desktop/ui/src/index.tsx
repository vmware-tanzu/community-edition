import { DockerMuiThemeProvider } from '@docker/docker-mui-theme';
import CssBaseline from '@mui/material/CssBaseline';
import { GlobalAppStateProvider, RoutesProvider } from './providers';
import React from 'react';
import ReactDOM from 'react-dom';

ReactDOM.render(
  <React.StrictMode>
    <DockerMuiThemeProvider>
      <CssBaseline />
      <GlobalAppStateProvider>
        <RoutesProvider />
      </GlobalAppStateProvider>
    </DockerMuiThemeProvider>
  </React.StrictMode>,
  document.getElementById('root'),
);
