// React imports
import React from 'react';

// Library imports
import { CdsButton } from '@cds/react/button';

//App imports
import { STATUS } from '../shared/constants/App.constants';
import { ChildProps } from '../shared/components/wizard/StepNav';

function TestRender (props: ChildProps | any) {
    return (
        <div>
            hello world

            <CdsButton onClick={() => {
                props.goToStep(props.currentStep + 1);
                const tabStatus = [...props.tabStatus];
                tabStatus[props.currentStep - 1] = STATUS.VALID;
                props.setTabStatus(tabStatus);
            }}>Next</CdsButton>
        </div>
    );
}

export default TestRender;