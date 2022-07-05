import './ContextualHelp.scss';

import { ClarityIcons, popOutIcon, searchIcon, viewColumnsIcon } from '@cds/core/icon';
import { ContextualHelpContentProps, ContextualHelpData, SearchProps } from './ContextualHelp.interface';
import { DrawerActions, drawerReducer, initialState } from '../Drawer/Drawer.store';
import React, { useContext, useEffect, useMemo, useReducer } from 'react';

import { CdsButton } from '@cds/react/button';
import { CdsControlAction } from '@cds/react/forms';
import { CdsIcon } from '@cds/react/icon';
import { CdsInput } from '@cds/react/input';
import { ContextualHelpState } from './ContextualHelp.store';
import DataAccordion from '../DataAccordion/DataAccordion';
import { DataAccordionConfig } from '../DataAccordion/DataAccordion.interface';
import Drawer from '../Drawer/Drawer';
import Fuse from 'fuse.js';
import { STORE_SECTION_UI } from '../../../state-management/reducers/Ui.reducer';
import { Store } from '../../../state-management/stores/Store';
import fuseIndex from '../../../assets/contextualHelpDocs/fuse-index.json';
import helpDocs from '../../../assets/contextualHelpDocs/data.json';
import useFuseSearch from '../../hooks/UseFuseSearch';

ClarityIcons.addIcons(viewColumnsIcon, searchIcon, popOutIcon);

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

function ContextualHelpContent({ title, keywords, externalLink, onClose, togglePin, ...drawerState }: ContextualHelpContentProps) {
    const { hits, onSearch } = useFuseSearch<ContextualHelpData>(helpDocs.data, fuseIndex, {
        keys: ['topicIds', 'topicTitle'],
        includeScore: true,
        useExtendedSearch: true,
        initialQuery: `=${keywords.join(' ').trim()}`,
    });

    const dataAccordionConfig: DataAccordionConfig<Fuse.FuseResult<ContextualHelpData>> = useMemo(() => {
        return {
            data: hits,
            key: (item: Fuse.FuseResult<ContextualHelpData>) => item.refIndex,
            title: (item: Fuse.FuseResult<ContextualHelpData>) => item.item.topicTitle,
            content: (item: Fuse.FuseResult<ContextualHelpData>) => item.item.htmlContent,
        };
    }, [hits]);

    return (
        <Drawer
            direction={drawerState.direction}
            open={drawerState.open}
            pinned={drawerState.pinned}
            onClose={onClose}
            togglePin={togglePin}
        >
            <div className="container w-full h-full" cds-layout="vertical">
                <div cds-layout="horizontal align-horizontal:center p-x:xxl m-b:lg">
                    <Search onSearch={onSearch} />
                </div>

                <h3 cds-layout="m-l:xxl m-y:sm" cds-text="title">
                    {title.pageTitle}
                </h3>

                <div className="content h-full" cds-layout="p:lg vertical align:stretch">
                    <DataAccordion config={dataAccordionConfig} />
                </div>

                <footer cds-layout="vertical" className="w-full">
                    <div cds-layout="horizontal align:right p:sm">
                        <CdsButton
                            cds-layout="vertical align:vertical-center"
                            className="learn-more"
                            size="sm"
                            onClick={() => {
                                window.open(externalLink, '_blank');
                            }}
                        >
                            <CdsIcon shape="pop-out"></CdsIcon>
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

function ContextualHelp() {
    const { state } = useContext(Store);
    const [drawerState, drawerDispatch] = useReducer(drawerReducer, initialState);
    const contextualHelpState: ContextualHelpState = state[STORE_SECTION_UI].contextualHelp;

    useEffect(() => {
        const mainContainer: HTMLElement | null = document.getElementById('main');
        if (mainContainer) {
            mainContainer.style.marginRight = drawerState.pinned ? '20rem' : '0px';
        }
    }, [drawerState.pinned]);

    return (
        <div className="contextual-help-container" cds-layout="align:right align:vertical-center">
            <CdsButton
                aria-label="contextual-help-info"
                cds-layout="m-r:md"
                className="more-info"
                action="flat"
                size="sm"
                onClick={() => drawerDispatch({ type: DrawerActions.OpenDrawer })}
            >
                Help: {contextualHelpState.title.contextTitle}
                <div className="panel-slide-r-l-solid" cds-layout="p:xs"></div>
            </CdsButton>

            {drawerState.open && (
                <ContextualHelpContent
                    title={contextualHelpState.title}
                    keywords={contextualHelpState.keywords}
                    externalLink={contextualHelpState.externalLink}
                    {...drawerState}
                    onClose={() => drawerDispatch({ type: DrawerActions.CloseDrawer })}
                    togglePin={() => drawerDispatch({ type: DrawerActions.TogglePin })}
                />
            )}
        </div>
    );
}

export default ContextualHelp;
