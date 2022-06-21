export interface DataAccordionConfig<T> {
    data: T[];
    key: (item: T) => number;
    title: (item: T) => string;
    content: (item: T) => string;
}
