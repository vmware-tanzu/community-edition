// App imports
import { combineReducers } from '../../shared/utilities/Reducer.utils';
import { formReducer } from './Form.reducer';
import { uiReducer } from './Ui.reducer';

// generic reducer used in wizards
export default combineReducers({
    data: formReducer,
    ui: uiReducer,
});
