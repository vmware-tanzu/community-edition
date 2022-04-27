// App imports
import { combineReducers } from '../../shared/utilities/Reducer.utils';
import { formReducer } from './Form.reducer';
import { uiReducer } from './Ui.reducer';

export default combineReducers({
    data: formReducer,
    ui: uiReducer
});