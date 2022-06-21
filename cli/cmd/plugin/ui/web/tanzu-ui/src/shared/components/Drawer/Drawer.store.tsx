import { Direction } from './Drawer.enum';

/**
 * State
 */
export interface DrawerState {
    open: boolean;
    pinned: boolean;
    direction: Direction;
}

export const initialState: DrawerState = {
    open: false,
    pinned: false,
    direction: Direction.Right,
};

/**
 * Actions
 */
export const enum DrawerActions {
    OpenDrawer,
    CloseDrawer,
    TogglePin,
}

/**
 * Reducer for the Drawer state.
 * @param state
 * @param action
 * @returns
 */

export const drawerReducer = (state: DrawerState, action: { payload?: any; type: DrawerActions }): DrawerState => {
    switch (action.type) {
        case DrawerActions.OpenDrawer:
            return {
                ...state,
                open: true,
            };
        case DrawerActions.TogglePin:
            return {
                ...state,
                pinned: !state.pinned,
            };
        case DrawerActions.CloseDrawer:
            return {
                ...state,
                open: false,
                pinned: false,
            };
    }
};
