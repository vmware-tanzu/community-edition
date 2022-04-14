// React imports
import React, { useContext } from 'react';
import { Link } from 'react-router-dom';

// Library imports
import { CdsIcon } from '@cds/react/icon';
import { ClarityIcons, homeIcon, compassIcon, deployIcon } from '@cds/core/icon';
import { CdsNavigation, CdsNavigationItem, CdsNavigationStart } from '@cds/react/navigation';

// App imports
import { Store } from '../../../state-management/stores/Store';
import { TOGGLE_NAV } from '../../../state-management/actions/Ui.actions';
import { NavRoutes } from '../../constants/NavRoutes.constants';

ClarityIcons.addIcons(homeIcon, compassIcon, deployIcon);

function SideNavigation(this: any) {

    const { state, dispatch } = useContext(Store);

    const onNavExpandedChange = () => {
        dispatch({
            type: TOGGLE_NAV
        });
    };

    return (
        <CdsNavigation expanded={state.ui.navExpanded} onExpandedChange={onNavExpandedChange}>
            <CdsNavigationStart></CdsNavigationStart>
            <CdsNavigationItem>
                <Link to={NavRoutes.WELCOME}>
                    <CdsIcon shape="home" size="sm"></CdsIcon>
                    Welcome
                </Link>
            </CdsNavigationItem>
            <CdsNavigationItem>
                <Link to={NavRoutes.GETTING_STARTED}>
                    <CdsIcon shape="compass" size="sm"></CdsIcon>
                    Getting Started
                </Link>
            </CdsNavigationItem>
            <CdsNavigationItem>
                <Link to={NavRoutes.MANAGEMENT_CLUSTER_LANDING}>
                    <CdsIcon shape="deploy" size="sm"></CdsIcon>
                    Management Cluster
                </Link>
            </CdsNavigationItem>
            <CdsNavigationItem>
                <Link to={NavRoutes.WORKLOAD_CLUSTER_LANDING}>
                    <CdsIcon shape="deploy" size="sm"></CdsIcon>
                    Workload Cluster
                </Link>
            </CdsNavigationItem>
            <CdsNavigationItem>
                <Link to="/progress">
                    <CdsIcon shape="deploy" size="sm"></CdsIcon>
                    Stream Logs - Temp
                </Link>
            </CdsNavigationItem>
        </CdsNavigation>
    );
}

export default SideNavigation;
