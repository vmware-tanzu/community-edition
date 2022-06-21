import { MouseEvent, ReactNode } from 'react';
import './Drawer.scss';
import { DrawerState } from './Drawer.store';

export interface DrawerProps extends DrawerEvents, DrawerState {
    children: ReactNode;
}

export interface DrawerEvents {
    onClose: (event: MouseEvent<HTMLElement>) => void;
    togglePin: (event: MouseEvent<HTMLElement>) => void;
}
