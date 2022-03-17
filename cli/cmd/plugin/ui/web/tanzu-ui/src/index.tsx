import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter } from 'react-router-dom';

import '@cds/core/global.css'; // pre-minified version breaks
import '@cds/city/css/bundles/default.min.css';
import '@cds/core/global.min.css';
import '@cds/core/styles/theme.dark.min.css';
import './index.scss';

import App from './App';
import { AppProvider } from './state-management/stores/store';

ReactDOM.render(
    <React.StrictMode>
        <AppProvider>
            <BrowserRouter basename="/ui">
                <App />
            </BrowserRouter>
        </AppProvider>
    </React.StrictMode>,
    document.getElementById('root')
);

