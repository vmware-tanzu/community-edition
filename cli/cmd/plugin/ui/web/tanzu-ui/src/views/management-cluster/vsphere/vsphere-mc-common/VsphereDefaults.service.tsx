import { VSphereDatastore, VSphereFolder, VSphereNetwork, VSphereVirtualMachine } from '../../../../swagger-api';
import { filterTemplates } from './VsphereOsImageUtil';
import { first } from '../../../../shared/utilities/Array.util';

export class VsphereDefaultsService {
    static selectDefaultOsImage = (osImages: VSphereVirtualMachine[]) => {
        return first<VSphereVirtualMachine>(filterTemplates(osImages));
    };

    static selectDefaultNetwork = (networks: VSphereNetwork[]) => {
        return first<VSphereNetwork>(networks);
    };

    static selectDefaultDatastore = (datastores: VSphereDatastore[]) => {
        return first<VSphereDatastore>(datastores);
    };

    static selectDefaultFolder = (folders: VSphereFolder[]) => {
        return first<VSphereFolder>(folders);
    };
}
