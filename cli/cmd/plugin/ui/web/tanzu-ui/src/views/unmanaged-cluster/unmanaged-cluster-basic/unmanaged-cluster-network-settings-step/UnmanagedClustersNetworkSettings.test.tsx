// React imports
import React from 'react';

// Library imports
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

// App imports
import { UmcProvider } from '../../../../state-management/stores/Store.umc';
import UnmanagedClusterNetworkSettings from './UnmanagedClustersNetworkSettings';

const testUnamangedClusterNodeToHostPortMapping = [
    {
        ip: '127.0.0.4',
        nodePort: '80',
        hostPort: '80',
        protocol: 'tcp',
        combined: '127.0.0.4:80:80/tcp',
    },
    {
        ip: '122.0.0.2',
        nodePort: '60',
        hostPort: '60',
        protocol: 'tcp',
        combined: '122.0.0.2:60:60/tcp',
    },
    {
        ip: '124.0.0.5',
        nodePort: '92',
        hostPort: '92',
        protocol: 'tcp',
        combined: '124.0.0.5:92:92/tcp',
    },
];

jest.mock('react-router-dom', () => ({
    ...(jest.requireActual('react-router-dom') as any),
    useNavigate: () => jest.fn(),
}));

describe('UnmanagedNetworksettings component', () => {
    userEvent.setup();
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
    it('should have cluster POD CIDR form', async () => {
        render(<UnmanagedClusterNetworkSettings />);
        expect(await screen.findByText('Cluster POD CIDR')).toBeInTheDocument();
    });
    it('should have Node to Host Port mapping form', async () => {
        render(<UnmanagedClusterNetworkSettings />);
        expect(await screen.findByText('Node to host port mapping')).toBeInTheDocument();
    });
    //Missing test for protocol change needs to be added after hostport event change
    it('make node to host port mapping string out of all input values', async () => {
        render(
            <UmcProvider>
                <UnmanagedClusterNetworkSettings />
            </UmcProvider>
        );
        await screen.findByText('IP Address');
        const ip_input = screen.getByPlaceholderText('127.0.0.1');
        const node_input = screen.getByPlaceholderText('Node port');
        const host_input = screen.getByPlaceholderText('Host port');

        for (const unmanagedClusterMapping of testUnamangedClusterNodeToHostPortMapping) {
            fireEvent.change(ip_input, { target: { value: unmanagedClusterMapping.ip } });
            fireEvent.change(node_input, { target: { value: unmanagedClusterMapping.nodePort } });
            fireEvent.change(host_input, { target: { value: unmanagedClusterMapping.hostPort } });
            const result = await screen.findByText(unmanagedClusterMapping.combined);
            expect(result).toBeInTheDocument();
        }
    });
    it('should have create button', async () => {
        render(<UnmanagedClusterNetworkSettings />);
        expect(await screen.findByTestId('create-cluster-btn')).toBeInTheDocument();
    });
});
