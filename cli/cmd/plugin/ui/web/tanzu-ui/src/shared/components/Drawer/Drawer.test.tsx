import { fireEvent, render, screen } from '@testing-library/react';
import React from 'react';
import Drawer from './Drawer';
import { Direction } from './Drawer.enum';

describe('Drawer Component', () => {
    test('Close and Toggle Pin Drawer', () => {
        const onClose = jest.fn();
        const togglePin = jest.fn();

        render(
            <Drawer direction={Direction.Right} open={true} pinned={false} onClose={onClose} togglePin={togglePin}>
                <div>Drawer content body</div>
            </Drawer>
        );

        expect(onClose).not.toHaveBeenCalled();
        expect(togglePin).not.toHaveBeenCalled();

        const closeButton = screen.getByLabelText('drawer-close');
        fireEvent.click(closeButton);

        expect(onClose).toHaveBeenCalled();

        const togglePinButton = screen.getByLabelText('drawer-toggle-pin');
        fireEvent.click(togglePinButton);

        expect(togglePin).toHaveBeenCalled();

        const contentBody = screen.getByText('Drawer content body');
        expect(contentBody).toBeInTheDocument();
    });
});
