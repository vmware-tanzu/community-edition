// App imports
import { groupedReducers } from '../../shared/utilities/Reducer.utils';
import { appReducerDescriptor } from './App.reducer';
import { formReducerDescriptor } from './Form.reducer';
import { uiReducerDescriptor } from './Ui.reducer';

export default groupedReducers({
    name: 'default app reducer',
    reducers: [appReducerDescriptor, formReducerDescriptor, uiReducerDescriptor],
});

// NOTE: this method's purpose is a side effect: to ensure the given data path will be valid for the state object
export function ensureDataPath(dataPath: string | undefined, separator: string, state: any): any {
    // if there is no data path, store the data at the top level of the state object
    if (!dataPath) {
        return state;
    }
    // This reducer keeps adding empty objects if given part of the path does not exist.
    // Thus if the path is foo.bar.eeyore and only the object foo exists, the reducer starts
    // with the existing foo and then adds foo.bar; then it takes foo.bar and adds foo.bar.eeyore.
    // It then returns the final object (in this case foo.bar.eeyore).
    // NOTE: the field name should NOT be part of the path! The field is added to the final object (elsewhere)
    return dataPath.split(separator).reduce<any>((accumulator, pathSegment) => {
        if (!accumulator[pathSegment]) {
            accumulator[pathSegment] = {};
        }
        return accumulator[pathSegment];
    }, state);
}

// returns the object at the given data path (or undefined)
export function getDataPath(dataPath: string | undefined, separator: string, state: any): any {
    if (!dataPath) {
        return state;
    }
    return getObject(dataPath.split(separator), state);
}

function getObject(pathElements: string[], state: any): any {
    return pathElements.reduce<any>((accumulator, pathSegment) => {
        return accumulator?.[pathSegment];
    }, state);
}

// given a data path, prune empty objects up the tree
export function pruneDataPath(dataPath: string | undefined, separator: string, state: any): any {
    if (!dataPath) {
        return state;
    }
    const arrPathElements = dataPath.split(separator);
    arrPathElements.reduceRight<any>((_, value: string, index: number, arrayPath) => {
        const attribute = value;
        const obj = getObject(arrayPath.slice(0, index), state);
        if (isEmptyObject(obj[attribute])) {
            delete obj[attribute];
        }
        return null;
    }, null);
    return state;
}

function isEmptyObject(obj: any): boolean {
    return !obj || Object.keys(obj).length === 0;
}
