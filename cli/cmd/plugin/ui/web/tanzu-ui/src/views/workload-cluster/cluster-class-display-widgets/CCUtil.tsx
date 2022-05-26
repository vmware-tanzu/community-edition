import { CCVariable } from '../../../shared/models/ClusterClass';

// NOTE: this fxn has the intentional side effect of populating default values for the children of the given ccVar
//       To change this to a "pure" fxn, we would need to do a deep clone.
// The incoming param object "defaults" is expected to have attributes corresponding to any child var of ccVar that has a default value
export function populateDefaults(defaults: any, ccVar: CCVariable): CCVariable {
    if (defaults) {
        Object.keys(defaults).forEach((key) => {
            const child = ccVar.children?.find((child) => child.name === key);
            if (child) {
                child.default = defaults[key];
            } else {
                console.warn(
                    `While populating default values for ${ccVar.name}, the default object (${JSON.stringify(defaults)}) ` +
                        `had a key of ${key}, but no such child field exists in our target object's children: ${JSON.stringify(ccVar)}`
                );
            }
        });
    }
    return ccVar;
}
