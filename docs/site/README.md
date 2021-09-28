# Website for [Template]

## Prerequisites

* [Hugo](https://github.com/gohugoio/hugo)
  * macOS: `brew install hugo`
  * Windows: `choco install hugo-extended -confirm`

## Build

```bash
hugo server --disableFastRender
```

## Serve

Serve site at [http://localhost:1313](http://localhost:1313)

## NPM dependencies

This docs site is _not_ an NPM module and is _not_ served using NodeJS.
It is a static Hugo site.
However, there are several NPM dependencies and there is a `package.json` file to reflect that.
This enables the maintainers to receive DependaBot security alerts through GitHub for those static dependencies.

For example, the following is declared as a dependency in the the `package.json` file:

```json
"docsearch.js": "2.6.3"
```

and is being pulled in statically via CDN in the base HTML file:

```html
<script type="text/javascript" src="https://cdn.jsdelivr.net/npm/docsearch.js@2.6.3/dist/cdn/docsearch.min.js" integrity="sha384-8uEk67aWSZHvjtAX9hf2AB+KzYcssy31vRRTi9oP81zHtyIj7PQGAykGbQpB1L2J" crossorigin="anonymous"></script>
```

Also note the `integrity` attribute.
This is the shasum value for that minified JavaScript file being pulled from through NPM and the CDN
and ensures that the JavaScript is not loaded in the case that the shasum changes or is malicously modified.

The shasum can be generated as follows:

```sh
curl https://cdn.jsdelivr.net/npm/docsearch.js@2.6.3/dist/cdn/docsearch.min.js | shasum -b -a 384 | awk '{ print $1 }' | xxd -r -p | base64
```

and is [documented in the MDN Web Docs here](https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity).
Whenever a dependency version is upgraded, ensure the `integrity` shasum is also updated.
