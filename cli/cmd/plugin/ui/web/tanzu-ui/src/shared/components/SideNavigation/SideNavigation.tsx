// React imports
import React, { useContext } from 'react';
import { Link } from 'react-router-dom';

// Library imports
import { CdsIcon } from '@cds/react/icon';
import { ClarityIcons, homeIcon, compassIcon, deployIcon, clusterIcon, nodeIcon, nodesIcon } from '@cds/core/icon';
import { CdsNavigation, CdsNavigationItem, CdsNavigationStart } from '@cds/react/navigation';
import { CdsDivider } from '@cds/react/divider';

// App imports
import { Store } from '../../../state-management/stores/Store';
import { TOGGLE_NAV } from '../../../state-management/actions/Ui.actions';
import { NavRoutes } from '../../constants/NavRoutes.constants';

ClarityIcons.addIcons(clusterIcon, compassIcon, deployIcon, nodeIcon, nodesIcon, homeIcon);

function SideNavigation(this: any) {
    const { state, dispatch } = useContext(Store);

    // toggles navigation panel open/closed
    const onNavExpandedChange = () => {
        dispatch({
            type: TOGGLE_NAV,
        });
    };

    // helper function to determine if nav item is active
    const isActiveNavItem = (route: string) => {
        return route === state.app.appRoute;
    };

    return (
        <CdsNavigation expanded={state.ui.navExpanded} onExpandedChange={onNavExpandedChange}>
            <CdsNavigationStart></CdsNavigationStart>
            <CdsNavigationItem active={isActiveNavItem(NavRoutes.WELCOME)}>
                <Link to={NavRoutes.WELCOME}>
                    <CdsIcon shape="home" size="sm"></CdsIcon>
                    Welcome
                </Link>
            </CdsNavigationItem>
            <CdsNavigationItem active={isActiveNavItem(NavRoutes.GETTING_STARTED)}>
                <Link to={NavRoutes.GETTING_STARTED}>
                    <CdsIcon shape="compass" size="sm"></CdsIcon>
                    Getting Started
                </Link>
            </CdsNavigationItem>
            <CdsDivider></CdsDivider>
            <CdsNavigationItem active={isActiveNavItem(NavRoutes.MANAGEMENT_CLUSTER_INVENTORY)}>
                <Link to={NavRoutes.MANAGEMENT_CLUSTER_INVENTORY}>
                    <CdsIcon shape="cluster" size="sm"></CdsIcon>
                    Management Clusters
                </Link>
            </CdsNavigationItem>
            <CdsNavigationItem active={isActiveNavItem(NavRoutes.WORKLOAD_CLUSTER_INVENTORY)}>
                <Link to={NavRoutes.WORKLOAD_CLUSTER_INVENTORY}>
                    <CdsIcon shape="nodes" size="sm"></CdsIcon>
                    Workload Clusters
                </Link>
            </CdsNavigationItem>
            <CdsNavigationItem active={isActiveNavItem(NavRoutes.UNMANAGED_CLUSTER_INVENTORY)}>
                <Link to={NavRoutes.UNMANAGED_CLUSTER_INVENTORY}>
                    <CdsIcon shape="node" size="sm"></CdsIcon>
                    Unmanaged Clusters
                </Link>
            </CdsNavigationItem>
            <CdsDivider></CdsDivider>
            <CdsNavigationItem active={isActiveNavItem(NavRoutes.WORKLOAD_CLUSTER_WIZARD)}>
                <Link to={NavRoutes.WORKLOAD_CLUSTER_WIZARD}>
                    <CdsIcon shape="deploy" size="sm"></CdsIcon>
                    WC Wizard - temp
                </Link>
            </CdsNavigationItem>
            <CdsNavigationItem>
                <Link to="/progress">
                    <CdsIcon shape="deploy" size="sm"></CdsIcon>
                    Deploy progress - temp
                </Link>
            </CdsNavigationItem>
        </CdsNavigation>

        // TODO: Add list of recent pages in order - create a GH issue in backlog
    );
}

export default SideNavigation;
