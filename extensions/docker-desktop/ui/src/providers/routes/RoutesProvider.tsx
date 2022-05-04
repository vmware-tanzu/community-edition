import React from 'react';
import { HashRouter, Route, Routes } from 'react-router-dom';
import { App, Intro, NotFoundView } from 'views';
import { GlobalAppState } from 'providers/globalAppState/GlobalAppStateProvider';

export default function RoutesProvider() {
  const { state: appState, dispatch } = React.useContext(GlobalAppState);
  
  return (
    <HashRouter>
      <Routes>
        <Route path="*" element={<NotFoundView />} />
        { (appState.isClusterStarted) && (<Route path="/" element={<App createCluster={false}/>}/>) || (<Route path="/" element={<Intro />}/>)}
        <Route path="/create-cluster" element={<App  createCluster={true} />} />
      </Routes>
    </HashRouter>
  );
}