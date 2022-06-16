import { Direction, DrawerActions } from './Drawer.enum';

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
