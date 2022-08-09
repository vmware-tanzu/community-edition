// React imports
import React from 'react';

// Library imports
import { render, screen, waitFor } from '@testing-library/react';

// App imports
import UnmanagedClusterNetworkSettings from './UnmanagedClustersNetworkSettings';

jest.mock('react-router-dom', () => ({
    ...(jest.requireActual('react-router-dom') as any),
    useNavigate: () => jest.fn(),
}));

describe('UnmanagedNetworksettings component', () => {

    it('should render', async () => {
        const view = render(<UnmanagedClusterNetworkSettings />);
        await waitFor(() => {
            expect(view).toBeDefined();
        });
    });
    it('should have CNI form', async () => {
        render(<UnmanagedClusterNetworkSettings />);
        expect(await screen.findByText('Container Network Interface (CNI) provider')).toBeInTheDocument();
    });
    it('should have cluster service CIDR form', async () => {
        render(<UnmanagedClusterNetworkSettings />);
        expect(await screen.findByText('Cluster service CIDR')).toBeInTheDocument();
    });
    it('should verify cluster service CIDR form', async () => {
        render(<UnmanagedClusterNetworkSettings />);
        expect(await screen.findByText('Cluster service CIDR')).toBeInTheDocument();
    });
    it('should have cluster POD CIDR form', async () => {
        render(<UnmanagedClusterNetworkSettings />);
        expect(await screen.findByText('Cluster POD CIDR')).toBeInTheDocument();
    });
    it('should have create button', async () => {
        render(<UnmanagedClusterNetworkSettings />);
        expect(await screen.findByTestId('create-cluster-btn')).toBeInTheDocument();
    });
});
