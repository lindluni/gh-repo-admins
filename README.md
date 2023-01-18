# Repo Admins

`gh-repo-admins` is a GitHub CLI extension used to retrieve a list of repository admins. This is useful when attempting
to respond to security incidents or to determine who to make administrative requests to.

**Note:** The extension will retrieve the username, public full name, and public email address of all repository admins. 
If a users name or email are not public, only the username will be returned.

**Note:** You must have at least `read` access to the repository in order to use this extension.

## Pre-requisites

- [GitHub CLI](https://cli.github.com/)

## Installation

```bash
gh extension install lindluni/gh-repo-admins
```

## Usage

```bash
NAME:
   repo-admins - query repository admins

USAGE:
   repo-admins [global options] command [command options] [arguments...]

VERSION:
   1.1.0

DESCRIPTION:
   gh repo-admins --owner [owner] --repo [repo]
   gh repo-admins --help

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --owner value, -o value  organization or user that owns the repository
   --repo value, -r value   repository name
   --file value, -f value   name of CSV output file (default: "repo-admins.csv")
   --delay value, -d value  delay between GitHub API requests in milliseconds. If you are hitting API rate limits, increase this value. (default: 500)
   --help, -h               show help (default: false)
   --version, -v            print the version (default: false)
```

## Examples

```bash
# Retrieve a list of repository admins for the lindluni/gh-repo-admins repository
gh repo-admins --owner lindluni --repo gh-repo-admins

# Retrieve a list of repository admins for the lindluni/gh-repo-admins repository and save the results to a different CSV file
gh repo-admins --owner lindluni --repo gh-repo-admins --file repo-admins.csv

# Retrieve a list of repository admins for the lindluni/gh-repo-admins repository and increase the delay between GitHub API requests
gh repo-admins --owner lindluni --repo gh-repo-admins --delay 1000
```