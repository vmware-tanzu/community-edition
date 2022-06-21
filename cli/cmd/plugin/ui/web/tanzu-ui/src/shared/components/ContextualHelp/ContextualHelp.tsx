import { ClarityIcons, searchIcon, viewColumnsIcon } from '@cds/core/icon';
import { CdsButton } from '@cds/react/button';
import { CdsControlAction } from '@cds/react/forms';
import { CdsIcon } from '@cds/react/icon';
import { CdsInput } from '@cds/react/input';
import React, { useEffect, useReducer } from 'react';
import helpDocs from '../../../assets/contextualHelpDocs/data.json';
import fuseIndex from '../../../assets/contextualHelpDocs/fuse-index.json';
import useFuseSearch from '../../hooks/UseFuseSearch';
import DataAccordion from '../DataAccordion/DataAccordion';
import Drawer from '../Drawer/Drawer';
import { DrawerActions, drawerReducer, initialState } from '../Drawer/Drawer.store';
import { ContextualHelpContentProps, ContextualHelpData, ContextualHelpProps, SearchProps } from './ContextualHelp.interface';
import './ContextualHelp.scss';
import * as ContextualHelpUtils from './ContextualHelp.utils';

ClarityIcons.addIcons(viewColumnsIcon, searchIcon);

const Search: React.FunctionComponent<SearchProps> = ({ value = '', onSearch }) => {
    const [searchValue, setSearchValue] = React.useState(value);

    const onChange = (value: string) => {
        setSearchValue(value);
        onSearch(value);
    };

    return (
        <CdsInput>
            <label>Search</label>
            <input
                type="text"
                value={searchValue}
                aria-label="contextual-help-search"
                className="search"
                placeholder="Search"
                onChange={(e) => onChange(e.target.value)}
            ></input>
            <CdsControlAction action="prefix" shape="search" aria-label="search"></CdsControlAction>
            {searchValue && (
                <CdsControlAction action="suffix" shape="close" aria-label="clear-search" onClick={() => onChange('')}></CdsControlAction>
            )}
        </CdsInput>
    );
};

const ContextualHelpContent: React.FunctionComponent<ContextualHelpContentProps> = ({ title, keywords, onClose, togglePin, ...state }) => {
    const { hits, onSearch } = useFuseSearch<ContextualHelpData>(helpDocs.data, fuseIndex, {
        keys: ['topicIds', 'topicTitle'],
        includeScore: true,
        defaultQuery: `=${keywords.join(' ').trim()}`,
    });

    useEffect(() => {
        const mainContainer: HTMLElement | null = document.getElementById('main');
        if (mainContainer) {
            mainContainer.style.marginRight = state.pinned ? '20rem' : '0px';
        }
    }, [state.pinned]);

    return (
        <Drawer direction={state.direction} open={state.open} pinned={state.pinned} onClose={onClose} togglePin={togglePin}>
            <div className="container w-full h-full" cds-layout="vertical">
                <h3 cds-layout="m-l:xxl m-y:sm" cds-text="title">
                    {title}
                </h3>

                <div cds-layout="horizontal align-horizontal:center p:xxl">
                    <Search onSearch={onSearch} />
                </div>

                <div className="content h-full" cds-layout="p:lg vertical align:stretch">
                    <DataAccordion accordionData={ContextualHelpUtils.generateDataAccordionContent(hits)} />
                </div>

                <footer cds-layout="vertical" className="w-full">
                    <div cds-layout="horizontal align:right p:sm">
                        <CdsButton className="learn-more" size="sm">
                            Learn More
                        </CdsButton>
                    </div>

                    <div className="footer-image" cds-layout="vertical align:stretch">
                        <div className="contextual-help-img"></div>
                    </div>
                </footer>
            </div>
        </Drawer>
    );
};

const ContextualHelp: React.FunctionComponent<ContextualHelpProps> = ({ title, keywords }) => {
    const [state, dispatch] = useReducer(drawerReducer, initialState);

    return (
        <div className="contextual-help-container" cds-layout="align:right align:vertical-center">
            <CdsButton
                aria-label="contextual-help-info"
                cds-layout="m-r:md"
                className="more-info"
                action="flat"
                onClick={() => dispatch({ type: DrawerActions.OpenDrawer })}
            >
                More Info
                <CdsIcon shape="view-columns"> </CdsIcon>
            </CdsButton>

            {state.open && (
                <ContextualHelpContent
                    title={title}
                    keywords={keywords}
                    {...state}
                    onClose={() => dispatch({ type: DrawerActions.CloseDrawer })}
                    togglePin={() => dispatch({ type: DrawerActions.TogglePin })}
                />
            )}
        </div>
    );
};

export default ContextualHelp;
