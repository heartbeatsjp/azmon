azmon
=====

A tool for Azure Monitor at Microsoft Azure that possible to collects metrics, and checks(as Nagios plugin) it.


## Usage

### Global options

- `--subscriptionID`
    - Set the subscription id 
- `--resource-group`
    - Set the resource group name
- `--namespace`
    - Set the metric namespace
- `--resource`
    - Set the target resource name
- `--metric-name`
    - Set the name of the metric
- `--aggregation`
    - Set the aggregation type. Choose from "Total", "Average", "Maximum", "Minimum" ("Count" is not supported)"
- `--auth-file`
    - Set the azure auth file path (default "/etc/nagios/azure.auth")
    - See also [Authentication](#authentication)

### Subcommands

`check`

```
azmon <global options> check --warning <warning threshold> --critical <critical threshold>
```

`metrics`

```
azmon <global options> metrics [--prefix <metric prefix key>]
```


### Authentication

azmon fetches metric data from Azure API, that required authentication with Azure API.    

[azure-sdk-for-go](https://github.com/Azure/azure-sdk-for-go) used in azmon provides several authentication methods.  
azmon supports file-based authentication.        
https://docs.microsoft.com/en-us/go/azure/azure-sdk-go-authorization#use-file-based-authentication  

Generate auth file can run below command.   

```bash
az ad sp create-for-rbac --sdk-auth > azure.auth
```

In the default, azmon reads `/etc/nagios/azure.auth` as auth file. You can change the auth file path by `--auth-file` option.
