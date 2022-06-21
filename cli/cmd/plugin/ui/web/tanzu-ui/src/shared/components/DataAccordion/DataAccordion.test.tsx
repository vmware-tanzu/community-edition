import { render, screen } from '@testing-library/react';
import React from 'react';
import DataAccordion from './DataAccordion';
import { DataAccordionItemData } from './DataAccordion.interface';

describe('DataAccordion Component', () => {
    const managementClusterData: DataAccordionItemData[] = [
        {
            id: 1,
            title: 'What is a management cluster',
            content: `Tenetur ullam rerum ad iusto possimus sequi mollitia dolore sunt quam praesentium. 
            Tenetur ullam rerum ad iusto possimus sequi mollitia dolore sunt quam praesentium.Tenetur ullam 
            rerum ad iusto possimus sequi mollitia dolore sunt quam praesentium.`,
        },
    ];
    test('Populate Data Accordion Items', () => {
        render(<DataAccordion accordionData={managementClusterData} />);
        expect(screen.getByText('What is a management cluster')).toBeInTheDocument();
    });
});
