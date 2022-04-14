import * as React from 'react';
import { CdsSelect } from '@cds/react/select';

// Expects:
// props.label      label for listbox
// props.onSelect   callback which takes a version string
// props.data       (array of strings)
// props.id         HTML id of listbox
function SelectBox(props) {
    return ( <div>
            <CdsSelect layout="vertical">
                <label>{props.label}</label>
                <select onChange={event => props.onSelect(event.target.value)} key={props.id}>
                    {props.leadingBlankItem ? <option id='' key=''></option> : ''}
                    {props.data ? props.data.map((key => <option id={key} key={key}>{key}</option> )) : ''}
                </select>
            </CdsSelect>
        </div>
    )
}

export default SelectBox
