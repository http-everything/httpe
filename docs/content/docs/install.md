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
User=httpe
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
Download and unpack the files to an appropriate location.

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

## Enable TLS aka HTTPS

HTTPE is capable to run with TLS enabled. Putting a reverse proxy in front of HTTP is not required.
To activate TLS you first need a certificate and the corresponding server key. You can generate this pair
using OpenSSL, or, if your server is exposed to the internet, you request a free Let's encrypt certificate.

### Retrieve a certificate

{{< alert context="info" text="The below documentation refers to Linux only." />}}

Install Certbot
```
apt install certbot
```

Make sure port 443 and 80 are exposed to the internet and not in use by another program.
Certbot will start its built-in webserver to validate the server address. Therefore, all
webservers listening on 80 or 443 must be stopped.

If your machine is behind NAT, create port forwarding for 80 and 443.

Once you got the certificate, you can remove the port forwarding, but this will stop the auto-renewal.

Request the certificate.

```shell
sudo certbot certonly --standalone \
  --agree-tos --register-unsafely-without-email -d <yourdomain.com>
```

### Activate TLS in HTTPE

Next change the group ownership so the httpe user can read the files.

```shell
chgrp -R httpe /etc/letsencrypt/archive/
chmod -R 0770 /etc/letsencrypt/archive/
chgrp -R httpe /etc/letsencrypt/live/
chmod -R 0770 /etc/letsencrypt/live/
chmod 0770 /var/log/httpe/
```

Now edit your `httpe.conf` file insert the paths to the certificate and key files.

```toml
[server]
## Defines the IP address and port the API server listens on.
## Environment variable HTTPE_SERVER_ADDRESS has precedence.
address = "0.0.0.0:443"

## A working directory
## Environment variable HTTPE_SERVER_DATA_DIR has precedence.
data_dir = "/var/lib/httpe"

## If both cert_file and key_file are specified, then rportd will use them to serve the API with TLS/https.
## Intermediate certificates should be included in cert_file if required.
## Environment variables HTTPE_SERVER_CERT_FILE and HTTPE_SERVER_KEY_FILE have precedence.
cert_file = "/etc/letsencrypt/live/<yourdomain.com>/fullchain.pem"
key_file = "/etc/letsencrypt/live/<yourdomain.com>/privkey.pem"

# ... snip snap, rest of your config
```

Start httpe and test it.
