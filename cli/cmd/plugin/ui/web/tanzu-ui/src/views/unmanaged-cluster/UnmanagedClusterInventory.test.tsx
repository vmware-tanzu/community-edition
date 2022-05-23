// React imports
import React from 'react';
import { BrowserRouter } from 'react-router-dom';

// Library imports
import { render } from '@testing-library/react';
import '@testing-library/jest-dom';

// App imports
import UnmanagedClusterInventory from './UnmanagedClusterInventory';

describe('UnmanagedClusterLanding', () => {
    test('should render', () => {
        const view = render(<UnmanagedClusterInventory />, { wrapper: BrowserRouter });
        expect(view).toBeDefined();
    });
});
