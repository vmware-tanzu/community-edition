import { render, screen } from '@testing-library/react';
import Fuse from 'fuse.js';
import React from 'react';
import { ContextualHelpData } from '../ContextualHelp/ContextualHelp.interface';
import DataAccordion from './DataAccordion';
import { DataAccordionConfig } from './DataAccordion.interface';

describe('DataAccordion Component', () => {
    const managementClusterData: Fuse.FuseResult<ContextualHelpData>[] = [
        {
            refIndex: 1,
            item: {
                topicTitle: 'Docker Daemon',
                htmlContent: '<p>Complete the following steps for a full Tanzu Community Edition implementation to AWS:</p>',
                topicIds: ['demo'],
            },
        },
    ];

    const dataAccordionConfig: DataAccordionConfig<Fuse.FuseResult<ContextualHelpData>> = {
        data: managementClusterData,
        key: (item: Fuse.FuseResult<ContextualHelpData>) => item.refIndex,
        title: (item: Fuse.FuseResult<ContextualHelpData>) => item.item.topicTitle,
        content: (item: Fuse.FuseResult<ContextualHelpData>) => item.item.htmlContent,
    };

    test('Populate Data Accordion Items', () => {
        render(<DataAccordion config={dataAccordionConfig} />);
        expect(screen.getByText('Docker Daemon')).toBeInTheDocument();
    });
});
