// React imports
import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';

// App css imports
import '@cds/core/global.css'; // pre-minified version breaks
import '@cds/city/css/bundles/default.min.css';
import '@cds/core/global.min.css';
import '@cds/core/styles/theme.dark.min.css';
import './index.scss';

// App imports
import App from './App';
import { AppProvider } from './state-management/stores/Store';

const container = document.getElementById('root') as HTMLElement;
const rootContainer = ReactDOM.createRoot(container);
rootContainer.render(
    <React.StrictMode>
        <AppProvider>
            <BrowserRouter basename="/ui">
                <App />
            </BrowserRouter>
        </AppProvider>
    </React.StrictMode>
);
