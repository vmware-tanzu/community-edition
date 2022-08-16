// Library imports
import { fireEvent, render, screen, waitFor } from '@testing-library/react';

// React imports
import React from 'react';
// App imports
import UnmanagedClusterInventory from './UnmanagedClusterInventory';
import { rest } from 'msw';
import { setupServer } from 'msw/lib/node';

const testUnamangedClusterArray = [
    {
        name: 'work-space-cluster',
        provider: 'MiniKube',
        status: 'Running',
    },
    {
        name: 'tanzu-cluster',
        provider: 'Kind',
        status: 'Stopped',
    },
    {
        name: 'ui-cluster',
        provider: 'Kind',
        status: 'Unknown',
    },
];

jest.mock('react-router-dom', () => ({
    ...(jest.requireActual('react-router-dom') as any),
    useNavigate: () => jest.fn(),
}));

describe('UnmanagedClusterInventory component', () => {
    const server = setupServer(
        rest.get('/api/unmanaged', (req, res, ctx) => {
            return res(ctx.status(200), ctx.json(testUnamangedClusterArray));
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
        const view = render(<UnmanagedClusterInventory />);
        await waitFor(() => {
            expect(view).toBeDefined();
        });
    });

    test('should always display a button to create an unmanaged cluster', async () => {
        render(<UnmanagedClusterInventory />);
        expect(await screen.findByText('create an unmanaged cluster')).toBeInTheDocument();
    });

    test('should display three unmanaged cluster cards when unmanaged clusters are present', async () => {
        render(<UnmanagedClusterInventory />);
        const unmanagedClusterCards = await screen.findAllByTestId('unmanaged-cluster-card');
        expect(unmanagedClusterCards.length).toBe(testUnamangedClusterArray.length);
    });

    test('should display name, provider and status for an unmanaged cluster card', async () => {
        render(<UnmanagedClusterInventory />);
        expect(await screen.findByText(testUnamangedClusterArray[0].name)).toBeInTheDocument();
        expect(await screen.findByText(testUnamangedClusterArray[0].provider)).toBeInTheDocument();
        expect(await screen.findByText(testUnamangedClusterArray[0].status)).toBeInTheDocument();
    });

    test('delete button should open modal confirmation', async () => {
        render(<UnmanagedClusterInventory />);

        await clickFirstDeleteButton();

        expect(await screen.findByTestId('confirm-delete-cluster-modal')).toBeInTheDocument();
    });

    test('delete modal confirmation cancel button should close modal window', async () => {
        render(<UnmanagedClusterInventory />);

        await clickFirstDeleteButton();

        const cancelBtn = await screen.findByText('Cancel');
        fireEvent.click(cancelBtn);
        expect(screen.queryByTestId('confirm-delete-cluster-modal')).not.toBeInTheDocument();

        const umanagedClusterCards = await screen.findAllByTestId('unmanaged-cluster-card');
        expect(umanagedClusterCards.length).toBe(3);
    });

    test('delete modal confirmation Delete button should delete unmanaged cluster', async () => {
        render(<UnmanagedClusterInventory />);

        server.use(
            rest.get('/api/unmanaged', (req, res, ctx) => {
                return res(
                    ctx.status(200),
                    ctx.json([
                        {
                            name: 'tanzu-cluster',
                            provider: 'Kind',
                            status: 'Stopped',
                        },
                        {
                            name: 'ui-cluster',
                            provider: 'MiniKube',
                            status: 'Unknown',
                        },
                    ])
                );
            })
        );

        await clickFirstDeleteButton();

        const deleteBtn = await screen.findByTestId('delete-cluster-btn');
        fireEvent.click(deleteBtn);

        const umanagedClusterCards = await screen.findAllByTestId('unmanaged-cluster-card');
        expect(umanagedClusterCards.length).toBe(testUnamangedClusterArray.length - 1);
    });

    test('should display messaging when no unmanaged clusters are present', async () => {
        render(<UnmanagedClusterInventory />);

        server.use(
            rest.get('/api/unmanaged', (req, res, ctx) => {
                return res(ctx.status(200), ctx.json([]));
            })
        );
        expect(await screen.findByTestId('no-clusters-messaging')).toBeInTheDocument();
    });
});
