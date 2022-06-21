import { ClarityIcons, searchIcon, viewColumnsIcon } from '@cds/core/icon';
import { CdsButton } from '@cds/react/button';
import { CdsControlAction } from '@cds/react/forms';
import { CdsIcon } from '@cds/react/icon';
import { CdsInput } from '@cds/react/input';
import Fuse from 'fuse.js';
import React, { useEffect, useMemo, useReducer } from 'react';
import helpDocs from '../../../assets/contextualHelpDocs/data.json';
import fuseIndex from '../../../assets/contextualHelpDocs/fuse-index.json';
import useFuseSearch from '../../hooks/UseFuseSearch';
import DataAccordion from '../DataAccordion/DataAccordion';
import { DataAccordionConfig } from '../DataAccordion/DataAccordion.interface';
import Drawer from '../Drawer/Drawer';
import { DrawerActions, drawerReducer, initialState } from '../Drawer/Drawer.store';
import { ContextualHelpContentProps, ContextualHelpData, ContextualHelpProps, SearchProps } from './ContextualHelp.interface';
import './ContextualHelp.scss';

ClarityIcons.addIcons(viewColumnsIcon, searchIcon);

function Search({ value = '', onSearch }: SearchProps) {
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
}

function ContextualHelpContent({ title, keywords, onClose, togglePin, ...state }: ContextualHelpContentProps) {
    const { hits, onSearch } = useFuseSearch<ContextualHelpData>(helpDocs.data, fuseIndex, {
        keys: ['topicIds', 'topicTitle'],
        includeScore: true,
        defaultQuery: `=${keywords.join(' ').trim()}`,
    });

    const dataAccordionConfig: DataAccordionConfig<Fuse.FuseResult<ContextualHelpData>> = useMemo(() => {
        return {
            data: hits,
            key: (item: Fuse.FuseResult<ContextualHelpData>) => item.refIndex,
            title: (item: Fuse.FuseResult<ContextualHelpData>) => item.item.topicTitle,
            content: (item: Fuse.FuseResult<ContextualHelpData>) => item.item.htmlContent,
        };
    }, [hits]);

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
                    <DataAccordion config={dataAccordionConfig} />
                </div>

                <footer cds-layout="vertical" className="w-full">
                    <div cds-layout="horizontal align:right p:sm">
                        <CdsButton
                            className="learn-more"
                            size="sm"
                            onClick={() => {
                                window.open('http://tanzucommunityedition.io', '_blank');
                            }}
                        >
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
}

function ContextualHelp({ title, keywords }: ContextualHelpProps) {
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
}

export default ContextualHelp;
