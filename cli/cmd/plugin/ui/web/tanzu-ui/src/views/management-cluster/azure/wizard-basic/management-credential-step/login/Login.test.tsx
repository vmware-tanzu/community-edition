import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import Login from './Login';
import { CONNECTION_STATUS } from '../../../../../../shared/components/ConnectionNotification/ConnectionNotification';

const azureEnvironment = ['Public Cloud', 'US Government Cloud'];
const formFieldItem = [
    { placeholder: 'Tenant ID', field: 'TENANT_ID' },
    { placeholder: 'Client ID', field: 'CLIENT_ID' },
    { placeholder: 'Client Secret', field: 'CLIENT_SECRET' },
    { placeholder: 'Subscription ID', field: 'SUBSCRIPTION_ID' },
];

const useFormMock = {
    ...jest.requireActual('react-hook-form'),
    useForm: () => ({
        formState: {
            errors: {},
        },
        setValue: jest.fn(),
        register: jest.fn(),
    }),
};
const methods = useFormMock.useForm();
const handleInputChange = jest.fn();
const handleConnect = jest.fn();

describe('Login component', () => {
    it('should render', async () => {
        const view = render(
            <Login
                status={CONNECTION_STATUS.DISCONNECTED}
                handleConnect={jest.fn}
                methods={methods}
                handleInputChange={handleInputChange}
                message={''}
            />
        );
        await waitFor(() => {
            expect(view).toBeDefined();
        });
    });

    it('should contain all Auzre Environment', async () => {
        render(
            <Login
                status={CONNECTION_STATUS.DISCONNECTED}
                handleConnect={jest.fn}
                methods={methods}
                handleInputChange={handleInputChange}
                message={''}
            />
        );
        for (let i = 0; i < azureEnvironment.length; i++) {
            const profileEl = screen.getByText(azureEnvironment[i]);
            expect(profileEl).toBeInTheDocument();
        }
    });

    it('should handle form field change', async () => {
        render(
            <Login
                status={CONNECTION_STATUS.DISCONNECTED}
                handleConnect={jest.fn}
                methods={methods}
                handleInputChange={handleInputChange}
                message={''}
            />
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
        expect(handleInputChange).toBeCalledWith('AZURE_ENVIRONMENT', 'AzurePublicCloud');
    });

    it('should show connected', async () => {
        render(
            <Login
                status={CONNECTION_STATUS.DISCONNECTED}
                handleConnect={handleConnect}
                methods={methods}
                handleInputChange={handleInputChange}
                message={''}
            />
        );
        fireEvent.click(screen.getByText('CONNECT'));
        expect(handleConnect).toHaveBeenCalled();
        render(
            <Login
                status={CONNECTION_STATUS.CONNECTED}
                handleConnect={handleConnect}
                methods={methods}
                handleInputChange={handleInputChange}
                message={'Connected to Azure'}
            />
        );
        fireEvent.click(screen.getByText('Connected to Azure'));
    });
});
