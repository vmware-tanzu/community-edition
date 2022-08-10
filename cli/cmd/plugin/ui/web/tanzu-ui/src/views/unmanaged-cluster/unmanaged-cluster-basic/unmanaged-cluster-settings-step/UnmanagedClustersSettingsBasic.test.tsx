// React imports
import React from 'react';

// Library imports
import { render, screen, waitFor } from '@testing-library/react';

// App imports
import UnmanagedClusterSettings from './UnmanagedClustersSettingsBasic';

jest.mock('react-router-dom', () => ({
    ...(jest.requireActual('react-router-dom') as any),
    useNavigate: () => jest.fn(),
}));

describe('Unmanagedsettings component', () => {
    it('should render', async () => {
        const view = render(<UnmanagedClusterSettings />);
        await waitFor(() => {
            expect(view).toBeDefined();
        });
    });
    it('should have next button', async () => {
        render(<UnmanagedClusterSettings />);
        expect(screen.getByText('NEXT')).toBeInTheDocument();
    });
});
