# k8s-dotenv

A commandline tool to fetch, merge and convert secrets and config maps in kubernetes to a dot env file.

## Install

Download from the [Releases](https://github.com/eiladin/k8s-dotenv/releases) page, extract and put in PATH.  

Alternatively, use [install.sh](https://raw.githubusercontent.com/eiladin/k8s-dotenv/main/install.sh) to download and extract the latest version automatically.

[Documentation](./docs/k8s-dotenv.md)

## Supported Resource Types
- cronjob
- deployment
- daemonset
- job
- pod
- statefulset

## Usage
```bash
k8s-dotenv get <resource_type> <RESOURCE_NAME>
```

## Examples

### Get Deployment and write to .env
```bash
k8s-dotenv get deploy my-deployment
```
### Get DaemonSet and write to out.txt file
```bash
k8s-dotenv get ds my-daemonset -f out.txt
```
### Get Job and write output to console (stdout)
```bash
k8s-dotenv get job my-job -c
```

## Help
```bash
k8s-dotenv --help
```

## Shell Completions

k8s-dotenv comes with shell completions for bash, zsh, fish and powershell built-in.  After install, run `k8s-dotenv completion --help` for information.