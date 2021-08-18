# Security Release Process

Tanzu Community Edition is a growing community of users and developers for a free open source distribution of Kubernetes and other cloud native technologies that is easy to install, easy to use, and easy to manage. The community has adopted this security disclosure and response policy to ensure we responsibly handle critical issues.

## Supported Versions

The Tanzu Community Edition project maintains the following [document on the release process and support matrix](./RELEASES.md). Please refer to it for release related details.

## Reporting a Vulnerability - Private Disclosure Process

Security is of the highest importance and all security vulnerabilities or suspected security vulnerabilities should be reported to Tanzu Community Edition privately, to minimize attacks against current users before they are fixed. Vulnerabilities will be investigated and patched on the next patch (or minor) release as soon as possible. This information could be kept entirely internal to the project.  

If you know of a publicly disclosed security vulnerability for Tanzu Community Edition or any of its components, please **IMMEDIATELY** [contact](https://github.com/vmware-tanzu/community-edition/security/policy#mailing-lists) the Tanzu Community Edition Security Team.
 **IMPORTANT: Do not file public issues on GitHub for security vulnerabilities**

To report a vulnerability or a security-related issue, please contact the [private email address](https://github.com/vmware-tanzu/community-edition/security/policy#mailing-lists) with the details of the vulnerability. The email will be fielded by the Tanzu Community Edition Security Team, which is made up of project maintainers who have committer and release permissions. Emails will be addressed within 3 business days, including a detailed plan to investigate the issue and any potential workarounds to perform in the meantime. Do not report non-security-impacting bugs through this channel. Use [GitHub issues](https://github.com/vmware-tanzu/community-edition/issues/new/choose) instead.

### Proposed Email Content

Provide a descriptive subject line and in the body of the email include the following information:

* Basic identity information, such as your name and your affiliation or company.
* Detailed steps to reproduce the vulnerability  (POC scripts, screenshots, and compressed packet captures are all helpful to us).
* Description of the effects of the vulnerability on Tanzu Community Edition (or its incorporated projects) and the related hardware and software configurations, so that the Security Team can reproduce it.
* How the vulnerability affects Tanzu Community Edition usage and an estimation of the attack surface, if there is one.
* List other projects or dependencies that were used in conjunction with Tanzu Community Edition to produce the vulnerability.

## When to report a vulnerability

* When you think Tanzu Community Edition has a potential security vulnerability.
* When you suspect a potential vulnerability but you are unsure that it impacts Tanzu Community Edition.
* When you know of or suspect a potential vulnerability on another project that is used by Tanzu Community Edition. For example, Community Edition incorporates Kubernetes.
  
## Patch, Release, and Disclosure

The Tanzu Community Edition Security Team will respond to vulnerability reports as follows:

1. The Security Team will investigate the vulnerability and determine its effects and criticality.
2. If the issue is not deemed to be a vulnerability, the Security Team will follow up with a detailed reason for rejection.
3. The Security Team will initiate a conversation with the reporter within 3 business days.
4. If a vulnerability is acknowledged and the timeline for a fix is determined, the Security Team will work on a plan to communicate with the appropriate community, including identifying mitigating steps that affected users can take to protect themselves until the fix is rolled out.
5. The Security Team will also create a [CVSS](https://www.first.org/cvss/specification-document) using the [CVSS Calculator](https://www.first.org/cvss/calculator/3.0). The Security Team makes the final call on the calculated CVSS; it is better to move quickly than making the CVSS perfect. Issues may also be reported to [Mitre](https://cve.mitre.org/) using this [scoring calculator](https://nvd.nist.gov/vuln-metrics/cvss/v3-calculator). The CVE will initially be set to private.
6. The Security Team will work on fixing the vulnerability and perform internal testing before preparing to roll out the fix.
7. A public disclosure date is negotiated by the Tanzu Community Edition Security Team and the bug submitter. We prefer to fully disclose the bug as soon as possible once a user mitigation or patch is available. It is reasonable to delay disclosure when the bug or the fix is not yet fully understood or the solution is not well-tested. The timeframe for disclosure is from immediate (especially if it’s already publicly known) to a few weeks. For a critical vulnerability with a straightforward mitigation, we expect the time from the report date to public disclosure date to be on the order of 14 business days. The Tanzu Community Edition Security Team holds the final say when setting a public disclosure date.
8. Once the fix is confirmed, the Security Team will patch the vulnerability in the next patch or minor release, and backport a patch release into all earlier supported releases. Upon release of the patched version of Tanzu Community Edition, we will follow the **Public Disclosure Process**.

### Public Disclosure Process

The Security Team publishes a public [advisory](https://github.com/vmware-tanzu/community-edition/security/advisories) to the Tanzu Community Edition community via GitHub. In most cases, additional communication via Slack, Twitter, mailing lists, blog and other channels will assist in educating users and rolling out the patched release to affected users.

The Security Team will also publish any mitigating steps users can take until the fix can be applied to their Tanzu Community Edition instances.  

## Mailing lists

* Use tanzu-ce-maintainers@lists.vmware.com to report security concerns to the Tanzu Community Edition Security Team, who uses the list to privately discuss security issues and fixes prior to disclosure.

## Confidentiality, integrity and availability

We consider vulnerabilities leading to the compromise of data confidentiality, elevation of privilege, or integrity to be our highest priority concerns. Availability, in particular in areas relating to DoS and resource exhaustion, is also a serious security concern. The Tanzu Community Edition Security Team takes all vulnerabilities, potential vulnerabilities, and suspected vulnerabilities seriously and will investigate them in an urgent and expeditious manner.

Note that we do not currently consider the default settings for Tanzu Community Edition to be secure-by-default. It is necessary for operators to explicitly configure settings, role based access control, and other resource related features in Tanzu Community Edition to provide a hardened environment. We will not act on any security disclosure that relates to a lack of safe defaults. Over time, we will work towards improved safe-by-default configuration, taking into account backwards compatibility.
