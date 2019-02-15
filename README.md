azmon
=====

A tool for Azure Monitor at Microsoft Azure that possible to collects metrics, and checks (as Nagios plugin) it.


## Usage

```bash
$ azmon <global options> check --metric-name "Percentage CPU" --warning-over 70 --critical-over 90
CRITICAL - <resource name> Percentage CPU is 95.885000 that over than 90.000000
```

```bash
$ azmon <global options> metric --metric-names "Percentage CPU,Network In,Network Out,Disk Read Bytes"
azure.Microsoft.ComputevirtualMachines.<resource name>.PercentageCPU.Average     5.932500        1550223420
azure.Microsoft.ComputevirtualMachines.<resource name>.NetworkIn.Average         37235.038462    1550223420
azure.Microsoft.ComputevirtualMachines.<resource name>.NetworkOut.Average        5743.250000     1550223420
azure.Microsoft.ComputevirtualMachines.<resource name>.DiskReadBytes.Average     0.000000        1550223420
```

### Global options

- `--subscriptionID`
    - Set the subscription id 
- `--resource-group`
    - Set the resource group name
- `--namespace`
    - Set the metric namespace
- `--resource`
    - Set the target resource name
- `--aggregation`
    - Set the aggregation type. Choose from "Total", "Average", "Maximum", "Minimum" ("Count" is not supported)"
- `--auth-file`
    - Set the azure auth file path (default "/etc/nagios/azure.auth")
    - See also [Authentication](#authentication)

### Subcommands

#### check

`check` sub-command fetches metric data and check it as nagios plugin.  

Options  

- `--metric-name`
    - Set the name of the metric
- `--warning-over`
    - Set the warning threshold. Occur warning level alert when metric data over than threshold 
- `--warning-under`
    - Set the warning threshold. Occur warning level alert when metric data under than threshold
- `--critical-over`
    - Set the critical threshold. Occur critical level alert when metric data over than threshold
- `--critical-under`
    - Set the critical threshold. Occur critical level alert when metric data under than threshold 

#### metric

`metric` sub-command fetches metric data and print it format sensu plugin.

Options  

- `--metric-names`
    - Set the names of the metric
- `--prefix`
    - Set the metric key prefix (default "azure")


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
