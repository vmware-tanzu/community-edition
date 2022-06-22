import { VSphereVirtualMachine } from '../../../../swagger-api';

// analyzeOsImages takes a datacenter name, a URL (for messaging) and an array of osImages. It returns an object that includes
// the number of images, the number of templates, and the relevant user message.
// If the datacenter name is empty, then an empty object is returned (the assumption is that the user hasn't selected a datacenter).
export function analyzeOsImages(
    datacenter: string,
    urlConvertOsImageToTemplate: string,
    osImages: VSphereVirtualMachine[]
): { msg: string; nImages: number; nTemplates: number } {
    // If no datacenter has been selected, then just return an empty object
    if (!datacenter) {
        return { msg: '', nTemplates: 0, nImages: 0 };
    }

    const nImages = osImages?.length || 0;
    const nTemplates = osImages?.reduce<number>((accum, image) => accum + (image.isTemplate ? 1 : 0), 0) || 0;
    let msg = '';
    if (nImages === 0) {
        msg = `No OS images are available! Please select a different data center or add an OS image to ${datacenter}`;
    } else if (nTemplates === 0) {
        const describeNumTemplates = nImages === 1 ? 'There is one OS image' : `There are ${nImages} OS images`;
        const notATemplate = nImages === 1 ? 'it is not a template' : 'none of them are templates';
        msg =
            `${describeNumTemplates} on data center ${datacenter}, but ${notATemplate}.` +
            `For information on how to convert an OS image to a template, see URL ${urlConvertOsImageToTemplate}`;
    }
    return { msg, nImages, nTemplates };
}
