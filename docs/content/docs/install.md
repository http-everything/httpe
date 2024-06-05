---
weight: 200
title: "Install & configure"
description: ""
icon: "article"
date: "2024-02-12T12:58:48+01:00"
lastmod: "2024-02-12T12:58:48+01:00"
draft: false
toc: true
slug: install
---

## Install the software

### On Linux

Download the most recent package for Linux, unpack and move the included files to the recommended locations.

```shell
sudo mkdir /etc/httpe
sudo mkdir /var/log/httpe
cd /tmp/
DOWNLOAD="httpe_{{% httpe-version %}}_Linux_$(uname -m).tar.gz"
curl -LO https://github.com/http-everything/httpe/releases/download/{{% httpe-version %}}/${DOWNLOAD}
tar xzf $DOWNLOAD
sudo mv httpe /usr/local/bin/httpe
chmod +x /usr/local/bin/httpe
sudo mv example.httpe.conf /etc/httpe/httpe.conf
sudo mv example.rules.unix.yaml /etc/httpe/rules.yaml
rm ${DOWNLOAD} example.rules.win.yaml
```

It's highly recommended not to run the httpe server as the root user. Create a dedicated user and change the ownership
of the previously created files.

```shell
sudo useradd -d /var/lib/httpe -m -U -r -s /bin/false httpe
sudo chgrp httpe /var/log/httpe/
sudo chgrp httpe /etc/httpe/
```

Finally, create a systemd service, to run the server as a background server and start it automatically on boot.

Create a file `/etc/systemd/system/httpe.service` with the following content:

```text
#
# HTTPE systemd service 
#
[Unit]
Description=HTTPE low code application server
ConditionFileIsExecutable=/usr/local/bin/httpe
StartLimitIntervalSec=5
StartLimitBurst=10
Documentation=https://httpe.io/docs/
After=network.target network-online.target

[Service]
ExecStart=/usr/local/bin/httpe "-c" "/etc/httpe/httpe.conf"
User=rport
Restart=always
RestartSec=120
EnvironmentFile=-/etc/sysconfig/httpe
AmbientCapabilities=CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target

```

Then start the server. Starting will very likely fail, because the configuration file is not yet ready. Continue reading
about the configuration below. 

### On Windows

Open a PowerShell terminal **with administrative rights**.
Download and unpack the files to an appropriated location.

```powershell
mkdir "C:\Program Files\httpe"
cd "C:\Program Files\httpe"
$download = "httpe_{{% httpe-version %}}_Windows_x86_64.zip"
iwr "https://github.com/http-everything/httpe/releases/download/{{% httpe-version%}}/$download" -OutFile $download
Expand-Archive -DestinationPath . $download
rm .\$download, .\example.rules.unix.yaml
mv example.httpe.conf httpe.conf
mv example.rules.win.yaml rules.yaml
```

The current version of HTTPE cannot run as a Windows service. This might change in the future. Meanwhile, you can use
[NSSM](https://nssm.cc/) to run the server as a Windows service.

```powershell
iwr "https://nssm.cc/release/nssm-2.24.zip" -OutFile nssm-2.24.zip
mv .\nssm-2.24\win64\nssm.exe .
rm .\nssm-2.24* -Force -Recurse
.\nssm.exe install httpe "C:\Program Files\httpe\httpe.exe" "-c `"C:\Program Files\httpe\httpe.conf`""
```

### On macOS

TBD

## Edit the configuration file

Open the configuration file, `/etc/httpe/httpe.conf` on Linux, `C:\Program Files\httpe\httpe.conf` on Windows and 
change it to your needs.

Remember to quote backslashes with backslashes to enter path values in the config file. For example:

```text
## Specifies the rules file
## Environment variable HTTPE_SERVER_RULES_FILE has precedence.
rules_file = "C:\\Program Files\\httpe\\rules.yaml"
```