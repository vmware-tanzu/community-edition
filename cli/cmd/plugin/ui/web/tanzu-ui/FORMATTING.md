# File and Code Formatting Best Practices

## All React Typescript Files

There are several types of React classes or components that you might create when developing in the
TCE Guided UI. Whether creating a new component, store, actions, or even constants; your file should be suffixed with `.tsx`.

### Naming Your File

All filename prefixes should be pascal case (example: MyReactComponent.tsx).

Here are some general guidelines for the filenames:

- Component - should consist of component name plus suffix. (example: MyReactComponent.tsx)
- Constants - should include `constants` between name and suffix (example: App.constants.tsx)
- Stores - should include `store` between name and suffix (example: App.store.tsx)
- Actions - should include `actions` between name and suffix (example: App.actions.tsx)
- Reducers - should include `reducers` between name and suffix (example: App.reducer.tsx)

### Code Formatting and Linting

Code formatting can be applied automatically by leveraging Prettier, which can be accessed via npm scripts.

To run formatting: `/community-edition/cli/cmd/plugin/ui/web/tanzu-ui npm run format`

To run formatting and unit tests (Pre-commit checks): `/community-edition/cli/cmd/plugin/ui/web/tanzu-ui npm run precommit`

Code linting checks can be manually, but are also a side-effect of starting the local dev environment via `npm run start`

To run linting alone manually: `/community-edition/cli/cmd/plugin/ui/web/tanzu-ui npm run lint`
