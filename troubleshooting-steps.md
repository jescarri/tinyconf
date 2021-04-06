Connectivity tests from the internet
------------------------------------

````
~/workspace/go-workspace/src/github.com/jescarri/tinyconf [master|✔]
05:10 $ curl http://3.88.100.61
^C
✘-INT ~/workspace/go-workspace/src/github.com/jescarri/tinyconf [master|✔]
05:11 $ sudo tcptraceroute 3.88.100.61 80
Selected device eth0, address 192.168.2.103, port 45541 for outgoing packets
Tracing the path to 3.88.100.61 on TCP port 80 (http), 30 hops max
 1  192.168.2.1  0.503 ms  0.401 ms  0.333 ms
 2  97.107.189.129  0.806 ms  0.668 ms  0.733 ms
 3  xe-4-0-3.cr0-van1.ip4.gtt.net (173.205.56.221)  0.767 ms  0.682 ms  0.660 ms
 4  ae19.cr9-chi1.ip4.gtt.net (141.136.108.189)  46.845 ms  47.051 ms  47.255 ms
 5  ip4.gtt.net (76.74.1.10)  47.016 ms  49.607 ms  53.774 ms
 6  150.222.76.133  48.729 ms  48.031 ms  47.960 ms
 7  52.95.63.147  47.088 ms  46.992 ms  47.564 ms
 8  * * *
 9  * * *
10  *^C
````

seems like something is blocking the traffic.

tests on the host
-----------------

check for iptables rules

`````
root@ip-172-31-255-65:~# iptables -L --line-numbers
Chain INPUT (policy ACCEPT)
num  target     prot opt source               destination
1    DROP       tcp  --  anywhere             anywhere             tcp dpt:http

Chain FORWARD (policy ACCEPT)
num  target     prot opt source               destination

Chain OUTPUT (policy ACCEPT)
num  target     prot opt source               destination
`````

confirmed, an iptables rule

Delete rule

````
root@ip-172-31-255-65:~# iptables -D INPUT 1
root@ip-172-31-255-65:~# iptables -L --line-numbers
Chain INPUT (policy ACCEPT)
num  target     prot opt source               destination

Chain FORWARD (policy ACCEPT)
num  target     prot opt source               destination

Chain OUTPUT (policy ACCEPT)
num  target     prot opt source               destination
````

test again from outside
-----------------------

