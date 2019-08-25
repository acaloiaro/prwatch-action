# prwatch
Github action for monitoring pull requests on a repository.

Current features
- Monitor the merability of all pull requests against master
- Transition Jira cards associated with in-conflict pull requests to a new status defined by `CONFLICT_ISSUE_STATUS`
- Comment on transitioned issues to let assignees know their cards have been "pushed back" due to conflicts in the pull
  request.

# Usage

To use this action, your Github Pull Requests must include in their description the key to any associated issues. E.g.
if your project name is `FOO` and the issue associated with your pull request is `1234`, then your Pull Request must
include `FOO-1234` somewhere in its descripton.

## Example Pull Request Description
```
This PR fixes the Thinger for FOO-1234
```

Note: While it is tempting to run this action on "push" to your master branch, doing so is quite ineffective. That's
because, immediately following a merge to "master", Github cannot determine the mergability of any pull requests
against it.

As a result, this action is best run on a schedule.

## Example: every 15 minutes
```yaml
---
'on':
  schedule:
    - cron: '*/15 * * * *'
name: Monitor Pull Requests
jobs:
  monitor:
    name: Monitor
    runs-on: ubuntu-latest
    steps:
      - name: Monitor
        uses: acaloiaro/prwatch@master
        env:
          CONFLICT_ISSUE_STATUS: In Development
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          JIRA_API_TOKEN: ${{ secrets.JIRA_API_TOKEN }}
          JIRA_HOST: companyname.atalassian.net
          JIRA_PROJECT_NAME: PROJNAME
          JIRA_USER: jira-bot
```

## Variables
`CONFLICT_ISSUE_STATUS`: The new status to assign issues when their corresponding PRs are in conflict

`JIRA_HOST`: The hostname for your Jira instance. If you are on Jira Cloud, this will be `companyname.atalassian.net`

`JIRA_PROJECT_NAME`: The name of the Jira project associated with the repository where this action is installed

`JIRA_USER`: The jira user with which to perform API requests

## Secrets
`GITHUB_TOKEN`: It is not necessary to set this, as it is available to all Github Actions

`JIRA_API_TOKEN`: The access token to authenticate `JIRA_USER` with your Jira instance
