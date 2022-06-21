import Fuse from 'fuse.js';
import debounce from 'lodash/debounce';
import { useCallback, useMemo, useState } from 'react';
import { FuseSearchOptions } from './Fuse.interface';

const useFuseSearch = <T,>(list: T[], index: any, options: FuseSearchOptions) => {
    const { defaultQuery, ...fuseOptions } = options;

    const [query, updateQuery] = useState<string>(defaultQuery);

    const myIndex = Fuse.parseIndex<T>(index);

    const fuse = useMemo<Fuse<T>>(() => new Fuse<T>(list, fuseOptions, myIndex), [list, fuseOptions, myIndex]);

    const hits = useMemo(() => (query ? fuse.search<T>(query) : fuse.search<T>(defaultQuery)), [fuse, query, defaultQuery]);

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
