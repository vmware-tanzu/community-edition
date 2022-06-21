import { CdsIcon } from '@cds/react/icon';
import React, { useReducer } from 'react';
import { DataAccordionActions, accordionReducer, initialState } from './DataAccordion.store';
import './DataAccordion.scss';
import { DataAccordionItemData } from './DataAccordion.interface';

const DataAccordionItem: React.FunctionComponent<{ item: DataAccordionItemData }> = ({ item }) => {
    const [state, dispatch] = useReducer(accordionReducer, initialState);
    const { title, content } = item;

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
};

const DataAccordion: React.FunctionComponent<{ accordionData: DataAccordionItemData[] }> = ({ accordionData }) => {
    return (
        <ul className="data-accordion" cds-layout="m:none p:none">
            {accordionData.map((item) => (
                <DataAccordionItem key={item.id} item={item} />
            ))}
        </ul>
    );
};

export default DataAccordion;
