// React imports
import React from 'react';
import { BrowserRouter } from 'react-router-dom';

// Library imports
import { render } from '@testing-library/react';
import '@testing-library/jest-dom';

// App imports
import App from './App';

describe('App', () => {
    test('should render', () => {
        const view = render(
            <BrowserRouter>
                <App />
            </BrowserRouter>
        );
        expect(view).toBeDefined();
    });
});
