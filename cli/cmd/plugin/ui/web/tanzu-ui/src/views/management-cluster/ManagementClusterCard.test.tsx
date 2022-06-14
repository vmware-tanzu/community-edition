// React imports
import React from 'react';

// Library imports
import { render, screen, waitFor } from '@testing-library/react';

// App imports
import ManagementClusterCard from './ManagementClusterCard';

describe('ManagementClusterCard component', () => {
    const mockProps = {
        name: 'test-mgmt-cluster',
        path: '/some/path',
        context: '/some/context',
        confirmDeleteCallback: (arg?: string) => {
            return;
        },
    };

    test('should render', async () => {
        const view = render(<ManagementClusterCard {...mockProps} />);
        await waitFor(() => {
            expect(view).toBeDefined();
        });
    });

    test('should display the management cluster name', async () => {
        render(<ManagementClusterCard {...mockProps} />);
        expect(await screen.findByText('test-mgmt-cluster')).toBeInTheDocument();
    });

    test('should display the management cluster path', async () => {
        render(<ManagementClusterCard {...mockProps} />);
        expect(await screen.findByText('/some/path')).toBeInTheDocument();
    });

    test('should display the management cluster context', async () => {
        render(<ManagementClusterCard {...mockProps} />);
        expect(await screen.findByText('/some/context')).toBeInTheDocument();
    });

    test('should display a delete button', async () => {
        render(<ManagementClusterCard {...mockProps} />);
        expect(await screen.findByText('Delete')).toBeInTheDocument();
    });
});
