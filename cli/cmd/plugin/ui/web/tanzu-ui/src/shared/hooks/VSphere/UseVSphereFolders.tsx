import { VSphereFolder, VsphereService } from '../../../swagger-api';
import { useEffect, useState } from 'react';

const useVSphereFolders = (dc: string) => {
    const [data, setData] = useState<VSphereFolder[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    useEffect(() => {
        setLoading(true);
        VsphereService.getVSphereFolders(dc)
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

export default useVSphereFolders;
