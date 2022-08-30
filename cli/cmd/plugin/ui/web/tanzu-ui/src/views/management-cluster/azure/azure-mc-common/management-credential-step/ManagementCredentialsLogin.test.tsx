// React imports
import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';

// Library imports
import { FormProvider } from 'react-hook-form';

// App imports
import { CONNECTION_STATUS } from '../../../../../shared/components/ConnectionNotification/ConnectionNotification';
import ManagementCredentialsLogin from './ManagementCredentialsLogin';

const azureEnvironment = ['Public Cloud', 'US Government Cloud'];
const formFieldItem = [
    { placeholder: 'Tenant ID', field: 'tenantId' },
    { placeholder: 'Client ID', field: 'clientId' },
    { placeholder: 'Client Secret', field: 'clientSecret' },
    { placeholder: 'Subscription ID', field: 'subscriptionId' },
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
    }),
};
const methods = useFormMock.useForm();
const handleInputChange = jest.fn();
const handleConnect = jest.fn();

describe('Login component', () => {
    it('should render', async () => {
        const view = render(
            <FormProvider {...methods}>
                <ManagementCredentialsLogin
                    status={CONNECTION_STATUS.DISCONNECTED}
                    handleConnect={jest.fn}
                    handleInputChange={handleInputChange}
                    message={''}
                />
            </FormProvider>
        );
        await waitFor(() => {
            expect(view).toBeDefined();
        });
    });

    it('should contain all Azure Environment', async () => {
        render(
            <FormProvider {...methods}>
                <ManagementCredentialsLogin
                    status={CONNECTION_STATUS.DISCONNECTED}
                    handleConnect={jest.fn}
                    handleInputChange={handleInputChange}
                    message={''}
                />
            </FormProvider>
        );
        for (let i = 0; i < azureEnvironment.length; i++) {
            const profileEl = screen.getByText(azureEnvironment[i]);
            expect(profileEl).toBeInTheDocument();
        }
    });

    it('should handle form field change', async () => {
        render(
            <FormProvider {...methods}>
                <ManagementCredentialsLogin
                    status={CONNECTION_STATUS.DISCONNECTED}
                    handleConnect={jest.fn}
                    handleInputChange={handleInputChange}
                    message={''}
                />
            </FormProvider>
        );
        for (let i = 0; i < formFieldItem.length; i++) {
            const input = screen.getByPlaceholderText(formFieldItem[i].placeholder);
            fireEvent.change(input, { target: { value: 'test' + formFieldItem[i].placeholder } });
            expect(handleInputChange).toHaveBeenCalled();
            expect(handleInputChange).toBeCalledWith(formFieldItem[i].field, 'test' + formFieldItem[i].placeholder);
        }
        const selectItem = screen.getByTestId('azure-environment-select');
        fireEvent.change(selectItem, { target: { value: 'AzurePublicCloud' } });
        expect(handleInputChange).toHaveBeenCalled();
        expect(handleInputChange).toBeCalledWith('azureEnvironment', 'AzurePublicCloud');
    });

    it('should show connected', async () => {
        render(
            <FormProvider {...methods}>
                <ManagementCredentialsLogin
                    status={CONNECTION_STATUS.DISCONNECTED}
                    handleConnect={handleConnect}
                    handleInputChange={handleInputChange}
                    message={''}
                />
            </FormProvider>
        );
        fireEvent.click(screen.getByText('CONNECT'));
        expect(handleConnect).toHaveBeenCalled();
        render(
            <FormProvider {...methods}>
                <ManagementCredentialsLogin
                    status={CONNECTION_STATUS.CONNECTED}
                    handleConnect={handleConnect}
                    handleInputChange={handleInputChange}
                    message={'Connected to Azure'}
                />
            </FormProvider>
        );
        fireEvent.click(screen.getByText('Connected to Azure'));
    });
});