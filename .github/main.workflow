action "Monitor Pull Requests" {
  uses = "./"
  secrets = ["GITHUB_TOKEN", "JIRA_API_TOKEN", "CONFLICT_ISSUE_STATUS", "JIRA_PROJECT_NAME", "JIRA_USER", "JIRA_HOST"]
}

workflow "Monitor Schedule" {
  on = "schedule(*/15 * * * *)"
  resolves = ["Monitor Pull Requests"]
}

