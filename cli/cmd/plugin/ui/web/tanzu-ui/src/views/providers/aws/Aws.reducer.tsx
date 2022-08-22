// App imports
import { formReducerDescriptor } from '../../../state-management/reducers/Form.reducer';
import { groupedReducers } from '../../../shared/utilities/Reducer.utils';
import { resourceReducerDescriptor } from '../../../state-management/reducers/Resources.reducer';
import { resourceWithDefaultReducerDescriptor } from '../../../state-management/reducers/ResourcesWithDefault.reducer';
import { uiReducerDescriptor } from '../../../state-management/reducers/Ui.reducer';

export default groupedReducers({
    name: 'Aws MC reducer',
    reducers: [uiReducerDescriptor, formReducerDescriptor, resourceReducerDescriptor, resourceWithDefaultReducerDescriptor],
});
