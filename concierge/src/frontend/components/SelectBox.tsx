import * as React from 'react';
import { CdsSelect } from '@cds/react/select';

// Expects:
// props.label      label for listbox
// props.onSelect   callback which takes an event (use event.target.value to get value)
// props.data       (array of strings)
// props.id         HTML id of listbox
function SelectBox(props) {
    return ( <div>
            <CdsSelect layout="vertical">
                <label>{props.label}</label>
                <select onChange={props.onSelect } key={props.id}>
                    {props.data ? props.data.map((key => <option key={key}>{key}</option> )) : ''}
                </select>
            </CdsSelect>
        </div>
    )
}

export default SelectBox
