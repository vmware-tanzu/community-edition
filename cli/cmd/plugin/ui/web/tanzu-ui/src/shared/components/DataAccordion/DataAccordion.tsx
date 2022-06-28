import { CdsIcon } from '@cds/react/icon';
import React, { useReducer } from 'react';
import { DataAccordionConfig } from './DataAccordion.interface';
import './DataAccordion.scss';
import { accordionReducer, DataAccordionActions, initialState } from './DataAccordion.store';

function DataAccordionItem({ title, content }: { title: string; content: string }) {
    const [state, dispatch] = useReducer(accordionReducer, initialState);

    return (
        <li className={`data-accordion-item ${state.active ? 'active' : ''}`}>
            <button
                cds-layout="p:lg horizontal"
                className="data-accordion-item-header w-full"
                onClick={() => dispatch({ type: DataAccordionActions.ToggleAccordion })}
            >
                <p cds-text={state.active ? 'subsection bold' : 'subsection'} className="data-accordion-item-header-title">
                    {title}
                </p>
                <CdsIcon cds-layout="align:right" shape="angle" className={state.active ? 'angle-down' : 'angle-right'}></CdsIcon>
            </button>
            <div className={`data-accordion-item-content-wrapper ${state.active ? 'open' : ''}`}>
                <p cds-text="body" className="data-accordion-item-content" dangerouslySetInnerHTML={{ __html: content }}></p>
            </div>
        </li>
    );
}

function DataAccordion<T>({ config }: { config: DataAccordionConfig<T> }) {
    const { data, ...selectors } = config;

    return (
        <ul className="data-accordion" cds-layout="m:none p:none">
            {data.map((item: T) => (
                <DataAccordionItem key={selectors.key(item)} title={selectors.title(item)} content={selectors.content(item)} />
            ))}
        </ul>
    );
}

export default DataAccordion;
