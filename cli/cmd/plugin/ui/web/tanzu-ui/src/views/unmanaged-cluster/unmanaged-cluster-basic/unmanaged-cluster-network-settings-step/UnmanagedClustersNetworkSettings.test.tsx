// React imports
import React from 'react';

// Library imports
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

// App imports
import { UNMANAGED_CLUSTER_FIELDS } from '../../unmanaged-cluster-common/UnmanagedCluster.constants';
import { UNMANAGED_PLACEHOLDER_VALUES } from '../../unmanaged-cluster-common/unmanaged.defaults';
import UnmanagedClusterNetworkSettings from './UnmanagedClustersNetworkSettings';
import { act } from 'react-dom/test-utils';

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
        protocol: 'udp',
        combined: '122.0.0.2:60:60/udp',
    },
    {
        ip: '124.0.0.5',
        nodePort: '92',
        hostPort: '92',
        protocol: 'sctp',
        combined: '124.0.0.5:92:92/sctp',
    },
];

jest.mock('react-router-dom', () => ({
    ...(jest.requireActual('react-router-dom') as any),
    useNavigate: () => jest.fn(),
}));

describe('UnmanagedNetworksettings component', () => {
    const user = userEvent.setup();
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
    it('make node to host port mapping string out of all input values', async () => {
        render(<UnmanagedClusterNetworkSettings />);
        await screen.findByText('IP Address');
        const input = screen.getByPlaceholderText('127.0.0.1');
        fireEvent.change(input, { target: { value: '127.0.0.4' } });
        // fireEvent.change(screen.getByPlaceholderText('Node port'), {
        //     target: { value: testUnamangedClusterNodeToHostPortMapping[0].nodePort },
        // });
        // fireEvent.change(screen.getByPlaceholderText('Host port'), {
        //     target: { value: testUnamangedClusterNodeToHostPortMapping[0].hostPort },
        // });
        // fireEvent.change(screen.getByText('tcp'), {
        //     target: { value: testUnamangedClusterNodeToHostPortMapping[0].protocol },
        // });
        console.log(screen.getByLabelText(/IP Address/i));
        console.log(testUnamangedClusterNodeToHostPortMapping[0].combined);
        let result = await screen.findByText(testUnamangedClusterNodeToHostPortMapping[0].ip);
        expect(result).toBeInTheDocument();
    });
    it('should have create button', async () => {
        render(<UnmanagedClusterNetworkSettings />);
        expect(await screen.findByTestId('create-cluster-btn')).toBeInTheDocument();
    });
});
