// Azure MC reducer
import { formReducerDescriptor } from '../../../state-management/reducers/Form.reducer';
import { groupedReducers } from '../../../shared/utilities/Reducer.utils';
import { uiReducerDescriptor } from '../../../state-management/reducers/Ui.reducer';
import { azureResourceReducerDescriptor } from './AzureResources.reducer';

export default groupedReducers({
    name: 'Azure MC reducer',
    reducers: [uiReducerDescriptor, formReducerDescriptor, azureResourceReducerDescriptor],
});
