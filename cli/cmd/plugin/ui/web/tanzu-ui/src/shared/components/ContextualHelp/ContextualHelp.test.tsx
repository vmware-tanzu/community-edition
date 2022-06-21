import { fireEvent, render, screen } from '@testing-library/react';
import React from 'react';
import ContextualHelp from './ContextualHelp';

describe('ContextualHelp Component', () => {
    test('Open Drawer', () => {
        render(<ContextualHelp keywords={['tce-docker']} title={'Management cluster'} />);
        const contextualHelpInfo = screen.getByLabelText('contextual-help-info');

        fireEvent.click(contextualHelpInfo);

        expect(screen.getByText('Management cluster')).toBeInTheDocument();
        expect(screen.getByText('Docker Daemon')).toBeInTheDocument();
    });
});
