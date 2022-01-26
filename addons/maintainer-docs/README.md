# Package Maintainer Documentation

This directory contains package owner and maintainer documentation.

## How To Run

The geekdocs theme has a number of misspellings and a lot of JavaScript. Let's not pollute our repo with that. The themes folder has been git ignored, so you'll have to download the theme before running the server.

```shell
mkdir -p themes/hugo-geekdoc/
curl -L https://github.com/thegeeklab/hugo-geekdoc/releases/latest/download/hugo-geekdoc.tar.gz | tar -xz -C themes/hugo-geekdoc/ --strip-components=1
```

This is a Hugo site. It can be started locally with the following command. The site is then accessible [locally](http://localhost:1313).

```shell
hugo server -D
```
