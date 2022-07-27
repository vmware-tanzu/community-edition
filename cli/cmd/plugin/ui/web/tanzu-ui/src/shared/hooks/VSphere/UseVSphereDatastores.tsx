import { VSphereDatastore, VsphereService } from '../../../swagger-api';
import { useEffect, useState } from 'react';

const useVSphereDatastores = (dc: string) => {
    const [data, setData] = useState<VSphereDatastore[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    useEffect(() => {
        setLoading(true);
        VsphereService.getVSphereDatastores(dc)
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

export default useVSphereDatastores;
