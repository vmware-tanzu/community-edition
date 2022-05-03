# Tanzu Guided UI Developer Docs

This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).

To learn React, check out the [React documentation](https://reactjs.org/).

## Prerequisites for Building and Running the UI

Node version 16.x.x (latest LTS version) on your local machine

- `node --version` to check which version you have

Node version can be set and managed by using NVM (Node Version Manager):

- `brew install nvm`
- `nvm install 16` (or `nvm use 16` to temporarily set node version)

If an alternate NPM registry is required to obtain the node dependencies, it should be configured either

- prior to running the make target, with `npm config set registry <register-url>`, or
- providing the URL in the CUSTOM_NPM_REGISTRY environment variable.



## Building and Running the UI Locally

*Note: Prior to running the UI in a local development mode, you must install all required npm packages (See `npm ci`).

In the `tanzu-ui` project directory, you can run:
### `npm ci`

Installs all required npm packages and versions from the `package-lock.json` file. Node version should be 16.x as noted above.

### `npm start`

Runs the app in the development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.\
You will also see any lint errors in the console.

### `npm test`

Launches the test runner in the interactive watch mode.\
See the section about [running tests](https://create-react-app.dev/docs/running-tests/) for more information.

See also [testing-library/React](https://testing-library.com/docs/react-testing-library/intro/) for reference to APIs used when writing tests for React components

### `npm lint`

Runs ESLint code linting against all js, ts, and tsx files in the `src` folder.

ESLint rules are defined in `tanzu-ui/package.json` and should be updated as developers see fit.

### `npm run generate-api`

Auto-generates all Swagger REST API models and methods in Typescript for UI

### `npm run build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.

## Developing in the UI

[Code/File formatting best practices](./FORMATTING.md)

Directory structure

React basics in practice


