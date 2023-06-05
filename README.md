# Prometheus Fohhn-Net exporter

This Fohhn-Net exporter is a possibility to provide data in a format, that Prometheus can ingest.  

Important: The implementation is independent and does not belong to the manufacturer.

### Requirements
Currently only Fohhn-Net RS-485 protocol is supported.
Each RS-485 based Fohhn-Net need an external ethernet adapter. The following adapters may work:

* Fohhn NA-4 via UDP (not tested yet)
* Moxa 5130 RS-485 to Ethernet device server via TCP

Some devices with integrated network interface can communicate directly via UDP. However, a few devices only use their interface for audio and can not be used for queries.

Currently only the first six IDs (1-6) are scraped. The result will contain all values from responding devices.   
Fohhn Text Protocol is not implemented yet.

#### Default Ports
```
UDP  |   2101
TCP  |   4001
```

## Moxa 5130 Setup

#### Adapter Pinout for Moxa 5130 device server on RJ-45

```

Fohhn-Net                              D-Sub female  

 ========                              |====== 
       1-|                             |       ======
       2-|                       ------|-+5          |
       3-|                       |     |        9    | 
       4-|+------- ( COLD - ) ---|-----|-+4          |
       5-|                       |     |        8    |
       6-|+------- ( HOT +  ) ---|-----|-+3          |
       7-|                       |     |        7    |
       8-|                       |     |  2          |
=Shield==                        |     |        6    |
       +_________________________|     |  1          |
                                       |       ====== 
                                       | ======
```

#### Adapter Pinout for Moxa 5130 device server on Screw Connector

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

### fohhnnet_exporter

```
Flags:                                                                       
  -h, --help           Show context-sensitive help (also try --help-long
                            and --help-man).                                 
      --web.listen-address=":2121"                                           
                            Address to listen on for web interface and       
                            telemetry.                                       
      --log.level="Debug"   LogLevel - Debug, Info, Warn, Error              
      --fnet.port.udp=2101  UDP Port for target devices                      
      --fnet.port.tcp=4001  TCP Port for target devices     
```

#### Examples
```sh
./fohhnnet_exporter
./fohhnnet_exporter --log.level Info
./fohhnnet_exporter --fnet.port.udp 4021 --log.level Info
./fohhnnet_exporter --fnet.port.udp 4021 --fnet.port.tcp 4005 --log.level Info
```
Visit the following urls,where terminalserver.localnetwork is the IP or DNS-Name of the your Terminalserver to get metrics from.

* http://localhost:2121/fohhnnetudp?target=terminalserver.localnetwork (for UDP-Requests)
* http://localhost:2121/fohhnnettcp?target=terminalserver.localnetwork (for TCP-Requests)



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

