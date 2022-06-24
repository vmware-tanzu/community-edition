import { groupedReducers } from '../../shared/utilities/Reducer.utils';
import { formReducerDescriptor } from '../../state-management/reducers/Form.reducer';
import { uiReducerDescriptor } from '../../state-management/reducers/Ui.reducer';
import { dynamicFormReducerDescriptor } from '../../state-management/reducers/DynamicForm.reducer';

export default groupedReducers({
    name: 'generic wizard reducer',
    reducers: [formReducerDescriptor, dynamicFormReducerDescriptor, uiReducerDescriptor],
});
