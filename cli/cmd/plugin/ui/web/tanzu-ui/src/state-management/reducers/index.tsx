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