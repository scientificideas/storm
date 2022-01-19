**Storm** is a tool to create storms in your docker environment. You can specify a group of containers that will stop and start in random order, or throw all available containers into the storm.

##### Installation
```
go install gitlab.n-t.io/atmz/storm@latest
```

##### Chaos level
```
storm -chaos hard 
```
The chaos level is simply how fast your containers will be affected by the tool (start, stop)

Available chaos levels: "easy", "medium" (default), "hard"

##### Filter targets
```
storm -filter containerGroup1,containerGroup2
```

If a filter is specified, the tool will work only with containers that contain the specified pattern in their names.


##### Indicate the targets that the storm will hit
If you only want to storm only certain containers, just list them:
```
storm -chaos hard -targets container1,container2,container3
```
No other containers will be affected.

##### Start containers immediately after stop
The containers will be restarted immediately after the fall, just like under the k8s orchestration
```
storm -startfast true
```