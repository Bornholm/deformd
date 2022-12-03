#!/sbin/openrc-run
  
command="/usr/bin/deformd"
command_args="run --config /etc/deformd/config.yml"

depend() {
    need net
}