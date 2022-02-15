# k8s-dotenv

A commandline tool to fetch, merge and convert secrets and config maps in kubernetes to a dot env file.

## Install

Download from the [Releases](https://github.com/eiladin/k8s-dotenv/releases) page, extract and put in PATH.  

Alternatively, use [install.sh](https://raw.githubusercontent.com/eiladin/k8s-dotenv/blob/main/install.sh) to download and extract the latest version automatically.

## Usage
### CronJob
```bash
k8s-dotenv get cj <DAEMONSET_NAME>
```

### Deployment
```bash
k8s-dotenv get deploy <DEPLOYMENT_NAME>
```
### DaemonSet
```bash
k8s-dotenv get ds <DAEMONSET_NAME>
```
### Job
```bash
k8s-dotenv get job <DAEMONSET_NAME>
```

## Help
```bash
k8s-dotenv --help
```

## Shell Completions

k8s-dotenv comes with shell completions for bash, zsh, fish and powershell built-in.  After install, run `k8s-dotenv completion --help` for information.