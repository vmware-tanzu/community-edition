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

## Building and Running the UI in Local Developer Mode

\*Note: Prior to running the UI in a local development mode, you must install all required npm packages (See `npm ci`).

In the `tanzu-ui` project directory, you can run:

### `npm ci`

Installs all required npm packages and versions from the `package-lock.json` file. Node version should be 16.x as noted
above.

### `npm run start`

Runs the app in the development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.\
You will also see any lint errors in the console.

### `npm run test`

Launches the test runner in the interactive watch mode.\
See the section about [running tests](https://create-react-app.dev/docs/running-tests/) for more information.

See also [testing-library/React](https://testing-library.com/docs/react-testing-library/intro/) for reference to APIs
used when writing tests for React components

### `npm run format`

Runs Prettier code formatting against all js, ts, and tsx files in the `src` folder.

### `npm run lint`

Runs ESLint code linting against all js, ts, and tsx files in the `src` folder.

ESLint rules are defined in `tanzu-ui/package.json` and should be updated as developers see fit.

### `npm run generate-api`

Auto-generates all Swagger REST API models and methods in Typescript for UI

### `npm run build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.

## Building and Running the UI Plugin Locally in Tanzu CLI

### Step 1: `/community-edition/cli/cmd/plugin/ui make ui-build`

Executing `make ui-build` from the `/community-edition/cli/cmd/plugin/ui` directory will generate all production-ready
UI assets needed for the Tanzu UI plugin.

### Step 2: `/community-edition make build-install-cli-plugins`

Executing `make build-install-cli-plugins` from the `/community-edition` directory will build and install all TCE
plugins into the Tanzu CLI, including the Tanzu UI plugin.

### Launching the Tanzu UI plugin

If your Tanzu CLI and the UI plugin are installed correctly, running `tanzu ui` will start the UI and launch a browser
window at `0.0.0.0:8080`

## Developing in the UI

[Code/File formatting best practices](./FORMATTING.md)

[REST API Consumption](./RESTAPIS.md)

## Build Contextual Help Docs

HTMl Contextual Help Docs are located at `src/assets/contextualHelpDocs`.\
To generate the docs index and JSON data from HTML files Run `npm run build-index`. This should be done everytime a new
topic is introduced.\
Converted HTML Docs to JSON Data are in `data.json`.\
Documents Index is generated in `fuse-index.json`.
