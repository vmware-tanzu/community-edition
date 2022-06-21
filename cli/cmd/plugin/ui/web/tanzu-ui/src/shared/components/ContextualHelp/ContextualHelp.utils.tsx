import Fuse from 'fuse.js';

export const generateDataAccordionContent = (content: Fuse.FuseResult<unknown>[]) =>
    content.map((item: { item: any; refIndex: any }) => ({
        id: item.refIndex,
        title: item.item.topicTitle,
        content: item.item.htmlContent,
    }));
