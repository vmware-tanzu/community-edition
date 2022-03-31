// App imports
import { combineReducers } from '../../shared/utilities/Reducer.utils';
import { appReducer } from './app.reducer';
import { formReducer } from './form.reducer';
import { uiReducer } from './ui.reducer';

export default combineReducers({
    app: appReducer,
    data: formReducer,
    ui: uiReducer
});