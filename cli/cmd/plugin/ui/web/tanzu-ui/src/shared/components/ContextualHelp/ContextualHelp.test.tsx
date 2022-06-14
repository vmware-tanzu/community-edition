import { fireEvent, render, screen } from '@testing-library/react';
import React from 'react';
import ContextualHelp from './ContextualHelp';

describe('ContextualHelp Component', () => {
    test('should render ContextualHelp Compononent', () => {
        render(<ContextualHelp />);
    });

    test('Open Drawer', () => {
        render(<ContextualHelp />);
        const contextualHelpInfo = screen.getByLabelText('contextual-help-info');

        fireEvent.click(contextualHelpInfo);
        expect(screen.getByText('Management clusters')).toBeInTheDocument();
    });
});
