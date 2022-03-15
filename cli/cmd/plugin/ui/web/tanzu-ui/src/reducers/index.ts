
import { combineReducers } from '../shared/utilities/reducerUtils';
import { formReducer } from './formReducer';
import { uiReducer } from './uiReducer';

export default combineReducers({
    data: formReducer,
    ui: uiReducer
});