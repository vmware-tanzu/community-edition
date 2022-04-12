// React imports
import React from 'react';
import { BrowserRouter, MemoryRouter } from 'react-router-dom';
import { createMemoryHistory } from 'history';

// Library imports
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import '@testing-library/jest-dom';

// App imports
import HeaderBar from './HeaderBar';
import VMWLogo from '../../../assets/vmw-logo.svg';

describe('HeaderBar', () => {
    test('should render', () => {
        const view = render(<BrowserRouter><HeaderBar /></BrowserRouter>);
        expect(view).toBeDefined();
    });

    test("renders the VMware logo", () => {
        render(<BrowserRouter><HeaderBar /></BrowserRouter>);
        const logo = screen.getByLabelText('header-logo');
        expect(logo.getAttribute('src')).toEqual(VMWLogo);
    });

    test('has correct title text (Tanzu)', () => {
        render(<BrowserRouter><HeaderBar /></BrowserRouter>);
        const title = screen.getByLabelText('header-title');
        expect(title).toHaveTextContent('Tanzu');
    });

    test('should route to the Welcome screen', () => {
        // const history = createMemoryHistory({ initialEntries: ['/ui/getting-started'] });
        // const route = '/some-route'
        const navigateHome = jest.fn(() => {});
        render(<BrowserRouter><HeaderBar /></BrowserRouter>);
        // const title = screen.getByLabelText('header-title');
        // expect(title).toHaveTextContent('Tanzu');
        const nav = screen.getByLabelText('header-logo');

        // expect(history.location.pathname).toBe('/ui/getting-started');
        fireEvent.click(nav);

        expect(navigateHome).toHaveBeenCalled();

        // history.push('/ui');
        // expect(history.location.pathname).toBe('/ui');
    });
});

//
//
// test('should pass', () => {
//     const history = createMemoryHistory({ initialEntries: ['/ui/getting-started'] });
//     const navigateHome = jest.fn(() => {
//         console.log('hi');
//     });
//     render(<BrowserRouter><HeaderBar /></BrowserRouter>);
//     const instance = screen.instance();
//     const logo = screen.getByLabelText('header-logo');
//     expect(history.location.pathname).toBe('/ui/getting-started');
//
//
//     fireEvent(screen.getByRole('button'), new MouseEvent('click', {
//         bubbles: true,
//         cancelable: true,
//     }));
//     console.log(history.location);
//     expect(navigateHome).toHaveBeenCalled();
//     // expect(history.location.pathname).toBe('/ui/getting-started');
// });

// test('should route to the Welcome screen', () => {
//     render(<BrowserRouter><HeaderBar /></BrowserRouter>);
//     const link = screen.getByLabelText('navigation-link');
//     expect(link).toHaveAttribute('to');
// });

// test('should route to the Welcome screen', () => {
//     const history = createMemoryHistory({ initialEntries: ['/ui/getting-started'] });
//     // const route = '/some-route'
//     history.push = jest.fn(() => {
//         // history.push('/ui');
//     })
//     render(<BrowserRouter><HeaderBar /></BrowserRouter>);
//     // const title = screen.getByLabelText('header-title');
//     // expect(title).toHaveTextContent('Tanzu');
//     const nav = screen.getByLabelText('navigation-link');
//     expect(history.location.pathname).toBe('/ui/getting-started');
//     fireEvent(nav, new MouseEvent('click', {
//         bubbles: true,
//         cancelable: true,
//     }));
//     expect(history.push).toHaveBeenCalledWith('/ui');
//
//     // history.push('/ui');
//     // expect(history.location.pathname).toBe('/ui');
// });