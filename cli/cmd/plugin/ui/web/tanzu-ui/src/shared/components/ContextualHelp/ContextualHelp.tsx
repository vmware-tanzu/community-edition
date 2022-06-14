import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import React, { useEffect, useReducer } from 'react';
import Drawer from '../Drawer/Drawer';
import { DrawerActions } from '../Drawer/Drawer.enum';
import { drawerReducer, initialState } from '../Drawer/Drawer.store';
import './ContextualHelp.scss';
import DataAccordion from '../DataAccordion/DataAccordion';
import { DataAccordionItemData } from '../DataAccordion/DataAccordion.interface';
import { ClarityIcons, viewColumnsIcon } from '@cds/core/icon';

ClarityIcons.addIcons(viewColumnsIcon);

const managementClusterData: DataAccordionItemData[] = [
    {
        id: 1,
        title: 'What is a management cluster',
        content: `Tenetur ullam rerum ad iusto possimus sequi mollitia dolore sunt quam praesentium. 
        Tenetur ullam rerum ad iusto possimus sequi mollitia dolore sunt quam praesentium.Tenetur ullam 
        rerum ad iusto possimus sequi mollitia dolore sunt quam praesentium.`,
    },
    {
        id: 2,
        title: 'Creation of management clusters',
        content: `Aperiam ab atque incidunt dolores ullam est, earum ipsa recusandae velit cumque. 
        Aperiam ab atque incidunt dolores ullam est, earum ipsa recusandae velit cumque.`,
    },
    {
        id: 3,
        title: 'What prerequistes do I need',
        content: 'Blanditiis aliquid adipisci quisquam reiciendis voluptates itaque.',
    },
    {
        id: 4,
        title: 'What is a management cluster',
        content: `Tenetur ullam rerum ad iusto possimus sequi mollitia dolore sunt quam praesentium. 
        Tenetur ullam rerum ad iusto possimus sequi mollitia dolore sunt quam praesentium.Tenetur ullam 
        rerum ad iusto possimus sequi mollitia dolore sunt quam praesentium.`,
    },
    {
        id: 5,
        title: 'Creation of management clusters',
        content: `Aperiam ab atque incidunt dolores ullam est, earum ipsa recusandae velit cumque. 
        Aperiam ab atque incidunt dolores ullam est, earum ipsa recusandae velit cumque.`,
    },
    {
        id: 6,
        title: 'What prerequistes do I need',
        content: 'Blanditiis aliquid adipisci quisquam reiciendis voluptates itaque.',
    },
    {
        id: 7,
        title: 'What is a management cluster',
        content: `Tenetur ullam rerum ad iusto possimus sequi mollitia dolore sunt quam praesentium. 
        Tenetur ullam rerum ad iusto possimus sequi mollitia dolore sunt quam praesentium.Tenetur ullam 
        rerum ad iusto possimus sequi mollitia dolore sunt quam praesentium.`,
    },
    {
        id: 8,
        title: 'Creation of management clusters',
        content: `Aperiam ab atque incidunt dolores ullam est, earum ipsa recusandae velit cumque. 
        Aperiam ab atque incidunt dolores ullam est, earum ipsa recusandae velit cumque.`,
    },
    {
        id: 9,
        title: 'What prerequistes do I need',
        content: 'Blanditiis aliquid adipisci quisquam reiciendis voluptates itaque.',
    },
];

const ContextualHelp = () => {
    const [state, dispatch] = useReducer(drawerReducer, initialState);

    useEffect(() => {
        const mainContainer: HTMLElement | null = document.getElementById('main');
        if (mainContainer) {
            mainContainer.style.marginRight = state.pinned ? '18rem' : '0px';
        }
    }, [state.pinned]);

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
                <Drawer
                    direction={state.direction}
                    open={state.open}
                    pinned={state.pinned}
                    onClose={() => dispatch({ type: DrawerActions.CloseDrawer })}
                    togglePin={() => dispatch({ type: DrawerActions.TogglePin })}
                >
                    <div className="container w-full h-full" cds-layout="vertical">
                        <h3 cds-layout="m-l:xxl m-y:sm" cds-text="title">
                            Management clusters
                        </h3>

                        <div className="content h-full" cds-layout="p:lg vertical align:stretch">
                            <DataAccordion data={managementClusterData} />
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
            )}
        </div>
    );
};

export default ContextualHelp;
