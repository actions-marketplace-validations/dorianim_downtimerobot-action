<h1 align="center">
    Downtimerobot-action
</h1>

<p align="center">
    <a href="https://www.gnu.org/licenses/agpl-3.0">
        <img src="https://img.shields.io/badge/License-AGPL%20v3-blue.svg" />
    </a>
    <a href="https://goreportcard.com/report/github.com/dorianim/downtimerobot-action">
        <img src="https://goreportcard.com/badge/github.com/dorianim/downtimerobot-action" alt="Go report" />
    </a>
</p>

# Terminology
- `Service`: A monitored website or other service

# TODO
- Implement other service types:
  - pattern
  - ping
  - port

# Environment
- `GITHUB_ACTIONS` -> If true, github mode will be used
- `GITHUB_REPOSITORY`, `GITHUB_REF_NAME` -> Will be used for api basepath in github mode: ` https://raw.githubusercontent.com/$GITHUB_REPOSITORY/$GITHUB_REF_NAME/public/data/*.json`