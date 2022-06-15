// App imports

// vSphere MC reducer
import { formReducerDescriptor } from '../../../state-management/reducers/Form.reducer';
import { groupedReducers } from '../../../shared/utilities/Reducer.utils';
import { uiReducerDescriptor } from '../../../state-management/reducers/Ui.reducer';
import { vsphereResourceReducerDescriptor } from './VsphereResources.reducer';

export default groupedReducers({
    name: 'vSphere MC reducer',
    reducers: [uiReducerDescriptor, formReducerDescriptor, vsphereResourceReducerDescriptor],
});
