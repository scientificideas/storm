**Storm** is a tool to create storms in your container runtime environment. You can specify a group of containers/pods that will stop and start in random order, or throw all available containers into the storm.

##### Installation
```
go install gitlab.n-t.io/atmz/storm@latest
```

<br />

##### Chaos level
```
storm -chaos hard 
```
The chaos level is simply how fast your containers will be affected by the tool (start, stop)

Available chaos levels: "easy", "medium" (default), "hard"

<br />

##### Filter targets
```
storm -filter containerGroup1,containerGroup2
```

If a filter is specified, the tool will work only with containers/pods that contain the specified pattern in their names.

<br />

##### Indicate the targets that the storm will hit
If you only want to storm only certain containers/pods, just list them:
```
storm -chaos hard -targets container1,container2,container3
```
No other containers will be affected.

<br />

##### Start containers immediately after stop
The containers will be restarted immediately after the fall
```
storm -startfast true
```

<br />

##### Choose container runtime/orchestrator which need a storm

```
storm -runtime=k8s -kube-namespace=some-namespace -kube-context=my-context
```

##### Kill k8s pods with pattern 'xyz' in pod's name

```
storm -filter=xyz -runtime=k8s -kube-namespace=some-namespace -kube-context=my-context
```