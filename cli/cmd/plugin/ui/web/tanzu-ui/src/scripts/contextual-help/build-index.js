/* eslint-disable @typescript-eslint/no-var-requires */
const fs = require('fs');
const fm = require('html-frontmatter');
const removeHtmlComments = require('remove-html-comments');
const Fuse = require('fuse.js');

const docsPath = 'src/assets/contextualHelpDocs';

const contextualHelpData = [];

const loadDoc = (filePath) => {
    const fileContent = fs.readFileSync(filePath, 'utf-8');
    const metaData = fm(fileContent);
    const htmlContent = removeHtmlComments(fileContent);

    contextualHelpData.push({
        ...metaData,
        htmlContent: htmlContent.data,
    });
};

const getAllFiles = (dir) => {
    return new Promise((resolve) => {
        const allPromises = [];
        fs.readdir(dir, (err, files) => {
            if (err) {
                console.log(err);
                return;
            }
            files.forEach((file) => {
                const filePath = `${dir}/${file}`;
                if (fs.statSync(filePath).isDirectory()) {
                    allPromises.push(getAllFiles(filePath));
                } else {
                    if (file.endsWith('html')) {
                        loadDoc(filePath);
                        allPromises.push(Promise.resolve());
                    }
                    allPromises.push(Promise.resolve());
                }
            });
            Promise.all(allPromises).then(resolve);
        });
    });
};

async function buildIndex() {
    const success = await getAllFiles(docsPath);

    if (success) {
        const myIndex = Fuse.createIndex(['topicTitle', 'topicIds', 'topicDescription'], contextualHelpData);
        fs.writeFileSync(docsPath + '/fuse-index.json', JSON.stringify(myIndex));
        fs.writeFileSync(docsPath + '/data.json', JSON.stringify({ data: [...contextualHelpData] }));

        console.log('Fuse Index created successfully.');
    } else {
        console.log('Fuse Index creation failed.');
    }
}

buildIndex();
