### Cron HPA
Cron HPA (Horizontal Pod Autoscaler) is a lightweight deployment that schedules changes at specific times (UTC) to a horizontal pod autoscaler.

<b>IMPORTANT:</b> this is solely an example that does not carry official support and should be used as a starting point, not a production ready deployment.

#### Config
```
version: 0.0.1            ==> printed when this application starts
namespace: default        ==> kubernetes namespace
hpaName: ssg              ==> the name of the HPA you'd like to target (kubectl get hpa)
scaleUpMaxReplicas: 3     ==> MaxReplicas when deployment is scaled up
scaleUpMinReplicas: 2     ==> MinReplicas when deployment is scaled up
scaleDownMaxReplicas: 2   ==> MaxReplicas when deployment is scaled down
scaleDownMinReplicas: 1   ==> MinReplicas when deployment is scaled down
schedule:
  enabled: true           ==> Enable/Disable the schedule 
  scaleUpTime: 16:00      ==> scale up time (UTC)
  scaleDownTime: 23:00    ==> scale down time (UTC)
```

#### Install or Upgrade
update config.yaml with your desired configuration
```
kubectl apply -k cron-hpa
```

#### Build Docker Image
```
docker build -t <host>:<port>/<project/username>/cron-hpa:<tag> .
```
Example
```
docker build -t docker.io/layer7api/cron-hpa:0.0.1 .
```

#### Push Docker Image
```
docker push <host>:<port>/<project/username>/cron-hpa:<tag>
```