// import { useContext } from "react";
// import { GlobalAppState } from "../providers";
// import { Actions } from "../providers/globalAppState/reducer";

// export function provisionIngress() {
//     const { dispatch } = useContext(GlobalAppState);

//     console.log("Provision ingress");
//     window.ddClient.extension.vm.cli.exec(`/backend/cluster-provision`)
//         .then((cmdResult: any) => {
//             dispatch({ type: "CLUSTER_LOGS", payload: "Cluster is being provisioned" } as Actions);
//             dispatch({ type: "PROVISION_INGRESS", payload: true } as Actions);
//         });
// }
export async function openBrowser(url: string) {
    return window.ddClient.host.openExternal(url)
}