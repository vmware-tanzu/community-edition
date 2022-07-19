import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { rest } from 'msw';
import { setupServer } from 'msw/lib/node';

import ManagementCredentials from './ManagementCredentials';

const regionsMock = ['West US', 'North central US', 'South central US', 'Central US', 'East US', 'East US 2'];
const azureEnvironment = ['Public Cloud', 'US Government Cloud'];
const formFieldItem = [
    { placeholder: 'Tenant ID', field: 'TENANT_ID' },
    { placeholder: 'Client ID', field: 'CLIENT_ID' },
    { placeholder: 'Client Secret', field: 'CLIENT_SECRET' },
    { placeholder: 'Subscription ID', field: 'SUBSCRIPTION_ID' },
];

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

    it('should handle form field change', async () => {
        const handleValueChange = jest.fn();
        render(<ManagementCredentials handleValueChange={handleValueChange} />);
        for (let i = 0; i < formFieldItem.length; i++) {
            const input = screen.getByPlaceholderText(formFieldItem[i].placeholder);
            fireEvent.change(input, {
                target: { value: 'test' + formFieldItem[i].placeholder },
            });
            expect(handleValueChange).toHaveBeenCalled();
            expect(handleValueChange).toBeCalledWith(
                'INPUT_CHANGE',
                formFieldItem[i].field,
                'test' + formFieldItem[i].placeholder,
                undefined,
                {}
            );
        }
        const selectItem = screen.getByTestId('azure-environment-select');
        fireEvent.change(selectItem, { target: { value: 'AzurePublicCloud' } });
        expect(handleValueChange).toHaveBeenCalled();
        expect(handleValueChange).toBeCalledWith('INPUT_CHANGE', 'AZURE_ENVIRONMENT', 'AzurePublicCloud', undefined, {});
    });

    it('should connect to Azure', async () => {
        render(<ManagementCredentials handleValueChange={jest.fn} />);
        fireEvent.click(screen.getByText('CONNECT'));
        expect(await screen.findByText('Connected to Azure')).toBeInTheDocument();
    });

    it('should select a region', async () => {
        render(<ManagementCredentials handleValueChange={jest.fn} />);
        fireEvent.click(screen.getByText('CONNECT'));
        expect(await screen.findByText('Connected to Azure')).toBeInTheDocument();
        const el = await screen.findByText('West US');
        for (let i = 0; i < regionsMock.length; i++) {
            const profileEl = screen.getByText(regionsMock[i]);
            expect(profileEl).toBeInTheDocument();
        }
        const keypairEl = screen.getByTestId('region-select');
        fireEvent.change(keypairEl, { target: { value: 'westus' } });
        console.log((el as HTMLOptionElement).value);
        expect((el as HTMLOptionElement).selected).toBeTruthy();
    });

    it('should change the button from connected to connect', async () => {
        render(<ManagementCredentials handleValueChange={jest.fn} />);
        fireEvent.click(screen.getByText('CONNECT'));
        expect(await screen.findByText('Connected to Azure')).toBeInTheDocument();
        await fireEvent.change(screen.getByPlaceholderText('Tenant ID'), {
            target: { value: 'myTestAccessKeyId' },
        });
        const keypairEl = screen.getByPlaceholderText('Tenant ID');
        console.log((keypairEl as HTMLElement).getAttribute('placeholder'));
        expect(await screen.findByText('CONNECT')).toBeInTheDocument();
    });
});
