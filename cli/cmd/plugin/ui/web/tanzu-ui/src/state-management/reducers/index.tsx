// App imports
import { combineReducers } from '../../shared/utilities/Reducer.utils';
import { appReducer } from './App.reducer';
import { formReducer } from './Form.reducer';
import { uiReducer } from './Ui.reducer';

export default combineReducers({
    app: appReducer,
    data: formReducer,
    ui: uiReducer
});

// NOTE: this method's purpose is a side effect: to ensure the given data path will be valid for the state object
export function ensureDataPath(dataPath: string | undefined, state: any): any {
    // if there is no data path, store the data at the top level of the state object
    if (!dataPath) {
        return state
    }
    // This reducer keeps adding empty objects if given part of the path does not exist.
    // Thus if the path is foo.bar.eeyore and only the object foo exists, the reducer starts
    // with the existing foo and then adds foo.bar; then it takes foo.bar and adds foo.bar.eeyore.
    // It then returns the final object (in this case foo.bar.eeyore).
    // NOTE: the field name should NOT be part of the path! The field is added to the final object (elsewhere)
    return dataPath.split('.').reduce<any>((accumulator, pathSegment) => {
        if (!accumulator[pathSegment]) {
            accumulator[pathSegment] = {}
        }
        return accumulator[pathSegment]
    }, state)
}

export function getDataPath(dataPath: string | undefined, state: any): any {
    if (!dataPath) {
        return state
    }
    return dataPath.split('.').reduce<any>((accumulator, pathSegment) => {
        return accumulator?.[pathSegment]
    }, state)
}
