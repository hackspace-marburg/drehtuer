# drehtuer

Connects the [door's status](https://hsmr.cc/Infrastruktur/Door), received via MQTT, with either an InfluxDB or Prometheus.
The result can be viewed in the [Grafana](https://grafana.hsmr.cc/dashboard/db/door-state), for example.


## Installation

```
go build

vim systemd/drehtuer.conf

sudo cp drehtuer /usr/bin
sudo cp systemd/drehtuer.conf /etc/conf.d
sudo cp systemd/drehtuer.service /etc/systemd/system

sudo systemctl daemon-reload
sudo systemctl enable drehtuer.service
sudo systemctl start drehtuer.service
```
