# RXTXSpeed Input Plugin

This plugin gathers metrics about rx/tx speed.  The measurements based on
the number of received and sent bytes from `/sys/class/net/*/statistics/`

### Configuration:

You can generate it using `telegraf --usage <plugin-name>`.

```toml
# Set the network interface
[[inputs.rxtxspeed]]
  Network_interface = "enp3s0"
```

### Metrics:

- rx/tx speed
  - fields:
    - rx speed (float, MiB)
    - tx speed (float, MiB)

### Example Output:

```
rxtxspeed,host=test-VirtualBox rx\ speed=0 1453831884664956455
```
