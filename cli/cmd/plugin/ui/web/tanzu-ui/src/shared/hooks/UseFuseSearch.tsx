import Fuse from 'fuse.js';
import debounce from 'lodash/debounce';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { FuseSearchOptions } from './Fuse.interface';

const useFuseSearch = <T,>(list: T[], index: any, options: FuseSearchOptions) => {
    const { initialQuery, ...fuseOptions } = options;
    const [query, updateQuery] = useState<string>(initialQuery);

    useEffect(() => {
        updateQuery(initialQuery);
    }, [initialQuery]);

    const myIndex = Fuse.parseIndex<T>(index);

    const fuse = useMemo<Fuse<T>>(() => new Fuse<T>(list, fuseOptions, myIndex), [list, fuseOptions, myIndex]);

    const hits = useMemo(() => (query ? fuse.search<T>(query) : fuse.search<T>(initialQuery)), [fuse, query, initialQuery]);

    const setQuery = useMemo(() => debounce(updateQuery, 300), []);

    const onSearch = useCallback((value: string) => setQuery(value.trim() || ''), [setQuery]);

    return {
        hits,
        onSearch,
        query,
        setQuery,
    };
};

export default useFuseSearch;
