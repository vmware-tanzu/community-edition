// React imports
import { useEffect } from 'react';

// Library imports
import { FieldErrors, FieldValues, useForm, UseFormProps, UseFormReturn } from 'react-hook-form';

// App imports
import { DependencyMap } from '../../components/DependencyMap';
import { RESET_DEPENDENT_FIELDS } from '../../state-management/actions/Form.actions';
import { StoreDispatch } from '../../state-management/stores/Store';
import { STATUS } from '../constants/App.constants';

export const useWizardForm = <TFieldValues extends FieldValues = FieldValues, TContext = any>
    (props?: UseFormProps<TFieldValues, TContext>): any => {
    const formObj: UseFormReturn<TFieldValues, TContext>= useForm({ ...props });
    const handleFormSubmit = (props: any) => {
        console.log(`Did you know that we got a form submit with currentStep=${props.currentStep}?`);
        if (!formObj.handleSubmit) {
            console.log('No form handleSubmit?!');
        }
        formObj.handleSubmit((data: any) => {
            console.log(`Before going to the next step, our current step is: ${props.currentStep}`);
            props.goToStep(props.currentStep + 1);
            console.log(`After going to the next step, our current step is: ${props.currentStep}`);
            const tabStatus = [...props.tabStatus];
            tabStatus[props.currentStep - 1] = STATUS.VALID;
            props.setTabStatus(tabStatus);
            props.submitForm(data);
        })();
    };
    return {
        ...formObj,
        handleFormSubmit
    };
};
interface DependenyProps {
    name: string,
    dependencyMap: DependencyMap,
    dispatch: StoreDispatch,

}
export const onBlurWithDependencies = (props: DependenyProps) => {
    const { name, dependencyMap, dispatch } = props;
    dispatch({
        type: RESET_DEPENDENT_FIELDS,
        payload: {
            fields: dependencyMap[name]
        } 
    });
};

export interface TabStatusProps {
    tabStatus: string[],
    setTabStatus: (tabStatus: string[]) => void,
    currentStep: number
}

export const useTabStatus = <TFieldValues extends FieldValues = FieldValues> (
    errors: FieldErrors<TFieldValues>,
    props: TabStatusProps) => {
    const numOfErrors = Object.keys(errors).length;
    const { tabStatus, currentStep, setTabStatus } = props;
    const curStatus = tabStatus[currentStep - 1];
    useEffect(() => {
        if (numOfErrors > 0 &&  (curStatus !== STATUS.VALID)) {
            console.log(`errors: ${JSON.stringify(errors)}`);
            console.log(`numOfErrors is ${numOfErrors} and curStatus is ${curStatus}, so setting status to INVALID`);
            tabStatus[currentStep - 1] = STATUS.INVALID;
            setTabStatus([...tabStatus]);
        } else if ( numOfErrors === 0 && curStatus === STATUS.INVALID) {
            tabStatus[currentStep - 1] = STATUS.VALID;
            setTabStatus([...tabStatus]);
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    },[numOfErrors]);
};

export const blurHandler =  <TFieldValues extends FieldValues = FieldValues>(
    errors: FieldErrors<TFieldValues>,
    submitForm: (data: {[key: string]: string}) => void,
    formValues: {[key: string]: string}) => {
    if (Object.keys(errors).length === 0) {
        submitForm(formValues);
    }
};
