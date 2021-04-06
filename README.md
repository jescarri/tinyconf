How to use it
-------------

the tinyconf config file is composed of the following sections:

install_pkgs:
-------------

This is a yaml list of packages and versions you want to install, is mandatory to use a version, in order to make the process repeatable.

````
install_pkgs
- name: apache2
  version: 2.4.41-4ubuntu3.1
- name: foo
  version: bar+1
````

remove_pkgs:
------------

Yaml list of packages and versions to remove.

````
remove_pkgs:
- name: minicom
  version: 2.7.1-1.1
````

files:
------

List of files to write, you can specify GID, UID and Mode and the content in multiline way, if you need to restart a service

if the file changes, you can also specify the service name to restart, this svc name should match a systemd unit,

if you need to restart a service  you can also specify the service name to restart, this svc name should match a systemd unit

````
- name: /var/www/html/index.php
  gid: 33
  uid: 33
  mode: 0600
  content: |+
    <?php
      header("Content-Type: text/plain");
      echo "Hello, world!\n";
    ?>
  service: apache2
````

onboot_svcs:
------------

List of systemd units you want to enable at boot time, it will enable the unit and start it.

````
onboot_svcs:
- apache2
````

Run Example
-----------

````
tinyconf -h
tinyconf usage:.
      --config-file string   config file (default "tinyconf.yaml")
      root@workstation:/home/jescarri/workspace/go-workspace/src/github.com/jescarri/tinyconf# ./build/tinyconf --config-file=tinyconf.yaml
      INFO[0000] Pkg: apache2 is already installed and on required version: 2.4.41-4ubuntu3.1
      INFO[0000] Pkg: php is already installed and on required version: 2:7.4+75
      INFO[0001] Pkg: libapache2-mod-php is already installed and on required version: 2:7.4+75
      INFO[0001] Installing pkg: minicom version: 2.7.1-1.1
      INFO[0014] Removing pkg: minicom version :2.7.1-1.1
      INFO[0027] Enabling systemd unit: apache2
      INFO[0027] Enabling systemd unit: apache2
      INFO[0028] Making shure that unit: apache2 is started
      INFO[0028] starting systemd unit: apache2
````

real example
------------

````
root@ip-172-31-255-15:~# curl -LO https://github.com/jescarri/tinyconf/raw/master/output/tinyconf
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   140  100   140    0     0    230      0 --:--:-- --:--:-- --:--:--   230
100 3733k  100 3733k    0     0  4088k      0 --:--:-- --:--:-- --:--:-- 56.1M
root@ip-172-31-255-15:~# curl -LO https://raw.githubusercontent.com/jescarri/tinyconf/master/prod.yaml
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   588  100   588    0     0  18967      0 --:--:-- --:--:-- --:--:-- 18967
root@ip-172-31-255-15:~# chmod +x ./tinyconf
root@ip-172-31-255-15:~# ./tinyconf -h
tinyconf usage:.
      --config-file string   config file (default "tinyconf.yaml")
root@ip-172-31-255-15:~# ./tinyconf --config-file=prod.yaml
INFO[0000] Installing pkg: apache2 version: 2.4.29-1ubuntu4.14
INFO[0011] Installing pkg: php version: 1:7.2+60ubuntu1
INFO[0025] Installing pkg: libapache2-mod-php version: 1:7.2+60ubuntu1
INFO[0028] File: /etc/apache2/mods-available/dir.conf requires change
INFO[0028] Service: apache2 will require restart
INFO[0028] Restarting systemd unit: apache2
INFO[0028] Enabling systemd unit: apache2
INFO[0028] Enabling systemd unit: apache2
INFO[0029] Making shure that unit: apache2 is started
INFO[0029] starting systemd unit: apache2
````

verification:
------------

````
$ for i in $(echo "18.207.149.79 52.90.233.103"); do curl -sv http://$i ; done
*   Trying 18.207.149.79:80...
* TCP_NODELAY set
* Connected to 18.207.149.79 (18.207.149.79) port 80 (#0)
> GET / HTTP/1.1
> Host: 18.207.149.79
> User-Agent: curl/7.68.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Tue, 06 Apr 2021 05:07:01 GMT
< Server: Apache/2.4.29 (Ubuntu)
< Content-Length: 14
< Content-Type: text/plain;charset=UTF-8
<
Hello, world!
* Connection #0 to host 18.207.149.79 left intact
*   Trying 52.90.233.103:80...
* TCP_NODELAY set
* Connected to 52.90.233.103 (52.90.233.103) port 80 (#0)
> GET / HTTP/1.1
> Host: 52.90.233.103
> User-Agent: curl/7.68.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Tue, 06 Apr 2021 05:07:01 GMT
< Server: Apache/2.4.29 (Ubuntu)
< Content-Length: 14
< Content-Type: text/plain;charset=UTF-8
<
Hello, world!
* Connection #0 to host 52.90.233.103 left intact
````
