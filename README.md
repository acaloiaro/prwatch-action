# prwatch-action [ archived ] 

## This repo is archived. The functionality it provides is largely provided by Github's new "real-time alerts". Go to https://github.com/settings/reminders and enable real-time alerts, and check "Your PR has merge conflicts". The code will remain on Github for future refernce. 

A Github action for monitoring pull requests on a repository.

Supported features
- Monitor the mergeability of all open pull requests in your repository
- When pull requests have conflicts, comment on them and `@mention` the owner
- When pull requests have conflicts, transition them to new statuses, e.g. 'To Be Shipped' -> 'In Progress'
- Configure globally for the entire repository or on a per-user basis

# Usage

To use this action, your Github Pull Requests must include in their description an associated issue tracker issue ID.
E.g. if your Jira project name is `FOO` and the issue associated with your pull request is `1234`, then your Pull
Request must include `FOO-1234` somewhere in its description.

## Example Pull Request Description
```
This PR fixes the Thinger for FOO-1234
```

## Setup

This action performs best when it is configured with `settings.dual_pass.enabled` and `on.push` to your mainline branch.
The idea is to trigger a mergeability check when feature branches merge upstream. See [the configuration
section](#configuration_file) for how to enable dual-pass mode when triggering this action `on.push`.

1. Add `./.github/workflows/prwatch.yml` workflow file to your repository [Example
   Workflows](https://github.com/acaloiaro/prwatch-action/tree/master/examples/workflows)
2. Add `./github-actions/prwatch-action/config.yaml` to your repository. [Example
   configuration](https://github.com/acaloiaro/prwatch-action/tree/master/examples/config.yaml)

## Run on Push
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
      - name: Checkout Branch
        uses: actions/checkout@v1
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

name: PRWatch Action
jobs:
  check:
    name: Check Pull Requests
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Branch
        uses: actions/checkout@v1
      - name: Check for conflicts
        uses: acaloiaro/prwatch-action@latest
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          JIRA_API_TOKEN: ${{ secrets.JIRA_API_TOKEN }}
```

## <a name="configuration_file"></a>Configuration File

This action is configured with a single yaml file. The configuration file lives in your repository at
`./github-actions/prwatch-action/config.yaml`. See [examples](https://github.com/acaloiaro/prwatch-action/tree/master/examples) for an example `config.yaml`.

User-specific settings are superseded by their global counterparts. i.e. if
`users.foobar.settings.issues.enable_transition` is _on_ for `foobar`, but _off_ globally, the feature will be turned
off for `foobar`.

If a global setting is unset globally, but on for a user, then it is still on for that user.

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
| users.`<github_username>`.settings.issues.enable_comment | Enable issue comments for a user | bool | |
| users.`<github_username>`.settings.issues.enable_transition | Enable issue transitions for a user | bool | |

## Secrets
`GITHUB_TOKEN`: _It is not necessary to set this, as it is available to all Github Actions_

`JIRA_API_TOKEN`: The access token associated with `settings.jira.user`.
