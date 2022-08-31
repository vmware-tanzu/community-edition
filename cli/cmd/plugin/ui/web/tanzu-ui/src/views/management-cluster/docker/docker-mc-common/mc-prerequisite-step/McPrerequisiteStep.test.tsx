// React imports
import React from 'react';

// Library imports
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { rest } from 'msw';
import { setupServer } from 'msw/node';
import '@testing-library/jest-dom';

// App imports
import McPrerequisiteStep from './McPrerequisiteStep';
import { STATUS } from '../../../../../shared/constants/App.constants';

describe('McPrerequisiteStep', () => {
    const server = setupServer(
        rest.get('/api/containerruntime', (req, res, ctx) => {
            return res(ctx.status(200));
        })
    );

    beforeAll(() => server.listen());
    afterEach(() => server.resetHandlers());
    afterAll(() => server.close());

    test('should render', async () => {
        const view = render(<McPrerequisiteStep />);
        await waitFor(() => {
            expect(view).toBeDefined();
        });
    });

    test('show successful message if docker daemon is connected:', async () => {
        render(<McPrerequisiteStep />);
        expect(await screen.findByText('Docker daemon is running')).toBeInTheDocument();
    });

    test('show an error message if docker daemon cannot be connected:', async () => {
        server.use(
            rest.get('/api/containerruntime', (req, res, ctx) => {
                return res(ctx.status(500), ctx.json({ message: 'Unable to verify Docker daemon' }));
            })
        );
        render(<McPrerequisiteStep />);
        expect(await screen.findByText('Unable to verify Docker daemon: Unable to verify Docker daemon')).toBeInTheDocument();
    });

    test('connect button', async () => {
        server.use(
            rest.get('/api/containerruntime', (req, res, ctx) => {
                return res(ctx.status(500), ctx.json({ message: 'Unable to verify Docker daemon' }));
            })
        );
        render(<McPrerequisiteStep />);
        const connectBtn = await screen.findByText('VERIFY DOCKER DAEMON');
        server.use(
            rest.get('/api/containerruntime', (req, res, ctx) => {
                return res(ctx.status(200));
            })
        );
        fireEvent(connectBtn, new MouseEvent('click'));
        expect(await screen.findByText('Docker daemon is running')).toBeInTheDocument();
    });

    test('next button', async () => {
        server.use(
            rest.get('/api/containerruntime', (req, res, ctx) => {
                return res(ctx.status(200));
            })
        );
        const mockProps = {
            tabStatus: [STATUS.INVALID],
            currentStep: 1,
            submitForm: jest.fn(),
            goToStep: jest.fn(),
        };

        render(<McPrerequisiteStep {...mockProps} />);
        await screen.findByText('VERIFY DOCKER DAEMON');
        const nextBtn = await screen.findByText('NEXT');
        await screen.findByText('Docker daemon is running');
        fireEvent.click(nextBtn);
        expect(mockProps.submitForm).toHaveBeenCalled();
        expect(mockProps.goToStep).toHaveBeenCalled();
    });
});
