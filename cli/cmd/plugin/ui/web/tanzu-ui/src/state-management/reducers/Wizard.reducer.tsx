// App imports
import { groupedReducers } from '../../shared/utilities/Reducer.utils';
import { formReducerDescriptor } from './Form.reducer';
import { uiReducerDescriptor } from './Ui.reducer';

// generic reducer used in wizards
export default groupedReducers({
    name: 'generic wizard reducer',
    reducers: [formReducerDescriptor, uiReducerDescriptor],
});
