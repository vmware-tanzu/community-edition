// React imports
import React from 'react';

// Library imports
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { rest } from 'msw';
import { setupServer } from 'msw/lib/node';

// App imports
import ManagementClusterInventory from './ManagementClusterInventory';

jest.mock('react-router-dom', () => ({
    ...(jest.requireActual('react-router-dom') as any),
    useNavigate: () => jest.fn(),
}));

describe('ManagementClusterInventory component', () => {
    const server = setupServer(
        rest.get('/api/management', (req, res, ctx) => {
            return res(
                ctx.status(200),
                ctx.json([
                    { name: 'aws-test-cluster-1', context: '/some/context/here/1', path: '/some/path/here/1', provider: 'aws' },
                    { name: 'vsphere-other-cluster', context: '/some/context/here/2', path: '/some/path/here/2', provider: 'vsphere' },
                    { name: 'docker-foobar-cluster', context: '/some/context/here/3', path: '/some/path/here/3', provider: 'docker' },
                    { name: 'azure-clown-cluster', context: '/some/context/here/4', path: '/some/path/here/4', provider: 'azure' },
                ])
            );
        })
    );

    beforeAll(() => server.listen({ onUnhandledRequest: 'bypass' }));
    afterEach(() => server.resetHandlers());
    afterAll(() => server.close());

    const clickFirstDeleteButton = async function () {
        const deleteBtns = await screen.findAllByText('Delete');
        const firstDeleteBtn = deleteBtns[0];
        fireEvent.click(firstDeleteBtn);
    };

    test('should render', async () => {
        const view = render(<ManagementClusterInventory />);
        await waitFor(() => {
            expect(view).toBeDefined();
        });
    });

    test('should always display a button to create a management cluster', async () => {
        render(<ManagementClusterInventory />);

        expect(await screen.findByText('create a management cluster')).toBeInTheDocument();
    });

    test('should display four management cluster cards when management clusters are present', async () => {
        render(<ManagementClusterInventory />);
        const managementClusterCards = await screen.findAllByTestId('management-cluster-card');
        expect(managementClusterCards.length).toBe(4);
    });

    test('should display title, path and context for a management cluster card', async () => {
        render(<ManagementClusterInventory />);
        expect(await screen.findByText('aws-test-cluster-1')).toBeInTheDocument();
        expect(await screen.findByText('/some/context/here/1')).toBeInTheDocument();
        expect(await screen.findByText('/some/path/here/1')).toBeInTheDocument();
    });

    test('delete button should open modal confirmation', async () => {
        render(<ManagementClusterInventory />);

        await clickFirstDeleteButton();

        expect(await screen.findByTestId('confirm-delete-cluster-modal')).toBeInTheDocument();
    });

    test('delete modal confirmation cancel button should close modal window', async () => {
        render(<ManagementClusterInventory />);

        await clickFirstDeleteButton();

        const cancelBtn = await screen.findByText('Cancel');
        fireEvent.click(cancelBtn);
        expect(screen.queryByTestId('confirm-delete-cluster-modal')).not.toBeInTheDocument();

        const managementClusterCards = await screen.findAllByTestId('management-cluster-card');
        expect(managementClusterCards.length).toBe(4);
    });

    test('delete modal confirmation Delete button should delete management cluster', async () => {
        render(<ManagementClusterInventory />);

        server.use(
            rest.get('/api/management', (req, res, ctx) => {
                return res(
                    ctx.status(200),
                    ctx.json([
                        { name: 'vsphere-other-cluster', context: '/some/context/here/2', path: '/some/path/here/2', provider: 'vsphere' },
                        { name: 'docker-foobar-cluster', context: '/some/context/here/3', path: '/some/path/here/3', provider: 'docker' },
                        { name: 'azure-clown-cluster', context: '/some/context/here/4', path: '/some/path/here/4', provider: 'azure' },
                    ])
                );
            })
        );

        await clickFirstDeleteButton();

        const deleteBtn = await screen.findByTestId('delete-cluster-btn');
        fireEvent.click(deleteBtn);

        const managementClusterCards = await screen.findAllByTestId('management-cluster-card');
        expect(managementClusterCards.length).toBe(3);
    });

    test('should display messaging when no management clusters are present', async () => {
        render(<ManagementClusterInventory />);

        server.use(
            rest.get('/api/management', (req, res, ctx) => {
                return res(ctx.status(200), ctx.json([]));
            })
        );
        expect(await screen.findByTestId('no-clusters-messaging')).toBeInTheDocument();
    });
});
