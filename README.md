# prwatch-action

A Github action for monitoring pull requests on a repository.

Current features
- Monitor the merability of all pull requests against master
- Transition Jira cards associated with in-conflict pull requests to a new status defined by `CONFLICT_ISSUE_STATUS`
- Comment on transitioned issues to let assignees know their cards have been "pushed back" due to conflicts in the pull
  request.

# Usage

To use this action, your Github Pull Requests must include in their description their associated issue tracker ID. E.g.
if your Jira project name is `FOO` and the issue associated with your pull request is `1234`, then your Pull Request
must include `FOO-1234` somewhere in its description.

## Example Pull Request Description
```
This PR fixes the Thinger for FOO-1234
```

## Run on merge with master
```yaml
---
'on':
  push:
      branches:
            - master

name: PRWatch Action
jobs:
  check:
    name: Check Pull Requests
    runs-on: ubuntu-latest
    steps:
      - name: Check for conflicts
        uses: acaloiaro/prwatch-action@latest
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          JIRA_API_TOKEN: ${{ secrets.JIRA_API_TOKEN }}
```

## Run every 15 minutes
```yaml
---
'on':
  schedule:
    - cron: '*/15 * * * *'
name: Monitor Pull Requests
jobs:
  check:
    name: Check Pull Requests
    runs-on: ubuntu-latest
    steps:
      - name: Check for conflicts
        uses: acaloiaro/prwatch-action@latest
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          JIRA_API_TOKEN: ${{ secrets.JIRA_API_TOKEN }}
```

## Configuration File

This action is configured with a single yaml file. The configuration file lives in your repository at
`./github-actions/prwatch-action/config.yaml`. See [examples](https://github.com/acaloiaro/prwatch-action/tree/master/examples) for an example `config.yaml`.

| key           | description                                                       | type | default |
| ------------- |:-----------------------------------------------------------------:|:----:|:--------|
| settings.dual_pass.enabled  | Dual-pass mode allows this action to be triggered on 'push' to a target branch while allowing Github time to recalculate the mergeability of PRs | bool | true |
| settings.dual_pass.wait_duration | The duration of time to wait between the first and second pass in dual pass mode. This period of time should be long enough for Github to determine the mergeability of all your open pull requests. e.g. `1m30s`. Note: The value of this variable must conform to the Golang duration format: https://golang.org/pkg/time/#ParseDuration | time | 60s |
| settings.issues.enable_comment | When merge conflicts occurr, comment on associated issues | bool | true |
| settings.issues.enable_transition | When merge conflicts occur, transition associated issues to new status | bool | true |
| settings.issues.conflict_status | When merge conflicts occur, the new issue status to transitions issues to | string | |
| settings.jira.enabled | Use Jira as your issue tracker | bool | true |
| settings.jira.host | The hostname of your Jira instance | string | |
| settings.jira.project_name | The name of the Jira project associated with your repository | string | |
| settings.jira.user | The "bot" user to use when transitioning and commenting on issues | string | |

## Secrets
`GITHUB_TOKEN`: _It is not necessary to set this, as it is available to all Github Actions_

`JIRA_API_TOKEN`: The access token associated with `settings.jira.user`.
