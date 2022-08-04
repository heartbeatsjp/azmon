azmon
=====

[![wercker status](https://app.wercker.com/status/c734acd73bb51fb728351ac8be88d83d/s/master "wercker status")](https://app.wercker.com/project/byKey/c734acd73bb51fb728351ac8be88d83d)
[![Go Report Card](https://goreportcard.com/badge/github.com/heartbeatsjp/azmon)](https://goreportcard.com/report/github.com/heartbeatsjp/azmon)

A tool for Azure Monitor at Microsoft Azure that possible to collect metrics and check it.  


## Usage

Azmon has two sub-commands named `check` and `metric`. When you invoke sub-commands, must be specified global options (see [Global options](#global-options)).  

`azmon check` can check metric data as Nagios plugin. You specify the target metric name in `--metric-name` option.  
Also, `azmon check` provides options `--*-over` and `--*-under`, because whether we should check "over than threshold" or "under than threshold" is different by type of metric data.  

```bash
$ azmon <global options> check --metric-name "Percentage CPU" --warning-over 70 --critical-over 90
CRITICAL - <resource name> Percentage CPU is 95.885000 that over than 90.000000
```

`azmon metric` can print metric data in the Sensu plugin format.  
The `--metric-names` option can specify one or more metric data with comma separated. When you invoke without `--metric-names` option, target metric data is all of the metric data contained namespace.  

```bash
$ azmon <global options> metric --metric-names "Percentage CPU,Network In,Network Out,Disk Read Bytes"
azure.Microsoft.ComputevirtualMachines.<resouce group>.<resource name>.PercentageCPU.Average     5.932500        1550223420
azure.Microsoft.ComputevirtualMachines.<resouce group>.<resource name>.NetworkIn.Average         37235.038462    1550223420
azure.Microsoft.ComputevirtualMachines.<resouce group>.<resource name>.NetworkOut.Average        5743.250000     1550223420
azure.Microsoft.ComputevirtualMachines.<resouce group>.<resource name>.DiskReadBytes.Average     0.000000        1550223420
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
- `--start-time`
    - Set the start time as unix timestamp, relative from now if 0 or negative value given (default -300)
- `--end-time`
    - Set the end time as unix timestamp, relative from now if 0 or negative value given (default 0)
- `interval-sec`
    - Set the interval seconds (supported ones are: 60,300,900,1800,3600,21600,43200,86400. default 60)
### Subcommands

#### check

`check` sub-command fetches metric data and check it as Nagios plugin.  

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

`metric` sub-command fetches metric data and print it in the Sensu plugin format.  

Options  

- `--metric-names`
    - Set the names of the metric
    - When omitting this option, print metric data contains in metric namespace specified by the global option
        - However except metrics that required dimension
- `--prefix`
    - Set the metric key prefix (default "azure")


### Authentication

Azmon fetches metric data from Azure API, that required authentication with Azure API.  

[azure-sdk-for-go](https://github.com/Azure/azure-sdk-for-go) used in azmon provides several authentication methods.  
Currently, azmon supports **file-based authentication** only.  
https://docs.microsoft.com/en-us/go/azure/azure-sdk-go-authorization#use-file-based-authentication  

To generate an authentication file, you can run the following command.  

```bash
az ad sp create-for-rbac --sdk-auth > azure.auth
```

In the default, azmon reads `/etc/nagios/azure.auth` as authentication file. You can change the authentication file path by `--auth-file` option.  


## License

[Apache License 2.0](https://github.com/heartbeatsjp/azmon/blob/master/LICENSE)  
