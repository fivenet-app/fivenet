---
title: FiveNet Helm Chart
---
{{ template "generatedDocsWarning" . }}

Installs [FiveNet](https://github.com/galexrt/fivenet).

## Prerequisites

* Kubernetes 1.19+
* Helm 3.x

See the [Helm support matrix](https://helm.sh/docs/topics/version_skew/) for more details.

## Installing

```console
helm repo add fivenet https://galexrt.moe/fivenet/
helm install fivenet/fivenet -f values.yaml

# Or upgrade
helm upgrade --install fivenet/fivenet -f values.yaml
```

For example settings, see the next section or [values.yaml](/charts/fivenet/values.yaml).

## Configuration

The following table lists the configurable parameters of the FiveNet chart and their default values.

{{ template "chart.valuesTable" . }}

## Uninstalling the Chart

To see the currently installed Rook chart:

```console
helm ls
```

To uninstall/delete the `fivenet` deployment:

```console
helm delete fivenet
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## License

Apache 2.0 License, see [LICENSE](/LICENSE).
