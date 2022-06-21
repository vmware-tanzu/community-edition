import Fuse from 'fuse.js';
import debounce from 'lodash/debounce';
import { useCallback, useMemo, useState } from 'react';
import { FuseSearchOptions } from './Fuse.interface';

const useFuseSearch = <T,>(list: T[], index: any, options: FuseSearchOptions) => {
    const { defaultQuery, ...fuseOptions } = options;

    const [query, updateQuery] = useState<string>(defaultQuery);

    const myIndex = Fuse.parseIndex<T>(index);

    const fuse = useMemo<Fuse<T>>(() => new Fuse<T>(list, fuseOptions, myIndex), [list, fuseOptions, myIndex]);

    // memoize results whenever the query or options change
    const hits = useMemo(
        // if query is empty and `matchAllOnEmptyQuery` is `true` then return all list
        // NOTE: we remap the results to match the return structure of `fuse.search()`
        () => (!query ? fuse.search<Fuse.FuseResult<T>[]>(defaultQuery) : fuse.search<Fuse.FuseResult<T>[]>(query)),
        [fuse, query, defaultQuery]
    );

    // debounce updateQuery and rename it `setQuery` so it's transparent
    const setQuery = useMemo(() => debounce(updateQuery, 300), []);

    // pass a handling helper to speed up implementation
    const onSearch = useCallback((value: string) => setQuery(value.trim() || ''), [setQuery]);

    // still returning `setQuery` for custom handler implementations
    return {
        hits,
        onSearch,
        query,
        setQuery,
    };
};

export default useFuseSearch;
