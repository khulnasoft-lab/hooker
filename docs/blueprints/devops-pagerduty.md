# Paging DevOps Teams

## Introduction
In this walkthrough, we will setup vulnerability scanning with [Tunnel](https://github.com/khulnasoft/tunnel) and send the results to Hooker for paging DevOps team members for critical vulnerabilities as they are introduced.

## Scenario
A DevOps team would like to configure alerts for scheduled vulnerability scans to notify them about any vulnerable images that they might be running in their clusters. For this they decide to install Tunnel, run it on a schedule and send the results to Hooker.

They decide to configure Hooker so that upon receiving such alerts, Hooker sends an event to PagerDuty which fires off an alert to inform DevOps teams to take necessary action. 

![img.png](assets/tunnel-pagerduty.png)

## Sample Configs
In this case a sample configuration for the components can be described as follows:

### Hooker Config
```yaml
routes:
- name: Tunnel Alerts to Pagerduty
  input: input.report.summary.criticalCount > 0
  actions: [alert-devops]
  template: tunnel-raw-json

# Templates are used to format a message
templates:
- name: tunnel-raw-json
  rego-package: hooker.rawmessage.json

# Actions are target services that should consume the messages
actions:
- name: alert-devops
  type: pagerduty
  enable: true
  pagerduty-auth-token: "<auth token>"
  pagerduty-routing-key: "<routing key>"
```