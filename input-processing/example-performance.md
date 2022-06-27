# This Program
```shell
andrewquerol@andrews-mbp input-processing % dd status=progress if=/dev/urandom bs=4096 count=2500000 | ./input-processing 
SP// Backend Developer Test - Input Processing
Standard error will contain the logging of this tool, standard out will only contain the filtered input

  10014150656 bytes (10 GB, 9550 MiB) transferred 39.001s, 257 MB/s 
2500000+0 records in
2500000+0 records out
10240000000 bytes transferred in 39.864104 secs (256872699 bytes/sec)
andrewquerol@andrews-mbp input-processing % 
```
# Grep with stdout to /dev/null
```shell
andrewquerol@andrews-mbp input-processing % dd status=progress if=/dev/urandom bs=4096 count=2500000 | grep "" > /dev/null
10226204672 bytes (10 GB, 9752 MiB) transferred 34.000s, 301 MB/s
2500000+0 records in
2500000+0 records out
10240000000 bytes transferred in 34.045677 secs (300772401 bytes/sec)
andrewquerol@andrews-mbp input-processing % 
```
