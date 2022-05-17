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
                <Link to="/">
                    <CdsIcon shape="home" size="sm"></CdsIcon>
                    Welcome
                </Link>
            </CdsNavigationItem>
            <CdsNavigationItem>
                <Link to="/getting-started">
                    <CdsIcon shape="compass" size="sm"></CdsIcon>
                    Getting Started
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