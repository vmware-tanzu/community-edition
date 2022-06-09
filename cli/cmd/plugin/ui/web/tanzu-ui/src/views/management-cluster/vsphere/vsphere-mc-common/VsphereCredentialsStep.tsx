// React imports
import React, { ChangeEvent, useContext, useState } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';
// Library imports
import { CdsButton } from '@cds/react/button';
import { CdsIcon } from '@cds/react/icon';
import { yupResolver } from '@hookform/resolvers/yup/dist/yup';
// App imports
import { StepProps } from '../../../../shared/components/wizard/Wizard';
import { vsphereCredentialFormSchema } from './vsphere.credential.form.schema';
import { CdsInput } from '@cds/react/input';
import { CdsControlMessage, CdsFormGroup } from '@cds/react/forms';
import { VSPHERE_FIELDS } from './VsphereManagementClusterCommon.constants';
import { VsphereStore } from '../../../../state-management/stores/Store.vsphere.mc';
import { INPUT_CHANGE } from '../../../../state-management/actions/Form.actions';
import { VSphereCredentials, VsphereService } from '../../../../swagger-api';

export interface FormInputs {
    [VSPHERE_FIELDS.SERVERNAME]: string;
    [VSPHERE_FIELDS.USERNAME]: string;
    [VSPHERE_FIELDS.PASSWORD]: string;
}

const SERVER_RESPONSE_BAD_CREDENTIALS = 403;

