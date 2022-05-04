import React from 'react';
import { Grid } from '@mui/material'
import { GlobalAppState } from 'providers/globalAppState/GlobalAppStateProvider';
import StatsGauge from './StatsGauge';

export default function StatsPane() {
    const { state: appState, dispatch } = React.useContext(GlobalAppState);

    return (
        <Grid container spacing={2} height="100%">
            <Grid item xs={6}>
                <StatsGauge
                    value={((appState.stats?.memory?.usage !== undefined) && (Number(appState.stats?.memory?.usage))) || 0}
                    scale={100}
                    label="Memory" />
            </Grid>
            <Grid item xs={6}>
                <StatsGauge
                    value={((appState.stats?.cpu?.usage !== undefined) && (Number(appState.stats?.cpu?.usage))) || 0}
                    scale={((appState.stats?.cpu?.number_cpus !== undefined) && (Number(appState.stats?.cpu?.number_cpus) * 100)) || 100}
                    label="CPU" />
            </Grid>
        </Grid>
    );
}

