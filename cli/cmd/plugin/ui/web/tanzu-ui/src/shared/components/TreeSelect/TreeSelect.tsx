import './TreeSelect.scss';

import { Control, useController } from 'react-hook-form';
import React, { useReducer } from 'react';
import { SelectionType, TreeSelectItem } from './TreeSelect.interface';
import { TreeSelectActions, initialState, treeSelectReducer } from './TreeSelect.store';

import { CdsCheckbox } from '@cds/react/checkbox';
import { CdsControlMessage } from '@cds/react/forms';
import { CdsIcon } from '@cds/react/icon';

function TreeSelect({
    control,
    name,
    selectionType,
    data,
    onChange,
}: {
    control: Control<any, any>;
    name: string;
    selectionType: SelectionType;
    data: TreeSelectItem[];
    onChange: (name: string, value: string) => void;
}) {
    const [state, dispatch] = useReducer(treeSelectReducer, { ...initialState, selectionType });

    const {
        field,
        fieldState: { error },
    } = useController({
        control,
        name,
        rules: { required: true },
    });

    const [value, setValue] = React.useState(field.value || new Map());

    const onCheckboxChange = (e: any) => {
        const mCopy = new Map();

        value.forEach((v: boolean, k: string) => {
            mCopy.set(k, false);
        });

        mCopy.set(e.target.value, e.target.checked);

        let selectedValue = '';

        mCopy.forEach((value, key) => {
            if (value) {
                selectedValue = key;
            }
        });

        // send data to react hook form
        field.onChange(selectedValue);

        //emit change event with the selected value
        onChange(name, selectedValue);

        // update local state
        setValue(mCopy);
    };

    const renderNode = (node: TreeSelectItem) => (
        <TreeItem
            key={node.id}
            node={node}
            onChange={onCheckboxChange}
            checked={state.checked.includes(node.id)}
            expanded={state.expanded.includes(node.id)}
            toggleChecked={(id: string) => dispatch({ type: TreeSelectActions.UpdateChecked, payload: id })}
            toggleExpanded={(id: string) => dispatch({ type: TreeSelectActions.UpdateExpanded, payload: id })}
        >
            {Array.isArray(node.children) ? node.children.map((node) => renderNode(node)) : null}
        </TreeItem>
    );

    const renderNodes = (nodes: TreeSelectItem[]) => nodes.map((node) => renderNode(node));

    return (
        <>
            <div>{renderNodes(data)}</div>
            {error && <CdsControlMessage status="error">{error?.message}</CdsControlMessage>}
        </>
    );
}

function TreeItem({
    node,
    checked,
    expanded,
    toggleChecked,
    toggleExpanded,
    children,
    onChange,
}: {
    node: TreeSelectItem;
    checked: boolean;
    expanded: boolean;
    toggleChecked: (id: string) => void;
    toggleExpanded: (id: string) => void;
    children: any;
    onChange: (e: any) => any;
}) {
    return (
        <div className="tree-container" cds-layout="vertical">
            <ul className="tree-list" cds-layout={children.length ? 'p-l:lg m-y:xs' : 'p-l:xll m-l:xs m-y:xs'}>
                <li className="tree-item">
                    <div className="tree-item-label">
                        {children?.length > 0 && (
                            <CdsIcon
                                aria-label="drawer-close"
                                cds-layout="align:right"
                                className="icon drawer-close"
                                shape="angle"
                                direction={expanded ? 'down' : 'right'}
                                size="md"
                                onClick={() => toggleExpanded(node.id)}
                            ></CdsIcon>
                        )}

                        <CdsCheckbox>
                            <label htmlFor="pool"> {node.label} </label>
                            <input
                                type="checkbox"
                                key={node.id}
                                id="pool"
                                checked={checked}
                                value={node.value}
                                onChange={(e) => {
                                    onChange(e);
                                    toggleChecked(node.id);
                                }}
                            />
                        </CdsCheckbox>
                    </div>
                </li>
                {expanded && children}
            </ul>
        </div>
    );
}

export default TreeSelect;
