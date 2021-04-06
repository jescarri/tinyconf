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
