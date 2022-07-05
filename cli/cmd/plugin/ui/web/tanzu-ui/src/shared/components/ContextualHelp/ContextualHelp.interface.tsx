import { DrawerEvents } from '../Drawer/Drawer.interface';
import { DrawerState } from '../Drawer/Drawer.store';
import { ContextualHelpState } from './ContextualHelp.store';

export interface SearchProps {
    value?: string;
    onSearch: (value: string) => void;
}

export interface ContextualHelpContentProps extends ContextualHelpState, DrawerEvents, DrawerState {}

export interface ContextualHelpData {
    topicTitle: string;
    topicIds: string[];
    htmlContent: string;
}
