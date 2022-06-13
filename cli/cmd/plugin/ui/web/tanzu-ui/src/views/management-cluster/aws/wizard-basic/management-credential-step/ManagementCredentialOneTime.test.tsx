// React imports
import React from 'react';

// Library imports
import { fireEvent, render, screen } from '@testing-library/react';
import ManagementCredentialOneTime from './ManagementCredentialOneTime';

// App imports

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
const testInputValues = ['testSecretAccessKey', 'testSessionToken', 'testAccessId'];

describe('ManagementCredentialOneTime component', () => {
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
    it('should render', async () => {
        const view = render(
            <ManagementCredentialOneTime
                initialRegion={'us-east-1'}
                handleInputChange={jest.fn()}
                handleSelectRegion={jest.fn}
                regions={regionsMock}
                methods={useFormMock.useForm()}
            />
        );
        expect(view).toBeDefined();
    });
    it('should render initial values', () => {
        render(
            <ManagementCredentialOneTime
                initialSecretAccessKey={testInputValues[0]}
                initialSessionToken={testInputValues[1]}
                initialAccessKeyId={testInputValues[2]}
                initialRegion={'us-east-1'}
                handleInputChange={jest.fn()}
                handleSelectRegion={jest.fn}
                regions={regionsMock}
                methods={useFormMock.useForm()}
            />
        );

        testInputValues.forEach((value) => {
            expect(screen.getByDisplayValue(value)).toBeInTheDocument();
        });
    });
    it('should handle secret acccess key input change', () => {
        const handleInputChangeMock = jest.fn();
        render(
            <ManagementCredentialOneTime
                initialSecretAccessKey={testInputValues[0]}
                initialSessionToken={testInputValues[1]}
                initialAccessKeyId={testInputValues[2]}
                initialRegion={'us-east-1'}
                handleInputChange={handleInputChangeMock}
                handleSelectRegion={jest.fn}
                regions={regionsMock}
                methods={useFormMock.useForm()}
            />
        );
        const input = screen.getByPlaceholderText('Secret access key');
        fireEvent.change(input, { target: { value: 'myTestSecretAccessKey' } });
        expect(handleInputChangeMock).toHaveBeenCalled();
        expect(handleInputChangeMock).toBeCalledWith('SECRET_ACCESS_KEY', 'myTestSecretAccessKey');
    });
    it('should handle session token input change', () => {
        const handleInputChangeMock = jest.fn();
        render(
            <ManagementCredentialOneTime
                initialSecretAccessKey={testInputValues[0]}
                initialSessionToken={testInputValues[1]}
                initialAccessKeyId={testInputValues[2]}
                initialRegion={'us-east-1'}
                handleInputChange={handleInputChangeMock}
                handleSelectRegion={jest.fn}
                regions={regionsMock}
                methods={useFormMock.useForm()}
            />
        );
        const input = screen.getByPlaceholderText('Session token');
        fireEvent.change(input, { target: { value: 'myTestSessionToken' } });
        expect(handleInputChangeMock).toHaveBeenCalled();
        expect(handleInputChangeMock).toBeCalledWith('SESSION_TOKEN', 'myTestSessionToken');
    });
    it('should handle access key id input change', () => {
        const handleInputChangeMock = jest.fn();
        render(
            <ManagementCredentialOneTime
                initialSecretAccessKey={testInputValues[0]}
                initialSessionToken={testInputValues[1]}
                initialAccessKeyId={testInputValues[2]}
                initialRegion={'us-east-1'}
                handleInputChange={handleInputChangeMock}
                handleSelectRegion={jest.fn}
                regions={regionsMock}
                methods={useFormMock.useForm()}
            />
        );
        const input = screen.getByPlaceholderText('Access key ID');
        fireEvent.change(input, { target: { value: 'myTestAccessKeyId' } });
        expect(handleInputChangeMock).toHaveBeenCalled();
        expect(handleInputChangeMock).toBeCalledWith('ACCESS_KEY_ID', 'myTestAccessKeyId');
    });
    it('select options should contain all regions', () => {
        render(
            <ManagementCredentialOneTime
                initialSecretAccessKey={testInputValues[0]}
                initialSessionToken={testInputValues[1]}
                initialAccessKeyId={testInputValues[2]}
                initialRegion={'us-east-1'}
                handleInputChange={jest.fn()}
                handleSelectRegion={jest.fn}
                regions={regionsMock}
                methods={useFormMock.useForm()}
            />
        );
        screen.getByText('us-east-1');
        for (let i = 0; i < regionsMock.length; i++) {
            const profileEl = screen.getByText(regionsMock[i]);
            expect(profileEl).toBeInTheDocument();
        }
    });
    it('select region should fire handler method', () => {
        const handleSelectRegionMock = jest.fn();
        render(
            <ManagementCredentialOneTime
                initialSecretAccessKey={testInputValues[0]}
                initialSessionToken={testInputValues[1]}
                initialAccessKeyId={testInputValues[2]}
                initialRegion={'us-east-1'}
                handleInputChange={jest.fn()}
                handleSelectRegion={handleSelectRegionMock}
                regions={regionsMock}
                methods={useFormMock.useForm()}
            />
        );
        screen.getByText('us-east-1');
        for (let i = 0; i < regionsMock.length; i++) {
            const regionEl = screen.getByText(regionsMock[i]);
            expect(regionEl).toBeInTheDocument();
        }
        fireEvent.change(screen.getByTestId('region-select'), { target: { value: 'us-east-2' } });
        expect(handleSelectRegionMock).toBeCalled();
        expect(handleSelectRegionMock).toBeCalledWith('us-east-2');
    });
});
