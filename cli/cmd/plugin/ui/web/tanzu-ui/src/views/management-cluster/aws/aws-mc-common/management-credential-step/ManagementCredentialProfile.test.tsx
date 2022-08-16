// React imports
import React from 'react';

// Library imports
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { rest } from 'msw';
import { setupServer } from 'msw/lib/node';

// App imports
import ManagementCredentialProfile from './ManagementCredentialProfile';
import { FormProvider } from 'react-hook-form';

const regionsMock = [
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
];
const useFormMock = {
    ...jest.requireActual('react-hook-form'),
    useForm: () => ({
        formState: {
            errors: {},
        },
        setValue: jest.fn(),
        register: (name: string, obj: any) => {
            return {
                name,
                ...obj,
            };
        },
        reset: jest.fn(),
    }),
};
const methods = useFormMock.useForm();
const mockSelectCallback = jest.fn();

describe('ManagementCredentialProfile component', () => {
    const server = setupServer(
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
        })
    );

    beforeAll(() => server.listen());
    afterEach(() => server.resetHandlers());
    afterAll(() => server.close());

    it('should render', async () => {
        const view = render(
            <FormProvider {...methods}>
                <ManagementCredentialProfile selectCallback={mockSelectCallback} />
            </FormProvider>
        );
        await waitFor(() => {
            expect(view).toBeDefined();
        });
    });
    it('select options should contain all profiles', async () => {
        render(
            <FormProvider {...methods}>
                <ManagementCredentialProfile selectCallback={mockSelectCallback} />
            </FormProvider>
        );
        await screen.findByText('profile1');
        const profiles = ['profile2', 'profile3', 'profile4'];
        for (let i = 0; i < profiles.length; i++) {
            const profileEl = await screen.findByText(profiles[i]);
            expect(profileEl).toBeInTheDocument();
        }
    });
    it('select profile should fire handler method', async () => {
        render(
            <FormProvider {...methods}>
                <ManagementCredentialProfile selectCallback={mockSelectCallback} />
            </FormProvider>
        );
        await screen.findByText('profile2');
        fireEvent.change(screen.getByTestId('profile-select'), { target: { value: 'profile2' } });
        expect(mockSelectCallback).toBeCalledTimes(1);
    });
    it('select options should contain all regions', async () => {
        render(
            <FormProvider {...methods}>
                <ManagementCredentialProfile selectCallback={mockSelectCallback} />
            </FormProvider>
        );
        await screen.findByText('us-east-1');
        for (let i = 0; i < regionsMock.length; i++) {
            const profileEl = screen.getByText(regionsMock[i]);
            expect(profileEl).toBeInTheDocument();
        }
    });
    it('select region should fire handler method', async () => {
        render(
            <FormProvider {...methods}>
                <ManagementCredentialProfile selectCallback={mockSelectCallback} />
            </FormProvider>
        );
        await screen.findByText('us-east-1');
        for (let i = 0; i < regionsMock.length; i++) {
            const regionEl = screen.getByText(regionsMock[i]);
            expect(regionEl).toBeInTheDocument();
        }
        fireEvent.change(screen.getByTestId('region-select'), { target: { value: 'us-east-2' } });
        expect(mockSelectCallback).toBeCalled();
    });
});
