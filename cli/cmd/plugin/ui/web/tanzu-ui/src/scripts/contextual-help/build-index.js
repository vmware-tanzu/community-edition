/* eslint-disable @typescript-eslint/no-var-requires */
const fs = require('fs');
const fm = require('html-frontmatter');
const stripHtmlComments = require('strip-html-comments');
const Fuse = require('fuse.js');

const docsPath = 'src/assets/contextualHelpDocs';

const fetchFileContent = (filePath) => {
    try {
        const fileContent = fs.readFileSync(filePath, 'utf-8');
        const metaData = fm(fileContent);
        const htmlContent = stripHtmlComments(fileContent);
        return {
            ...metaData,
            htmlContent,
        };
    } catch (e) {
        console.err(`Error Reading File Content for path: ${filePath}`, e);
        return;
    }
};

function* fetchFiles(path, fileType) {
    const entries = fs.readdirSync(path);
    for (let file of entries) {
        const filePath = `${path}/${file}`;

        if (fs.statSync(filePath).isDirectory()) {
            yield* fetchFiles(filePath, fileType);
        } else if (file.endsWith(fileType)) {
            yield fetchFileContent(filePath);
        }
    }
}

const buildIndex = () => {
    try {
        const files = [...fetchFiles(docsPath, 'html')];
        const myIndex = Fuse.createIndex(['topicTitle', 'topicIds', 'topicDescription'], [...files]);
        fs.writeFileSync(docsPath + '/fuse-index.json', JSON.stringify(myIndex));
        fs.writeFileSync(docsPath + '/data.json', JSON.stringify({ data: [...files] }));
        console.log('Fuse Index created successfully.');
    } catch (e) {
        console.err('Fuse Index creation failed.');
    }
};

buildIndex();
