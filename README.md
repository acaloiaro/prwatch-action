# prwatch
Github action for monitoring pull requests on a repository.

Current features
* Monitor the merability of all pull requests against master
* Transition Jira cards associated with in-conflict pull requests to a new status defined by `CONFLICT_ISSUE_STATUS`
* Comment on transitioned cards to let assignees know that their cards have been "pushed back" due to PR conflict

# Usage

Note: While it is tempting to run this action on "push" to your master branch, doing so will be quite ineffective. The
reason for that is because, immediately following a merge to "master", it is impossible to determine the mergability of the
pull requests against it.

For that reason, it's better to run on a schedule.

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
