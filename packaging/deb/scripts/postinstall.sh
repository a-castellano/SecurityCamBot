#!/bin/sh

mkdir -p /etc/windmaker-security-cam-bot

echo "### NOT starting on installation, please execute the following statements to configure windmaker-security-cam-bot to start automatically using systemd"
echo " sudo /bin/systemctl daemon-reload"
echo " sudo /bin/systemctl enable windmaker-security-cam-bot"
echo "### You can start grafana-server by executing"
echo " sudo /bin/systemctl start windmaker-security-cam-bot"
