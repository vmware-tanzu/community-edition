# Metapackage-based Installation Experience

* Proposal: https://github.com/vmware-tanzu/community-edition/issues/2848

## Summary

In order to simplify the experience of installing multiple packages, we wanted to leverage the ability of creating a "metapackage" using Carvel tooling. 

This metapackage is essentially a package of packages wherein it references a Bill of Materials that includes existing packages in the repository. Being able to install, update & remove these packages in concert with each other simplifies the complexity significantly. 

This is a method in practice in other projects leveraging Carvel tooling. 

## Motivation

The motivation for introducing this is to reduce the complexity that comes up with installing multiple carvel packages at a time. This will allow newcomers to the project rapidly establish a working technology stack.

* Describe why the change is important and the benefits to users.  

### Goals

* Demonstrate the ability that multiple subordinate packages can be installed using a single metapackage. 
* Be able to provide a single values.yaml file that can be used for all the subordinate package values.yaml.

### Non-Goals/Future Work

* Define each profile (and the packages it will have) that will use this metapackage technology. 

## Proposal

The proposal will describe what a metapackage is and how it will be implemented. 

The proposal can be leveraged for various user scenarios. One example user scenario is described below. 

**Terminology**

* Metapackage - concept to bring two or more packages together. 
* Profile - Enabling users to decide by a config value which part of the metapackage the user wants to install. 



### Example Scenario

The below scenario is one use-case of the metapackage technology. This is **not** the only way to use it. 


- foo metapackage
- bar1, bar2 and bar3 are the three subordinate packages part of foo. 
- **lite** and **full** are the two profiles we will use


1. Listing all the packages  

`tanzu package available list`

```
NAME                    DISPLAY-NAME          SHORT-DESCRIPTION
  
foo.tanzu.vmware.com    Foo                   Used to explain metapackages
bar1.tanzu.vmware.com   Bar1                  First subordinate package
bar2.tanzu.vmware.com   Bar2                  Second subordinate package
bar3.tanzu.vmware.com   Bar3                  Third subordinate package
```

The subordinate packages as part of each of the profiles. 

| lite  | full  |
| ----- | ----- |
| bar 1 | bar 1 |
| bar 2 | bar 2 |
|       | bar 3 |

2. Installing a profile
* List version information for the package by running:

```
tanzu package available list foo.tanzu.vmware.com 
```
* Create a `foo-values.yaml` file. This configuration file will have configuration required to deploy foo and the subordinate packages. 

Lite Profile
```
profile: lite
bar1_config: value
bar2_config: value 
```
Full Profile

```
profile: full
bar1_config: value
bar2_config: value
bar3_config: value 
```

* Install the profile
```
tanzu package install foo -p foo.tanzu.vmware.com -v 1.0.0 --values-file foo-values.yml
```



### User Stories

#### Story 1
* As a platform installer, I want to be able to install multiple packages with a single command so that I can reduce the toil necessary. 

#### Story 2
* As a platform installer, I want to be able to provide configuration for multiple packages in a single configuration file. 

### Requirements

* Package Repository that sits on a OCI registry.
* Any necessary security is handled by the OCI registry itself. 
* Tanzu CLI with package plug-in
* kapp controller is installed on the target cluster (what version do we need of kapp controller?)


### Implementation Details/Notes/Constraints


* All the three of these packages include their own 
    * metadata
    * versioned package manifest - this includes 
        * version name, 
        * imgpkg bundle(s) and 
        * values schema
    * package configuration 
        * if it is a subordinate package, it deals with its own information
        * if it is a metadata package, this includes references to the subordinate packages and their necessary information. 

* Can you nest metapackages? Yes, but it can get complicated pretty quickly. 

### Security Model

N/A

### Risks and Mitigations

* UX has already been reviewed as this has been adopted by other projects using carvel technology. 

## Compatibility

* tanzu cli with package plugin
* kapp controller version
* TCE-repo package repository

## Alternatives

The alternative today is to install each package individually with it's own values file and it's own upgrade planning. 

## Upgrade Strategy

This is simply a matter of pointing you to a repository that has the new version of all the necessary packages. 

Original Installation

* foo 1.0 which includes 
    * bar1 1.0 and 
    * bar2 1.0

Upgrade 

* foo 1.1 which includes
    * bar1 1.1 and
    * bar2 1.0

We would perform the upgrade of foo, and it will automatically reconcile the subordinate packages to the appropriate version. You don't have to upgrade individual packages. 

It is the upgraded version manifest that points to the validated necessary subordinate packages. 
