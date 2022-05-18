// React imports
import React from 'react';
import { BrowserRouter } from 'react-router-dom';

// Library imports
import { render } from '@testing-library/react';
import '@testing-library/jest-dom';

// App imports
import UnmanagedClusterLanding from './UnmanagedClusterLanding';

describe('UnmanagedClusterLanding', () => {
    test('should render', () => {
        const view = render(<UnmanagedClusterLanding />, { wrapper: BrowserRouter });
        expect(view).toBeDefined();
    });
});
