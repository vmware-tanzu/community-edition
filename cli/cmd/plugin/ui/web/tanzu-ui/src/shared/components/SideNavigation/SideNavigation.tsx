// React imports
import React, { useContext } from 'react';
import { Link } from 'react-router-dom';

// Library imports
import { CdsIcon } from '@cds/react/icon';
import {
    chatBubbleIcon,
    ClarityIcons,
    clusterIcon,
    compassIcon,
    computerIcon,
    deployIcon,
    homeIcon,
    listIcon,
    nodeIcon,
    nodesIcon,
} from '@cds/core/icon';
import { CdsNavigation, CdsNavigationItem, CdsNavigationStart } from '@cds/react/navigation';
import { CdsDivider } from '@cds/react/divider';

// App imports
import { AppFeature, featureAvailable } from '../../services/AppConfiguration.service';
import { NavRoutes } from '../../constants/NavRoutes.constants';
import { Store } from '../../../state-management/stores/Store';
import { STORE_SECTION_UI } from '../../../state-management/reducers/Ui.reducer';
import { STORE_SECTION_APP } from '../../../state-management/reducers/App.reducer';
import { TOGGLE_NAV } from '../../../state-management/actions/Ui.actions';

ClarityIcons.addIcons(clusterIcon, chatBubbleIcon, compassIcon, computerIcon, deployIcon, homeIcon, listIcon, nodeIcon, nodesIcon);

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
        return route === state[STORE_SECTION_APP].appRoute;
    };

    const workloadClusterSupport = featureAvailable(AppFeature.WORKLOAD_CLUSTER_SUPPORT);
    const unmanagedClusterSupport = featureAvailable(AppFeature.UNMANAGED_CLUSTER_SUPPORT);
    const navigationItemLinkLayout = state[STORE_SECTION_UI].navExpanded ? 'p-l:md' : 'p-l:xs';

    return (
        <CdsNavigation expanded={state[STORE_SECTION_UI].navExpanded} onExpandedChange={onNavExpandedChange}>
            <CdsNavigationStart></CdsNavigationStart>
            <CdsNavigationItem cds-layout="m-t:sm" active={isActiveNavItem(NavRoutes.WELCOME)}>
                <Link cds-layout={navigationItemLinkLayout} to={NavRoutes.WELCOME}>
                    <CdsIcon shape="home" size="sm"></CdsIcon>
                    Welcome
                </Link>
            </CdsNavigationItem>
            <CdsNavigationItem active={isActiveNavItem(NavRoutes.GETTING_STARTED)}>
                <Link cds-layout={navigationItemLinkLayout} to={NavRoutes.GETTING_STARTED}>
                    <CdsIcon shape="compass" size="sm"></CdsIcon>
                    Getting started
                </Link>
            </CdsNavigationItem>
            <CdsDivider cds-layout="p-y:sm"></CdsDivider>
            <CdsNavigationItem active={isActiveNavItem(NavRoutes.MANAGEMENT_CLUSTER_INVENTORY)}>
                <Link cds-layout={navigationItemLinkLayout} to={NavRoutes.MANAGEMENT_CLUSTER_INVENTORY}>
                    <CdsIcon shape="blocks-group" size="sm"></CdsIcon>
                    Management Clusters
                </Link>
            </CdsNavigationItem>
            {workloadClusterSupport && (
                <CdsNavigationItem active={isActiveNavItem(NavRoutes.WORKLOAD_CLUSTER_INVENTORY)}>
                    <Link cds-layout={navigationItemLinkLayout} to={NavRoutes.WORKLOAD_CLUSTER_INVENTORY}>
                        <CdsIcon shape="nodes" size="sm"></CdsIcon>
                        Workload Clusters
                    </Link>
                </CdsNavigationItem>
            )}
            {unmanagedClusterSupport && (
                <CdsNavigationItem active={isActiveNavItem(NavRoutes.UNMANAGED_CLUSTER_INVENTORY)}>
                    <Link cds-layout={navigationItemLinkLayout} to={NavRoutes.UNMANAGED_CLUSTER_INVENTORY}>
                        <CdsIcon shape="computer" size="sm"></CdsIcon>
                        Unmanaged Clusters
                    </Link>
                </CdsNavigationItem>
            )}
            <CdsDivider cds-layout="p-y:sm"></CdsDivider>
            {workloadClusterSupport && (
                <CdsNavigationItem active={isActiveNavItem(NavRoutes.WORKLOAD_CLUSTER_WIZARD)}>
                    <Link cds-layout={navigationItemLinkLayout} to={NavRoutes.WORKLOAD_CLUSTER_WIZARD}>
                        <CdsIcon shape="deploy" size="sm"></CdsIcon>
                        WC Wizard - temp
                    </Link>
                </CdsNavigationItem>
            )}
            <CdsNavigationItem>
                <Link cds-layout={navigationItemLinkLayout} to="/progress">
                    <CdsIcon shape="deploy" size="sm"></CdsIcon>
                    Deploy progress - temp
                </Link>
            </CdsNavigationItem>
            <CdsDivider cds-layout="p-y:sm"></CdsDivider>
            {/* TODO: Determine links for external pages */}
            <CdsNavigationItem>
                <Link cds-layout={navigationItemLinkLayout} to={'/'}>
                    <CdsIcon shape="list" size="sm"></CdsIcon>
                    FAQ
                </Link>
            </CdsNavigationItem>
            <CdsNavigationItem>
                <Link cds-layout={navigationItemLinkLayout} to={'/'}>
                    <CdsIcon shape="chat-bubble" size="sm"></CdsIcon>
                    Feedback
                </Link>
            </CdsNavigationItem>
            <CdsNavigationItem>
                <Link cds-layout={navigationItemLinkLayout} to={'/'}>
                    <CdsIcon shape="computer" size="sm"></CdsIcon>
                    Contribute
                </Link>
            </CdsNavigationItem>
        </CdsNavigation>

        // TODO: Add list of recent pages in order - create a GH issue in backlog
    );
}

export default SideNavigation;
