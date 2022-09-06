// React imports
import React, { useEffect } from 'react';

// Library imports
import { CdsIcon } from '@cds/react/icon';
import { ClarityIcons, moonIcon, sunIcon } from '@cds/core/icon';

// App imports
import { AppThemes, useThemeService } from '../../services/Themes.service';
import './ThemeToggle.component.scss';

ClarityIcons.addIcons(moonIcon, sunIcon);

function ThemeToggle() {
    const themeService = useThemeService();

    // TODO: refactor custom hook so useEffect not needed
    useEffect(() => {
        themeService.initializeTheme();
    }, []); // eslint-disable-line react-hooks/exhaustive-deps

    function displayAlternatingThemeSwitch() {
        const theme = themeService.getTheme();
        const icon = theme === AppThemes.LIGHT ? 'moon' : 'sun';
        const text = theme === AppThemes.LIGHT ? 'Dark' : 'Light';

        return (
            <>
                <CdsIcon shape={icon} size="sm" cds-layout="m-l:md"></CdsIcon>
                <span cds-layout="p-x:xs">{text}</span>
            </>
        );
    }

    return (
        <>
            <span onClick={themeService.toggleTheme} aria-label="toggle color theme">
                {displayAlternatingThemeSwitch()}
            </span>
        </>
    );
}

export default ThemeToggle;
