// App imports
import { combineReducers } from '../../shared/utilities/reducerUtils';
import { formReducer } from './form.reducer';
import { uiReducer } from './ui.reducer';

export default combineReducers({
    data: formReducer,
    ui: uiReducer
});