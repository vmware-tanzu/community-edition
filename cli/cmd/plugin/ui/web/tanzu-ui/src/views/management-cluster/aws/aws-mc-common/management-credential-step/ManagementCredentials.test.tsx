// React imports
import React from 'react';

// Library imports
import { act } from 'react-dom/test-utils';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { rest } from 'msw';
import { setupServer } from 'msw/lib/node';
import userEvent from '@testing-library/user-event';

// App imports
import ManagementCredentials from './ManagementCredentials';

describe('ManagementCredential component', () => {
    const user = userEvent.setup();
    const server = setupServer(
        rest.get('/api/provider/aws/profiles', (req, res, ctx) => {
            return res(ctx.status(200), ctx.json(['profile1', 'profile2', 'profile3', 'profile4']));
        }),
        rest.get('/api/provider/aws/keypair', (req, res, ctx) => {
            return res(
                ctx.status(200),
                ctx.json([
                    { id: '1', name: 'us-west-2-kp', thumbprint: '' },
                    { id: '2', name: 'eu-west-1-kp', thumbprint: '' },
                    { id: '3', name: 'eu-west-2-kp', thumbprint: '' },
                ])
            );
        }),
        rest.get('/api/provider/aws/profiles', (req, res, ctx) => {
            return res(ctx.status(200), ctx.json(['profile1', 'profile2', 'profile3', 'profile4']));
        }),
        rest.get('/api/provider/aws/regions', (req, res, ctx) => {
            return res(
                ctx.status(200),
                ctx.json([
                    'us-east-1',
                    'us-east-2',
                    'us-west-1',
                    'us-west-2',
                    'eu-central-1',
                    'eu-east-1',
                    'eu-east-2',
                    'ap-east-1',
                    'ap-south-1',
                    'ca-central-1',
                ])
            );
        }),
        rest.post('/api/provider/aws', (req, res, ctx) => {
            return res(ctx.status(200));
        }),
        rest.get('api/provider/aws/osimages', (req, res, ctx) => {
            return res(ctx.status(200));
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
    it('should have credential profile selected by default', async () => {
        render(<ManagementCredentials />);
        expect(await screen.findByText('profile1')).toBeInTheDocument();
    });
    it('should show one time credential', async () => {
        render(<ManagementCredentials />);
        await user.click(screen.getByDisplayValue('ONE_TIME'));
        expect(screen.getByText('Access key ID')).toBeInTheDocument();
    });
    it('should connect to aws', async () => {
        render(<ManagementCredentials />);
        await screen.findByText('us-east-1');
        fireEvent.change(screen.getByTestId('region-select'), { target: { value: 'us-east-2' } });
        fireEvent.click(screen.getByText('CONNECT'));
        expect(await screen.findByText('Connected to AWS')).toBeInTheDocument();
        expect(screen.getByText('CONNECT')).toHaveAttribute('disabled');
    });
    it('should select a key pair', async () => {
        render(<ManagementCredentials />);
        await screen.findByText('us-east-1');
        fireEvent.change(screen.getByTestId('region-select'), { target: { value: 'us-east-2' } });
        fireEvent.click(screen.getByText('CONNECT'));
        expect(await screen.findByText('Connected to AWS')).toBeInTheDocument();
        expect(screen.getByText('CONNECT')).toHaveAttribute('disabled');
        const el = await screen.findByText('us-west-2-kp');
        const keypairEl = screen.getByTestId('ec2keypair-select');
        await fireEvent.change(keypairEl, { target: { value: 'us-west-2-kp' } });
        expect((el as HTMLOptionElement).selected).toBeTruthy();
    });
    it('should click on refresh button', async () => {
        render(<ManagementCredentials />);
        await screen.findByText('us-east-1');
        fireEvent.change(screen.getByTestId('region-select'), { target: { value: 'us-east-2' } });
        fireEvent.click(screen.getByText('CONNECT'));
        expect(await screen.findByText('Connected to AWS')).toBeInTheDocument();
        expect(screen.getByText('CONNECT')).toHaveAttribute('disabled');
        server.use(
            rest.get('/api/provider/aws/keypair', (req, res, ctx) => {
                return res(
                    ctx.status(200),
                    ctx.json([
                        { id: '1', name: 'us-west-2-kp-refreshed', thumbprint: '' },
                        { id: '2', name: 'eu-west-1-kp-refreshed', thumbprint: '' },
                        { id: '3', name: 'eu-west-2-kp-refreshed', thumbprint: '' },
                    ])
                );
            })
        );
        fireEvent.click(screen.getByText('REFRESH'));
        expect(await screen.findByText('us-west-2-kp-refreshed')).toBeInTheDocument();
    });
    describe('should change the button from connected to connect', () => {
        it('change input value', async () => {
            render(<ManagementCredentials />);
            await user.click(screen.getByDisplayValue('ONE_TIME'));
            await screen.findByText('us-east-1');
            fireEvent.change(screen.getByTestId('region-select'), { target: { value: 'us-east-2' } });
            fireEvent.click(screen.getByText('CONNECT'));
            expect(await screen.findByText('Connected to AWS')).toBeInTheDocument();
            expect(screen.getByText('CONNECT')).toHaveAttribute('disabled');
            fireEvent.change(screen.getByPlaceholderText('Access key ID'), { target: { value: 'myTestAccessKeyId' } });
            expect(await screen.findByText('CONNECT')).toBeInTheDocument();
        });
        it('change profile value', async () => {
            render(<ManagementCredentials />);
            await screen.findByText('us-east-1');
            fireEvent.change(screen.getByTestId('region-select'), { target: { value: 'us-east-2' } });
            fireEvent.click(screen.getByText('CONNECT'));
            expect(await screen.findByText('Connected to AWS')).toBeInTheDocument();
            expect(screen.getByText('CONNECT')).toHaveAttribute('disabled');
            await screen.findByText('profile2');
            fireEvent.change(screen.getByTestId('profile-select'), { target: { value: 'profile2' } });
            expect(await screen.findByText('CONNECT')).toBeInTheDocument();
        });
    });
});
