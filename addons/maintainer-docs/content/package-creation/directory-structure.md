---
title: Directory Structure
---

Packages should conform to the following directory structure.

```shell
├── 1.2.3
│   ├── README.md
│   ├── bundle
│   │   ├── .imgpkg
│   │   │   └── images.yml
│   │   ├── config
│   │   │   ├── overlays
│   │   │   │   ├── overlay-a.yaml
│   │   │   │   └── overlay-b.yaml
│   │   │   ├── upstream
│   │   │   │   ├── upstream-a.yaml
│   │   │   │   └── upstream-b.yaml
│   │   │   └── values.yaml
│   │   ├── schema.yml
│   │   ├── vendir.lock.yml
│   │   └── vendir.yml
│   └── package.yaml
├── metadata.yaml
└── test
    ├── Makefile
    ├── README.md
    ├── e2e
    │   └── test-a.go
    └── unittest
        └── test-b.go
```