````
05:18 $ curl -v http://3.88.100.61
*   Trying 3.88.100.61:80...
* TCP_NODELAY set
* Connected to 3.88.100.61 (3.88.100.61) port 80 (#0)
> GET / HTTP/1.1
> Host: 3.88.100.61
> User-Agent: curl/7.68.0
> Accept: */*
>
^C
````

I can connect but something is blocking the http response.

check index page to see what's going on.

look for apache or any process that looks like is listenning on port 80.

looks that apache is not installed but theres a netcat process listenning on port 80.

````
root@ip-172-31-255-65:~# ps -efa | grep nc
systemd+   488     1  0 Mar31 ?        00:00:30 /lib/systemd/systemd-timesyncd
root      1207     1  0 Mar31 ?        00:00:00 nc -k -l 0.0.0.0 80
root      5600  4716  0 05:21 pts/0    00:00:00 grep --color=auto nc
````

vitals check
------------

free space.

````
root@ip-172-31-255-65:~# df -h
Filesystem      Size  Used Avail Use% Mounted on
udev            224M     0  224M   0% /dev
tmpfs            48M  5.5M   43M  12% /run
/dev/xvda1      7.7G  7.7G     0 100% /
tmpfs           238M     0  238M   0% /dev/shm
tmpfs           5.0M     0  5.0M   0% /run/lock
tmpfs           238M     0  238M   0% /sys/fs/cgroup
/dev/loop0       98M   98M     0 100% /snap/core/9993
/dev/loop1       29M   29M     0 100% /snap/amazon-ssm-agent/2012
tmpfs            48M     0   48M   0% /run/user/0
````

there's something that is filling the disk.

Let's look at the process list.

```
root      1207     1  0 Mar31 ?        00:00:00 nc -k -l 0.0.0.0 80
root      1212     1  0 Mar31 ?        00:00:00 /bin/sh /sbin/named
```

this is a webserver why named is running?

let's see what kind of file is named.

````
file /sbin/named
/sbin/named: POSIX shell script, ASCII text executable
root@ip-172-31-255-65:/# cat /sbin/named-bash: cannot create temp file for here-document: No space left on device

#!/bin/sh
set -e
TMP="/tmp/tmp.zzBKjVA0Gg"
exec 3>"$TMP"
dd bs="104857600" count="200" if="/dev/zero" of="$TMP" || :
rm -f "$TMP"
kill -STOP "$$"
````

this process is filling the disk!.

````
root@ip-172-31-255-65:/# kill -9 1212
root@ip-172-31-255-65:/# df -h
Filesystem      Size  Used Avail Use% Mounted on
udev            224M     0  224M   0% /dev
tmpfs            48M  1.4M   47M   3% /run
/dev/xvda1      7.7G  1.2G  6.6G  15% /
tmpfs           238M     0  238M   0% /dev/shm
tmpfs           5.0M     0  5.0M   0% /run/lock
tmpfs           238M     0  238M   0% /sys/fs/cgroup
/dev/loop0       98M   98M     0 100% /snap/core/9993
/dev/loop1       29M   29M     0 100% /snap/amazon-ssm-agent/2012
tmpfs            48M     0   48M   0% /run/user/0
````

let's also kill netcat

```
root@ip-172-31-255-65:/# kill -9 1207
root@ip-172-31-255-65:/# netstat -an|grep LIST
tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN
tcp6       0      0 :::22                   :::*                    LISTEN
```

nothing listenning.

apparently can't curl, and it's because resolv.conf is not present

configure resolved
------------------

vi /etc/systemd/resolved.conf

````
[Resolve]
DNS=8.8.8.8
FallbackDNS=1.1.1.1
````

restart resolved unit
---------------------

````
systemctl restart systemd-resolved.service^C
root@ip-172-31-255-65:/etc# cat /etc/resolv.conf
# This file is managed by man:systemd-resolved(8). Do not edit.
#
# This is a dynamic resolv.conf file for connecting local clients to the
# internal DNS stub resolver of systemd-resolved. This file lists all
# configured search domains.
#
# Run "systemd-resolve --status" to see details about the uplink DNS servers
# currently in use.
#
# Third party programs must not access this file directly, but only through the
# symlink at /etc/resolv.conf. To manage man:resolv.conf(5) in a different way,
# replace this symlink by a static file or a different symlink.
#
# See man:systemd-resolved.service(8) for details about the supported modes of
# operation for /etc/resolv.conf.

nameserver 127.0.0.53
options edns0
search ec2.internal
root@ip-172-31-255-65:/etc# dig google.com

; <<>> DiG 9.11.3-1ubuntu1.13-Ubuntu <<>> google.com
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 53039
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 65494
;; QUESTION SECTION:
;google.com.			IN	A

;; ANSWER SECTION:
google.com.		252	IN	A	172.217.15.78

;; Query time: 1 msec
;; SERVER: 127.0.0.53#53(127.0.0.53)
;; WHEN: Tue Apr 06 05:50:49 UTC 2021
;; MSG SIZE  rcvd: 55
````

Install apache using tinyconf
-----------------------------

````
root@ip-172-31-255-65:/# curl -LO https://github.com/jescarri/tinyconf/raw/master/output/tinyconf
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   140  100   140    0     0    243      0 --:--:-- --:--:-- --:--:--   243
100 3733k  100 3733k    0     0  5148k      0 --:--:-- --:--:-- --:--:-- 5148k
root@ip-172-31-255-65:/# chmod +x tinyconf
root@ip-172-31-255-65:/# curl -LO https://raw.githubusercontent.com/jescarri/tinyconf/master/prod.yaml
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   588  100   588    0     0   7170      0 --:--:-- --:--:-- --:--:--  7170
root@ip-172-31-255-65:/# chmod +x tinyconf

apt-get update
Get:1 http://security.ubuntu.com/ubuntu bionic-security InRelease [88.7 kB]
Get:2 http://archive.ubuntu.com/ubuntu bionic InRelease [242 kB]
Get:3 http://security.ubuntu.com/ubuntu bionic-security/main amd64 Packages [1639 kB]
...
....
./tinyconf --config-file=prod.yaml
INFO[0000] Installing pkg: apache2 version: 2.4.29-1ubuntu4.14
INFO[0047] Installing pkg: php version: 1:7.2+60ubuntu1
INFO[0068] Installing pkg: libapache2-mod-php version: 1:7.2+60ubuntu1
INFO[0071] File: /etc/apache2/mods-available/dir.conf requires change
INFO[0071] Service: apache2 will require restart
INFO[0071] Restarting systemd unit: apache2
INFO[0071] Enabling systemd unit: apache2
INFO[0071] Enabling systemd unit: apache2
INFO[0072] Making shure that unit: apache2 is started
INFO[0072] starting systemd unit: apache2
````

test local
----------

````
root@ip-172-31-255-65:/# curl -sv http://127.0.0.1
* Rebuilt URL to: http://127.0.0.1/
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 80 (#0)
> GET / HTTP/1.1
> Host: 127.0.0.1
> User-Agent: curl/7.58.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Tue, 06 Apr 2021 05:56:11 GMT
< Server: Apache/2.4.29 (Ubuntu)
< Content-Length: 14
< Content-Type: text/plain;charset=UTF-8
<
Hello, world!
* Connection #0 to host 127.0.0.1 left intact

````

test from the internet
----------------------

````
 $ curl -v http://3.88.100.61
 *   Trying 3.88.100.61:80...
 * TCP_NODELAY set
 * Connected to 3.88.100.61 (3.88.100.61) port 80 (#0)
 > GET / HTTP/1.1
 > Host: 3.88.100.61
 > User-Agent: curl/7.68.0
 > Accept: */*
 >
 * Mark bundle as not supporting multiuse
 < HTTP/1.1 200 OK
 < Date: Tue, 06 Apr 2021 05:56:42 GMT
 < Server: Apache/2.4.29 (Ubuntu)
 < Content-Length: 14
 < Content-Type: text/plain;charset=UTF-8
 <
 Hello, world!
 * Connection #0 to host 3.88.100.61 left intact
````
