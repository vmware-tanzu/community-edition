<!-- The sections in this template are highly-suggested, but optional. They exist to
provide an idea of what an in-depth proposal **might** include. If a different
format fits your proposal, feel free to use it, with the expectation that:

1. You must link to the proposal issue (github) under the title
1. If you omit sections that are relevant, you may be asked by maintainers to
   add them in order to best review the proposal. -->

# Title

Keep it simple and descriptive.

* Proposal: < GITHUB ISSUE URL >
* See Also:
  * < RELEVANT URL >
  * < RELEVANT URL >

## Summary

Provide a short summary of what is being proposed. The `Summary` section is
incredibly important for producing high quality user-focused documentation such
as release notes or a development roadmap.  It should be possible to collect
this information before implementation begins in order to avoid requiring
implementers to split their attention between writing release notes and
implementing the feature itself.

A good summary is probably at least a paragraph in length.

## Motivation

This section is for explicitly listing the motivation, goals and non-goals of
this proposal. To support what is being proposed, provide details on why this
change is important. If user stories are relevant to this proposal, add them
below.

* Describe why the change is important and the benefits to users.  

### Goals

* List the specific high-level goals of the proposal.
* How will we know that this has succeeded?

### Non-Goals/Future Work

* What high-levels are out of scope for this proposal?
* Listing non-goals helps to focus discussion and make progress.

## Proposal

This is where we get down to the nitty gritty of what the proposal actually is.
Provide a detailed proposal of **what** would be done and **how**; feel free to
add multiple subsections, diagrams, and config/code snippets

* What is the plan for implementing this feature?
* What data model changes, additions, or removals are required?
* Provide a scenario, or example.
* Use diagrams to communicate concepts, flows of execution, and states.

### User Stories

* Detail the things that people will be able to do if this proposal is implemented.
* Include as much detail as possible so that people can understand the "how" of the system.
* The goal here is to make this feel real for users without getting bogged down.

#### Story 1

#### Story 2

### Requirements

Some authors may wish to use requirements in addition to user stories.
Technical requirements should derived from user stories, and provide a trace from
use case to design, implementation and test case. Requirements can be prioritised
using the MoSCoW (MUST, SHOULD, COULD, WON'T) criteria.

The FR and NFR notation is intended to be used as cross-references across a CAEP.

The difference between goals and requirements is that between an executive summary
and the body of a document. Each requirement should be in support of a goal,
but narrowly scoped in a way that is verifiable or ideally - testable.

#### Functional Requirements

Functional requirements are the properties that this design should include.

##### FR1

##### FR2

#### Non-Functional Requirements

Non-functional requirements are user expectations of the solution. Include
considerations for performance, reliability and security.

##### NFR1

##### NFR2

### Implementation Details/Notes/Constraints

* What are some important details that didn't come across above.
* What are the caveats to the implementation?
* Go in to as much detail as necessary here.
* Talk about core concepts and how they relate.

### Security Model

Document the intended security model for the proposal, including implications
on the Kubernetes RBAC model. Questions you may want to answer include:

* Does this proposal implement security controls or require the need to do so?
  * If so, consider describing the different roles and permissions with tables.
* Are their adequate security warnings where appropriate.
* Are regex expressions going to be used, and are their appropriate defenses
  against DOS.
* Is any sensitive data being stored in a secret, and only exists for as long as necessary?

### Risks and Mitigations

* What are the risks of this proposal and how do we mitigate? Think broadly.
* How will UX be reviewed and by whom?
* How will security be reviewed and by whom?
* Consider including folks that also work outside the SIG or subproject.

## Compatibility

If this change impacts compatibility of previous versions of TCE or software
integrated with TCE, please call it out here. If incompatibilities can be
mitigated, please add it here.

## Alternatives

The `Alternatives` section is used to highlight and record other possible
approaches to delivering the value proposed by a proposal.

## Upgrade Strategy

If applicable, how will the component be upgraded? Make sure this is in the test
plan.

Consider the following in developing an upgrade strategy for this enhancement:

* What changes (in invocations, configurations, API use, etc.) is an existing
  cluster required to make on upgrade in order to keep previous behavior?

## Additional Details

### Test Plan

**Note:** *Section not required until targeted at a release.*

Consider the following in developing a test plan for this enhancement:

* Will there be e2e and integration tests, in addition to unit tests?
* How will it be tested in isolation vs with other components?

No need to outline all of the test cases, just the general strategy.  Anything
that would count as tricky in the implementation and anything particularly
challenging to test should be called out.

All code is expected to have adequate tests (eventually with coverage
expectations).

### Graduation Criteria

**Note:** *Section not required until targeted at a release.*

Define graduation milestones.

These may be defined in terms of API maturity, or as something else. Initial
proposal should keep this high-level with a focus on what signals will be looked
at to determine graduation.

Consider the following in developing the graduation criteria for this
enhancement:

* [Maturity levels (`alpha`, `beta`, `stable`)][maturity-levels]
* [Deprecation policy][deprecation-policy]

Clearly define what graduation means by either linking to the [API doc
definition](https://kubernetes.io/docs/concepts/overview/kubernetes-api/#api-versioning),
or by redefining what graduation means.

In general, we try to use the same stages (alpha, beta, GA), regardless how the
functionality is accessed.

[maturity-levels]: https://git.k8s.io/community/contributors/devel/sig-architecture/api_changes.md#alpha-beta-and-stable-versions
[deprecation-policy]: https://kubernetes.io/docs/reference/using-api/deprecation-policy/

### Version Skew Strategy

If applicable, how will the component handle version skew with other components?
What are the guarantees? Make sure this is in the test plan.

<!-- Links -->
[community meeting]: https://hackmd.io/CiuO4V0AT6WL_TgA47MXBA
