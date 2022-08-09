// App imports

// Aws MC reducer
import { formReducerDescriptor } from '../../../state-management/reducers/Form.reducer';
import { groupedReducers } from '../../../shared/utilities/Reducer.utils';
import { uiReducerDescriptor } from '../../../state-management/reducers/Ui.reducer';
import { awsResourceReducerDescriptor } from './AwsResources.reducer';

export default groupedReducers({
    name: 'Aws MC reducer',
    reducers: [uiReducerDescriptor, formReducerDescriptor, awsResourceReducerDescriptor],
});
