# Roadmap Triage and Prioritization

The TCE project operates on carefully chosen
[labels](https://github.com/vmware-tanzu/tce/labels),
[milestones](https://github.com/vmware-tanzu/tce/milestones), and
[projects](https://github.com/vmware-tanzu/tce/projects). These enable us to
operate the project appropriately and automatically provide easy reporting to
VMware and our community.

## Triage Process

Every new issue in TCE is marked as
[triage/needs-triage](https://github.com/vmware-tanzu/tce/labels/triage%2Fneeds-triage).
These issues have **not** been assigned a
[milestone](#milestones).
When an issue is triaged, at minimum it is:

- Added to a milestone
- Assigned
  [labels](#labels)
    - minimum: `kind/*` and `area/*`

![triage-flow](../../../img/triage-flow.png)

As a project with multiple decencies, some issues may required opening an
ancillary issue in a dependency's repository. When this takes place, the
original TCE issue will becoming a tracking issue for the dependency's issue.

## Milestones

Milestones associate issues and PRs with a release. Each release (thus
milestone) has a target date associated with it. Any issue that is slated for a
milestone and not assigned to a person is open for contribution. However, we do
recommend preferring those labeled [help
wanted](https://github.com/vmware-tanzu/tce/labels/help%20wanted).

The only milestone that is not associated with a release or date is the
[icebox](https://github.com/vmware-tanzu/tce/milestone/6). This is where all
issues not currently associated with a known milestone are placed. We welcome
[contributions](https://github.com/vmware-tanzu/tce/blob/main/CONTRIBUTING.md')
to help us speed up the work on any icebox'd items. When work begins or a PR is
opened against an item in the icebox, we'll work with the contributor to
determine if we should commit it to another milestone.

When creating a **new** milestone, engineering and product management review
issues in the icebox to determine which should be prioritized.

## Projects

Projects are a 1:1 reflection of a
[milestone](#milestones).
They offer engineering a means to understand what is:

- **Todo:** Planned but not actively worked on
- **In Progress:** Being worked on
    - These issues **must** have an associated PR linked to them

        ![issue association](../../../img/issue-association.png)

- **In Review:** PR is open
    - These PRs must be **not** marked as drafts
- **Completed:** PR has merged and issue is closed

Pull requests should **not** be added to projects, instead they are linked to
issues in the project.

## Labels

Each label has a prefix/suffix where the prefix is a category for the suffix.
The prefixes and their meaning are as follows:

- **triage/*:** This issue is not being worked on and requires additional
  consideration before moving forward.
- **kind/*:** The type of issue being proposed. This weighs into prioritization
  as some issues may be bugs, enhancements to existing functionality, or new
functionality altogether.
    - Every triaged issue **must** include a kind.
- **area/*:** The area(s) of the codebase this impacts. This metadata is used to
  determine who is best to assign an issue.
    - Every triaged issue **must** include an area.

## Roadmap

Our roadmap is captured in a [long-running
issue](https://github.com/vmware-tanzu/tce/issues/1293). This issue is an
aggregate of one or many milestone(s) to articulate to the community where we
are headed and what we have done.

## Common Questions and Answers:

- **Q:** How do I know what is slated for a release?
    - **A:** There are multiple ways:
        - View the [roadmap](https://github.com/vmware-tanzu/tce/issues/1293).
        - View the corresponding milestone ([for example,
          v1.0.0](https://github.com/vmware-tanzu/tce/milestone/5)).
- **Q:** How do I view the progress on an in-flight release?
    - **A:** View the corresponding project for that release/milestone ([for
      example, v1.0.0](https://github.com/vmware-tanzu/tce/projects/11)).
- **Q:** How do I view bugs specifically being targeting for a milestone?
    - **A:** Create an issue filter for `kind/bug` on the given milestone ([for
      example, bugs fixed in
v1.0.0](https://github.com/vmware-tanzu/tce/issues?q=label%3Akind%2Fbug+milestone%3Av1.0.0+)).
- **Q**: How do I view issues relating to documentation for a milestone?
    - **A:** Create an issue filter for `area/docs` on the given milestone ([for
      example, docs issues in
v1.0.0](https://github.com/vmware-tanzu/tce/issues?q=label%3Aarea%2Fdocs+milestone%3Av1.0.0+)).
