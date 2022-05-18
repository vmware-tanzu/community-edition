// React imports
import React from 'react';
import { BrowserRouter } from 'react-router-dom';

// Library imports
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import '@testing-library/jest-dom';

// App imports
import HeaderBar from './HeaderBar';
import VMWLogo from '../../../assets/vmw-logo.svg';

// Inject jest mock function for useNavigate()
const mockedNavigator = jest.fn();
jest.mock('react-router-dom', () => ({
    ...(jest.requireActual('react-router-dom') as any),
    useNavigate: () => mockedNavigator,
}));

describe('HeaderBar', () => {
    test('should render', () => {
        const view = render(<HeaderBar />, { wrapper: BrowserRouter });
        expect(view).toBeDefined();
    });

    test('renders the VMware logo', () => {
        render(<HeaderBar />, { wrapper: BrowserRouter });
        const logo = screen.getByLabelText('header-logo');
        expect(logo.getAttribute('src')).toEqual(VMWLogo);
    });

    test('has correct title text (Tanzu)', () => {
        render(<HeaderBar />, { wrapper: BrowserRouter });
        const title = screen.getByLabelText('header-title');
        expect(title).toHaveTextContent('Tanzu');
    });

    test('should route to the Welcome screen', () => {
        render(<HeaderBar />, { wrapper: BrowserRouter });
        userEvent.click(screen.getByLabelText('navigate-to-welcome'));
        expect(mockedNavigator).toHaveBeenCalled();
    });
});
