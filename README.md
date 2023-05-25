# Prometheus Fohhn-Net exporter

This Fohhn-Net exporter is a possibility to provide data in a format, that Prometheus can ingest.

This exporter uses the Fohhn-Net RS-485 protocol.
Devices with Fohhn-Net RS-485 require an external ethernet adapter. The following adapters may work:

* Fohhn NA-4 via UDP (not tested yet)
* Moxa 5130 RS-485 to Ethernet device server via TCP

Some devices with integrated network interface can communicate directly via UDP. However, other devices only use their interface for audio.

Fohhn Text Protocol is not implemented yet.

#### Default Ports
```
UDP  |   2101
TCP  |   4001
```

## Moxa 5130 Setup

#### Adapter Pinout for Moxa 5130 device server

```

Fohhn-Net                              D-Sub female  

 ========                              |====== 
 |       |                             |       ======
-  - o   |-------                ------|-+5          |
 |       |      |                |     |        9    | 
-  G o   |----  -- ( COLD - ) ---|-----|-+4          |
 |       |   |                   |     |        8    |
-  + o   |+--|---- ( HOT +  ) ---|-----|-+3          |
 |       |   |                   |     |        7    |
 =========   ---------------------     |  2          |
                                       |        6    |
                                       |  1          |
                                       |       ====== 
                                       | =======
```



#### Moxa Nport 5130 configuration
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

### fohhn-cli

There is a simple fohhn-cli. It can be used to check, if a device is responding correctly. 

#### Usage
```
Flags:                                                                       
      --help           Show context-sensitive help (also try --help-long and
                        --help-man).                                                                   
      --id              Query device with id                              
      --scan            Scan for devices                                     
      --all             Show data from all scanned devices                        
      --port            Port of target device                                
  -p, --protocol        Use tcp or udp                                       
                                                                             
Args:                                                                        
  <host>  Host or IP of target device                                        

fohhn-cli --scan --all

``` 

#### Examples
```
fohhn-cli --scan --all 10.0.0.2
fohhn-cli --id 5 --port 4028 --protocol tcp 10.0.0.2
```

## Installation

Clone this repository from github to your go directory. Within this repository run:

```
make configure
make build
sudo make install
```

Make install copies the two binares to your /usr/local/bin.

Setup exporter as systemd service

```
sudo install -o root -g root fohhnnet_exporter.service /usr/lib/systemd/system

sudo systemctl enable fohhnnet_exporter.service
sudo systemctl start fohhnnet_exporter.service

```


## Contributions
Thanks to https://github.com/prometheus/snmp_exporter. This project was used as example and for inspriration while realizing this exporter.

Thanks to Fohhn Company for their support.  
Command line utility written in C: https://github.com/Fohhn/fohhn-net-command-line-utility  
Command Reference: https://www.fohhn.com/en/technologies/integration-in-media-controls/

