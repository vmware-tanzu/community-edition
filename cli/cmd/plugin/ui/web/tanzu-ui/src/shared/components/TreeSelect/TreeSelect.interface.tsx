export interface TreeSelectItem {
    id: string;
    parentId: string;
    label: string;
    value: string;
    type?: string;
    children?: TreeSelectItem[];
}

export const enum SelectionType {
    Single = 'single',
    Multi = 'multi',
}
