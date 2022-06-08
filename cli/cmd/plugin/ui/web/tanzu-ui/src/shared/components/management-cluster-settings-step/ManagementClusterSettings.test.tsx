// React imports
import React from 'react';

// Library imports
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

// App imports
import ManagementClusterSettings from './ManagementClusterSettings';

describe('ManagementClusterSettings', () => {
    it('should render', async () => {
        const view = render(<ManagementClusterSettings />);
        expect(view).toBeDefined();
    });

    it('show a message if one option is provided', () => {
        render(<ManagementClusterSettings message="test message" />);
        expect(screen.getByText('test message')).toBeInTheDocument();
    });

    it('show three options for the control plane node profile', () => {
        render(<ManagementClusterSettings />);
        const labels = ['Single node', 'High availability', 'Production-ready (High availability)'];
        labels.forEach((label) => {
            expect(screen.getByText(label)).toBeInTheDocument();
        });
        expect(screen.getAllByRole('radio')[0].hasAttribute('checked')).toBe(true);
    });

    it('should update cluster name', () => {
        const handleValueChangeMock = jest.fn();
        render(<ManagementClusterSettings handleValueChange={handleValueChangeMock} />);
        const input = screen.getByLabelText('cluster name');
        fireEvent.change(input, { target: { value: 'myTestCluster' } });
        expect(handleValueChangeMock).toHaveBeenCalled();
    });

    it('select the third option in the control plane profile', async () => {
        const user = userEvent.setup();
        render(<ManagementClusterSettings />);
        const thirdRadio = screen.getAllByRole('radio')[2];
        await user.click(thirdRadio);
        expect(screen.getAllByTestId('cds-radio')[2]).toHaveAttribute('_checked');
    });

    it('create management cluster button can be triggered', () => {
        const deployMock = jest.fn();
        render(<ManagementClusterSettings deploy={deployMock} />);
        fireEvent.click(screen.getByText('Create Management cluster'));
        expect(deployMock.call.length).toBe(1);
    });
});
