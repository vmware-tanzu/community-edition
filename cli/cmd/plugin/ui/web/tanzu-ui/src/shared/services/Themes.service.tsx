// React imports
import { useContext } from 'react';

// App imports
import { SET_APP_THEME } from '../../state-management/actions/Ui.actions';
import { Store } from '../../state-management/stores/Store';
import { STORE_SECTION_UI } from '../../state-management/reducers/Ui.reducer';

export enum AppThemes {
    LIGHT = 'light',
    DARK = 'dark',
}

const useThemes = () => {
    const { state, dispatch } = useContext(Store);
    const localStorageThemeVar = 'cds-theme';

    const getStoredTheme = () => {
        return localStorage.getItem(localStorageThemeVar);
    };

    const setStoredTheme = (theme: AppThemes): void => {
        localStorage.setItem(localStorageThemeVar, theme);
    };

    const getDefaultTheme = (): AppThemes => {
        let theme = AppThemes.LIGHT;
        try {
            const stored = getStoredTheme();
            if (stored) {
                theme = stored as AppThemes;
            } else if (window.matchMedia('(prefers-color-scheme)').media !== 'not all') {
                theme = window.matchMedia('(prefers-color-scheme: light)').matches ? AppThemes.LIGHT : AppThemes.DARK;
                setStoredTheme(theme);
            }
        } catch (err) {
            console.log(`Error retrieving ${localStorageThemeVar} from local storage: ${err}`);
        }

        return theme;
    };

    const getTheme = (): AppThemes => {
        return state[STORE_SECTION_UI].theme;
    };

    const initializeTheme = (): void => {
        const theme = getDefaultTheme();
        dispatch({
            type: SET_APP_THEME,
            payload: theme,
        });
        applyTheme(theme);
    };

    const toggleTheme = (): void => {
        const oldTheme = state[STORE_SECTION_UI].theme || AppThemes.LIGHT;
        const theme = oldTheme === AppThemes.LIGHT ? AppThemes.DARK : AppThemes.LIGHT;

        dispatch({
            type: SET_APP_THEME,
            payload: theme,
        });
        applyTheme(theme);
        setStoredTheme(theme);
    };

    const applyTheme = (theme: AppThemes): void => {
        document.body.setAttribute('cds-theme', theme);
        document.body.setAttribute('class', theme);
    };

    return {
        getTheme,
        initializeTheme,
        toggleTheme,
    };
};

// Exports theme service
export const useThemeService = (): any => {
    return useThemes();
};
