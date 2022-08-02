import { VSphereManagementObject, VsphereService } from '../../../swagger-api';
import { useEffect, useState } from 'react';

const useVSphereComputeResources = (dc: string) => {
    const [data, setData] = useState<VSphereManagementObject[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    useEffect(() => {
        setLoading(true);
        VsphereService.getVSphereComputeResources(dc)
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

export default useVSphereComputeResources;
