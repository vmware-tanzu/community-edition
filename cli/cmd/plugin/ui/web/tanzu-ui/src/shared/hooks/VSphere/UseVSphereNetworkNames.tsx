import { VSphereNetwork, VsphereService } from '../../../swagger-api';
import { useEffect, useState } from 'react';

const useVSphereNetworkNames = (dc: string) => {
    const [data, setData] = useState<VSphereNetwork[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    useEffect(() => {
        setLoading(true);
        VsphereService.getVSphereNetworks(dc)
            .then((response) => {
                setData(response);
            })
            .catch((err) => {
                setError(err);
            })
            .finally(() => {
                setLoading(false);
            });
    }, [dc]);
    return { data, loading, error };
};

export default useVSphereNetworkNames;
