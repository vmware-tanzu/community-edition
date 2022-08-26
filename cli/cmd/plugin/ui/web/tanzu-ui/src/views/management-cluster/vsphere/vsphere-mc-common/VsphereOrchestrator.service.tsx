// React imports
// Library imports
// App imports
import { FormAction, StoreDispatch } from '../../../../shared/types/types';
import { VSPHERE_FIELDS } from '../VsphereManagementCluster.constants';
import { VsphereDefaultsService } from './VsphereDefaults.service';
import {
    VSphereDatastore,
    VSphereFolder,
    VSphereManagementObject,
    VSphereNetwork,
    VsphereService,
    VSphereVirtualMachine,
} from '../../../../swagger-api';
import React from 'react';
import { DefaultOrchestrator } from '../../default-orchestrator/DefaultOrchestrator';
import { SET_DEFAULTS } from '../../../../state-management/actions/Form.actions';

export interface VsphereOrchestratorProps {
    vsphereState: { [key: string]: any };
    vsphereDispatch: StoreDispatch;
    errorObject: { [fieldName: string]: any };
    setErrorObject: (newErrorObject: { [key: string]: any }) => void;
}
export async function initFolders(
    errorObject: { [fieldName: string]: any },
    setErrorObject: React.Dispatch<React.SetStateAction<any>>,
    dispatch: StoreDispatch,
    datacenter: string
) {
    const result = await DefaultOrchestrator.initResources<VSphereFolder>({
        errorObject,
        setErrorObject,
        resourceName: VSPHERE_FIELDS.VMFolder,
        dispatch,
        fetcher: () => VsphereService.getVSphereFolders(datacenter),
        fxnSelectDefault: (folders) => VsphereDefaultsService.selectDefaultFolder(folders),
    });
    return result;
}

export async function initDatastores(
    errorObject: { [fieldName: string]: any },
    setErrorObject: React.Dispatch<React.SetStateAction<any>>,
    dispatch: StoreDispatch,
    datacenter: string
) {
    await DefaultOrchestrator.initResources<VSphereDatastore>({
        errorObject,
        setErrorObject,
        resourceName: VSPHERE_FIELDS.DataStore,
        dispatch,
        fetcher: () => VsphereService.getVSphereDatastores(datacenter),
        fxnSelectDefault: (stores) => VsphereDefaultsService.selectDefaultDatastore(stores),
    });
}

export async function initNetworks(
    errorObject: { [fieldName: string]: any },
    setErrorObject: React.Dispatch<React.SetStateAction<any>>,
    dispatch: StoreDispatch,
    datacenter: string
) {
    await DefaultOrchestrator.initResources<VSphereNetwork>({
        errorObject,
        setErrorObject,
        resourceName: VSPHERE_FIELDS.Network,
        dispatch,
        fetcher: () => VsphereService.getVSphereNetworks(datacenter),
        fxnSelectDefault: (networks) => VsphereDefaultsService.selectDefaultNetwork(networks),
    });
}

export async function initResources(
    errorObject: { [fieldName: string]: any },
    setErrorObject: React.Dispatch<React.SetStateAction<any>>,
    dispatch: StoreDispatch,
    datacenter: string
) {
    // NOTE: we do not set a default, so we don't send a fxnSelectDefault function
    await DefaultOrchestrator.initResources<VSphereManagementObject>({
        errorObject,
        setErrorObject,
        resourceName: VSPHERE_FIELDS.Pool,
        dispatch,
        fetcher: () => VsphereService.getVSphereComputeResources(datacenter),
    });
}

export async function initOsImages(datacenter: string, props: VsphereOrchestratorProps) {
    const { errorObject, setErrorObject, vsphereDispatch } = props;
    return await DefaultOrchestrator.initResources<VSphereVirtualMachine>({
        errorObject,
        setErrorObject,
        dispatch: vsphereDispatch,
        resourceName: VSPHERE_FIELDS.VMTEMPLATE,
        fetcher: () => VsphereService.getVSphereOsImages(datacenter),
        fxnSelectDefault: VsphereDefaultsService.selectDefaultOsImage,
    });
}

export function initDefaults(dispatch: StoreDispatch) {
    dispatch({ type: SET_DEFAULTS, payload: VsphereDefaultsService.getStaticDefaults(), field: '' } as FormAction);
}

export class VsphereOrchestrator {}
