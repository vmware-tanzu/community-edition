import { DockerMuiThemeProvider } from '@docker/docker-mui-theme';
import CssBaseline from '@mui/material/CssBaseline';
import { GlobalAppStateProvider, RoutesProvider } from './providers';
import React from 'react';
import ReactDOM from 'react-dom';

import reportWebVitals from './reportWebVitals';

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

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals(console.log);