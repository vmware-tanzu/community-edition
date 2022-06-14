// App imports

// vSphere MC reducer
import { combineReducers } from '../../../shared/utilities/Reducer.utils';
import { formReducer } from '../../../state-management/reducers/Form.reducer';
import { uiReducer } from '../../../state-management/reducers/Ui.reducer';
import { resourcesReducer } from '../../../state-management/reducers/Resources.reducer';

export default combineReducers({
    data: formReducer,
    ui: uiReducer,
    resources: resourcesReducer,
});