export function VsphereCredentialsStep(props: Partial<StepProps>) {
    const { vsphereState } = useContext(VsphereStore);
    const { handleValueChange, currentStep, goToStep, submitForm } = props;
    const [connected, setConnection] = useState(false);
    const [connectionErrorMessage, setConnectionErrorMessage] = useState('');
    const methods = useForm<FormInputs>({
        resolver: yupResolver(vsphereCredentialFormSchema),
    });

    const {
        register,
        handleSubmit,
        setValue,
        formState: { errors },
    } = methods;

    // This factory produces an event handler which takes a change event (from a field) and updates both the store value and
    // the validation/error-status of the field
    const handleCredentialFieldChangeFactory = (fieldName: VSPHERE_FIELDS): ((event: ChangeEvent<HTMLInputElement>) => void) => {
        return (event: ChangeEvent<HTMLInputElement>) => {
            setConnection(false);
            setConnectionErrorMessage('');
            if (handleValueChange) {
                handleValueChange(INPUT_CHANGE, fieldName, event.target.value, currentStep, errors);
                setValue(fieldName, event.target.value, { shouldValidate: true });
            }
        };
    };

    const canContinue = (): boolean => {
        return connected && Object.keys(errors).length === 0;
    };

    const onSubmit: SubmitHandler<FormInputs> = (data) => {
        if (canContinue() && goToStep && currentStep && submitForm) {
            goToStep(currentStep + 1);
            submitForm(currentStep);
        }
    };

    // TODO: return true if reasonable data is entered in the form fields
    const connectionDataEntered = (): boolean => {
        return (
            !errors[VSPHERE_FIELDS.SERVERNAME] &&
            vsphereState.data[VSPHERE_FIELDS.SERVERNAME] &&
            vsphereState.data[VSPHERE_FIELDS.USERNAME] &&
            vsphereState.data[VSPHERE_FIELDS.PASSWORD]
        );
    };

    const handleConnect = () => {
        // TODO: use verifyThumbprint
        setConnectionErrorMessage('');
        setConnection(false);
        const vSphereCredentials = {
            username: vsphereState.data[VSPHERE_FIELDS.USERNAME],
            host: vsphereState.data[VSPHERE_FIELDS.SERVERNAME],
            password: vsphereState.data[VSPHERE_FIELDS.PASSWORD],
        } as VSphereCredentials;
        VsphereService.setVSphereEndpoint(vSphereCredentials).then(
            () => {
                setConnection(true);
            },
            (reason: any) => {
                console.warn(`When trying to connect to the server, encountered error ${JSON.stringify(reason)}`);
                let msg = 'Unable to connect to server! (See console for details.)';
                if (reason?.status === SERVER_RESPONSE_BAD_CREDENTIALS) {
                    msg = 'Incorrect username/password combination! Please try again. (See console for technical details.)';
                }
                setConnectionErrorMessage(msg);
            }
        );
    };

    // NOTE: these error objects are defined to keep IDE happy (rather than inline in the HTML below)
    const errorServerName = errors[VSPHERE_FIELDS.SERVERNAME];
    const errorUserName = errors[VSPHERE_FIELDS.USERNAME];
    const errorPassword = errors[VSPHERE_FIELDS.PASSWORD];
    // TODO: add form fields (and disconnect on changes to them)
    return (
        <div>
            <div className="wizard-content-container">
                <h2 cds-layout="m-t:lg">vSphere Credentials</h2>
                <p cds-layout="m-y:lg" className="description">
                    Provide the vCenter server user credentials to create the Management Servicer on vSphere.
                </p>
                <p cds-layout="m-y:lg" className="description">
                    Don&apos;t have vSphere credentials? View our guide on{' '}
                    <a href="/" className="text-blue">
                        creating vSphere credentials
                    </a>
                    .
                </p>
                <div cds-layout="p-t:lg">
                    <CdsFormGroup layout="vertical-inline" control-width="shrink">
                        <div cds-layout="horizontal gap:lg align:vertical-center">
                            <CdsInput layout="compact">
                                <label cds-layout="p-b:md">vSphere server</label>
                                <input
                                    {...register(VSPHERE_FIELDS.SERVERNAME)}
                                    placeholder="vSphere server"
                                    onChange={handleCredentialFieldChangeFactory(VSPHERE_FIELDS.SERVERNAME)}
                                    defaultValue={vsphereState.data[VSPHERE_FIELDS.SERVERNAME]}
                                ></input>
                                {errorServerName && <CdsControlMessage status="error">{errorServerName.message}</CdsControlMessage>}
                            </CdsInput>
                            <CdsInput layout="compact">
                                <label cds-layout="p-b:md">Username</label>
                                <input
                                    {...register(VSPHERE_FIELDS.USERNAME)}
                                    placeholder="username"
                                    onChange={handleCredentialFieldChangeFactory(VSPHERE_FIELDS.USERNAME)}
                                    defaultValue={vsphereState.data[VSPHERE_FIELDS.USERNAME]}
                                ></input>
                                {errorUserName && <CdsControlMessage status="error">{errorUserName.message}</CdsControlMessage>}
                            </CdsInput>
                            <CdsInput layout="compact">
                                <label cds-layout="p-b:md">Password</label>
                                <input
                                    {...register(VSPHERE_FIELDS.PASSWORD)}
                                    placeholder="password"
                                    type="password"
                                    onChange={handleCredentialFieldChangeFactory(VSPHERE_FIELDS.PASSWORD)}
                                    defaultValue={vsphereState.data[VSPHERE_FIELDS.PASSWORD]}
                                ></input>
                                {errorPassword && <CdsControlMessage status="error">{errorPassword.message}</CdsControlMessage>}
                            </CdsInput>
                        </div>
                    </CdsFormGroup>
                    <div>
                        <br />
                    </div>
                    <div>
                        {connectionErrorMessage && (
                            <div>
                                <CdsControlMessage status="error">{connectionErrorMessage}</CdsControlMessage>
                                <br />
                            </div>
                        )}
                        {connected && (
                            <div>
                                <CdsControlMessage status="success">Connection established</CdsControlMessage>
                                <br />
                            </div>
                        )}
                        {!connectionErrorMessage && !connected && (
                            <div>
                                <br />
                            </div>
                        )}
                    </div>
                    <CdsButton onClick={handleConnect} disabled={connected || !connectionDataEntered()}>
                        <CdsIcon shape="connect" size="md"></CdsIcon>
                        {connected ? 'CONNECTED' : 'CONNECT'}
                    </CdsButton>
                </div>
                <CdsButton onClick={handleSubmit(onSubmit)} disabled={!canContinue()}>
                    NEXT
                </CdsButton>
            </div>
        </div>
    );
}
