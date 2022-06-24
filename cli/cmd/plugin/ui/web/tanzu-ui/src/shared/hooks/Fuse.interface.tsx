import Fuse from 'fuse.js';

export interface FuseSearchOptions extends Fuse.IFuseOptions<any> {
    defaultQuery: string;
}
