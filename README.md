# pkgfy - let's pretend it's a package

Install golang binaries from a URL, e.g. GitHub release, as if it was a package from a package
manager.

## Milestones

### 1 - Support Install

Given a URL, install the binary from it.

Must handle compressed and non-compressed formats.

Must "remember" what was installed, and expose command to view what is installed.

Installation destination is `~/bin` by default, but parameter or configuration file can specify
a different location.

"db" of installed packages is located at `~/.config/pkgfy/installed.db` by default.
Parameter or configuration file can specify a different location.

Configuration file is located at `~/.config/pkgfy/pkgfy.yaml` by default.
Parameter can specify a different location.

### 2 - Support Upgrades

Given a package that was installed from GitHub, check git repository for updates.
Offer update option.

### 3 - Introduce pkgfy Repositories

Create a YAML file to define a set of packages that can be installed.
Support digests and signatures for installation verification.
