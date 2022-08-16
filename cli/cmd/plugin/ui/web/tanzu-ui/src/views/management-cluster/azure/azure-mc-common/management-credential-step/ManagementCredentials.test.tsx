// React imports
import React from 'react';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';

// Library imports
import { rest } from 'msw';
import { setupServer } from 'msw/lib/node';

// App imports
import ManagementCredentials from './ManagementCredentials';

const regionsMock = ['West US', 'North central US', 'South central US', 'Central US', 'East US', 'East US 2'];
const azureEnvironment = ['Public Cloud', 'US Government Cloud'];

describe('ManagementCredential component', () => {
    const server = setupServer(
        rest.post('/api/provider/azure', (req, res, ctx) => {
            return res(ctx.status(200));
        }),
        rest.get('/api/provider/azure/regions', (req, res, ctx) => {
            return res(
                ctx.status(200),
                ctx.json([
                    {
                        name: 'westus',
                        displayName: 'West US',
                    },
                    {
                        name: 'northcentralus',
                        displayName: 'North central US',
                    },
                    {
                        name: 'southcentralus',
                        displayName: 'South central US',
                    },
                    {
                        name: 'centralus',
                        displayName: 'Central US',
                    },
                    {
                        name: 'eastus',
                        displayName: 'East US',
                    },
                    {
                        name: 'eastus2',
                        displayName: 'East US 2',
                    },
                ])
            );
        })
    );

    beforeAll(() => server.listen());
    afterEach(() => server.resetHandlers());
    afterAll(() => server.close());

    it('should render', async () => {
        const view = render(<ManagementCredentials />);
        await waitFor(() => {
            expect(view).toBeDefined();
        });
    });

    it('should contain all Azure Environment', async () => {
        render(<ManagementCredentials />);
        for (let i = 0; i < azureEnvironment.length; i++) {
            const profileEl = screen.getByText(azureEnvironment[i]);
            expect(profileEl).toBeInTheDocument();
        }
    });

    it('should connect to Azure', async () => {
        render(<ManagementCredentials />);
        fireEvent.click(screen.getByText('CONNECT'));
        expect(await screen.findByText('Connected to Azure')).toBeInTheDocument();
    });

    it('should select a region', async () => {
        render(<ManagementCredentials />);
        fireEvent.click(screen.getByText('CONNECT'));
        expect(await screen.findByText('Connected to Azure')).toBeInTheDocument();
        const el = await screen.findByText('West US');
        for (let i = 0; i < regionsMock.length; i++) {
            const profileEl = screen.getByText(regionsMock[i]);
            expect(profileEl).toBeInTheDocument();
        }
        const keypairEl = screen.getByTestId('region-select');
        // TODO: The issue caused by setValue method. This logic should be revisited.
        // eslint-disable-next-line testing-library/no-unnecessary-act
        await act(async () => {
            fireEvent.change(keypairEl, { target: { value: 'westus' } });
        });
        expect((el as HTMLOptionElement).selected).toBeTruthy();
    });

    it('should change the button from connected to connect', async () => {
        render(<ManagementCredentials />);
        fireEvent.click(screen.getByText('CONNECT'));
        expect(await screen.findByText('Connected to Azure')).toBeInTheDocument();
        // TODO: The issue caused by setValue method. This logic should be revisited.
        // eslint-disable-next-line testing-library/no-unnecessary-act
        await act(async () => {
            fireEvent.change(screen.getByPlaceholderText('Tenant ID'), {
                target: { value: 'myTestAccessKeyId' },
            });
        });
        expect(await screen.findByText('CONNECT')).toBeInTheDocument();
    });
});
