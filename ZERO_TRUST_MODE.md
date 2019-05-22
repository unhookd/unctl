# Deploying an application with Unhookd Zero Trust
## Overview
The zero trust mode of Unhookd requires that an application be registered with Unhookd in order to deploy. This ensures that:

- An application is explicitly authorized to be deployed
- The application's configuration as specified by a `values.yaml` file has been code reviewed and checked in
- The application is green in CI, the latest valid sha of a given branch, and being deployed to the correct cluster.

It's called `zero trust` because the intent is to ensure that there is a trail from application configuration, to source code change, to tests passing, to deploy.

## Prerequisites
- All applications being deployed with Unhookd zero trust mode *must* have a configured umbrella chart. For more information on how to create and use an umbrella chart, see the [charts](https://github.com/org/charts/tree/master/README.md) readme.

## Configuring an Unhookd deployed application
Deployments parameters are defined in the `lookup-config.yaml` in this repository. An application requires the following parameters:

```
base-chart-name:
  umbrella-chart-name:
    release: # the name the applications helm release will have
    namespace: # the namespace the application will be deployed to
    repo: # the repository the application lives in
    cluster: # the cluster the application should be deployed to
    branch: # the branch that deploys the application
    chart: # the chart to be used to deploy the application
    notifications: [] # and optional array of notification channels and their parameters
```

Below is an example of what that code looks like for our sample repo.

```
kube-hello-world:
  kube-hello-world-production:
    release: kube-hello-world
    namespace: kube-hello-world
    repo: org/kube-hello-world
    cluster: production
    branch: master
    chart: org/kube-hello-world-production
    notifications:
      - provider: slack
        channel: "#ops-notes"
        text: "Deployed kube hello world for production"
      - provider: slack
        channel: "#ops-notes"
        text: "Deployed kube hello world for production yet again"
```

## Configuring notifications for Unhookd

Unhookd supports sending a notification to a given notification channel upon deploy. Currently only `slack` is supported.

### Slack
A Slack notification is configured like this:

```
notifications:
  - provider: slack
    channel: "#my-slack-channel"
    text: "The message I want to send to Slack"
```

Once configured a slack message will be sent on a successful deploy.

### Adding a deploy job

A deploy job should make the correct request to Unhookd to kick off the deploy. This will only work if your application has been configured correctly in the `lookup-config.yaml`.

A deploy job is pretty simple:
- Set your docker image to be the latest sha of `unhookd`. Feel free to check another Unhookd configured repo to see what that looks like.
- Run the `unhookd deploy` command with your `application` and `release` and the `sha` you are deploying. The sha comes from the Circle CI environment.

A deploy step might look something like this:

```
  deploy-master:
    docker:
      - image: 673102273038.dkr.ecr.us-west-2.amazonaws.com/unhookd:git-52da8b5d2f7b5d6dd6f997d45b6dbd04a021c250
    steps:
      - run:
          name: deploy
          command: |
            unhookd deploy kube-hello-world kube-hello-world-engineering "${CIRCLE_SHA1}"
```
