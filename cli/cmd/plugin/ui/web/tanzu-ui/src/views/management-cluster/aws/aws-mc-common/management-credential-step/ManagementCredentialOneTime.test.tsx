// React imports
import React from 'react';

// Library imports
import { fireEvent, render, screen } from '@testing-library/react';
import ManagementCredentialOneTime from './ManagementCredentialOneTime';
import { FormProvider } from 'react-hook-form';

// App imports

// const regionsMock = [
//     'us-east-1',
//     'us-east-2',
//     'us-west-1',
//     'us-west-2',
//     'eu-central-1',
//     'eu-east-1',
//     'eu-east-2',
//     'ap-east-1',
//     'ap-south-1',
//     'ca-central-1',
// ];
// const testInputValues = ['testSecretAccessKey', 'testSessionToken', 'testAccessId'];
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

describe('ManagementCredentialOneTime component', () => {
    it('should render', async () => {
        const mockSelectCallback = jest.fn();
        const view = render(
            <FormProvider {...methods}>
                <ManagementCredentialOneTime selectCallback={mockSelectCallback} />
            </FormProvider>
        );
        expect(view).toBeDefined();
    });
    it('should handle region change', () => {
        const mockSelectCallback = jest.fn();
        render(
            <FormProvider {...methods}>
                <ManagementCredentialOneTime selectCallback={mockSelectCallback} />
            </FormProvider>
        );
        const region = screen.getByTestId('region-select');
        fireEvent.change(region, { target: { value: 'us-east-2' } });
        expect(mockSelectCallback).toHaveBeenCalled();
    });
    //TODO: More test cases should be added such as check initial values and text input change.
});
