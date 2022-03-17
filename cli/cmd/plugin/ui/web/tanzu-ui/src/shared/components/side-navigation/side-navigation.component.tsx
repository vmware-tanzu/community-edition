// React imports
import React, { useContext } from 'react';
import { Link } from 'react-router-dom';

// Library imports
import { CdsIcon } from '@cds/react/icon';
import { ClarityIcons, homeIcon, compassIcon } from '@cds/core/icon';
import { CdsNavigation, CdsNavigationItem, CdsNavigationStart } from '@cds/react/navigation';

// App imports
import { Store } from '../../../state-management/stores/store';
import { TOGGLE_NAV } from '../../../state-management/actions/actionTypes';

ClarityIcons.addIcons(homeIcon, compassIcon);

function SideNavigationComponent(this: any) {

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
        </CdsNavigation>
    );
}

export default SideNavigationComponent;