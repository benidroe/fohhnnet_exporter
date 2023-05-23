# Fohhn-Net exporter

Fohhn-Net Exporter for Prometheus.

This exporter uses the Fohhn-Net RS-485 Protocol.
Each Fohhn-Net network requires one RS-485 to TCP device server, running in 2 wire mode. Like the Moxa 5130.

Implementation of Fohhn-Net UDP and Fohhn Text Protocol will be done later.

Default telnet-listening port is 4001.

####Adapter Pinout:

```

RJ45                                   D-Sub female  

=========                              |====== 
        1|                             |       ======
         |                       ------|-+5          |
==       |                       |     |        9    | 
  |     4|+------ ( COLD - ) ----|-----|-+4          |
  |      |                       |     |        8    |
==      6|+-------( HOT +  ) ----|-----|-+3          |
         |                       |     |        7    |
Shield   |                       |     |  2          |
=========                        |     |        6    |
+------------------( GND   ) -----     |  1          |
                                       |       ====== 
                                       | =======
```



####Moxa Nport 5130 configuration
```
   ## serial settings Port 1
   BAUDRATE     19200
   DATA BITS    8
   STOP BITS    1
   PARITY       NONE
   FLOW CONTROL NONE
   FIFO         ENABLE
   INTERFACE    RS-485 2-Wire
   
   ## Operating Settings Port 1
   Operation Mode       TCP Server Mode
   TCP alive check time 7
   Inactivity time      0
   Max Connection       1        # important to avoid conflicts
   Ignore Jammed IP     No
   Allow driver control No
   Packet length        0
   Delimiter 1          0
   Delimiter 2          0
   Delimiter process    Do Nothing
   Force Transmit       0
   Local TCP Port       4001
   Command Port         966    
```


## Usage

```sh
./fohhnnet_exporter
```

Visit http://localhost:2121/fohhnnet?target=terminalserver.localnetwork where terminalserver.localnetwork is the IP or DNS-Name of the your Terminalserver to get metrics from.

## Installation

Clone this repository from github to your go directory. Within this repository run:

```
make build
```

Copy build to your /usr/local/bin

```
cp fohhnnet_exporter /usr/local/bin
```


Change permissions
```
chown prometheus:prometheus /usr/local/bin/fohhnnet_exporter
```

Setup exporter as daemon

```
Todo with systemctl script...
vim /etc/default/legaetx_exporter
START=yes

```


## Contributions
Thanks to https://github.com/prometheus/snmp_exporter. This project was used as example and for inspriration while realizing this exporter.

Thanks to Fohhn Company for their support.  
Command line utility written in C: https://github.com/Fohhn/fohhn-net-command-line-utility  
Command Reference: https://www.fohhn.com/en/technologies/integration-in-media-controls/

