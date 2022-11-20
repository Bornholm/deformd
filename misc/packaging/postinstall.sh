#!/bin/sh

use_systemctl="True"
systemd_version=0
if ! command -V systemctl >/dev/null 2>&1; then
  use_systemctl="False"
else
    systemd_version=$(systemctl --version | head -1 | sed 's/systemd //g')
fi

service_name=deformd

cleanup() {
    if [ "${use_systemctl}" = "False" ]; then
        rm -f /usr/lib/systemd/system/${service_name}.service
    else
        rm -f /etc/chkconfig/${service_name}
        rm -f /etc/init.d/${service_name}
    fi
}

cleanInstall() {
    printf "\033[32m Post Install of an clean install\033[0m\n"
    if [ "${use_systemctl}" = "False" ]; then
        if command -V chkconfig >/dev/null 2>&1; then
          chkconfig --add ${service_name}
        fi

        service ${service_name} restart ||:
    else
        # rhel/centos7 cannot use ExecStartPre=+ to specify the pre start should be run as root
        # even if you want your service to run as non root.
        if [ "${systemd_version}" -lt 231 ]; then
            printf "\033[31m systemd version %s is less then 231, fixing the service file \033[0m\n" "${systemd_version}"
            sed -i "s/=+/=/g" /usr/lib/systemd/system/${service_name}.service
        fi
        printf "\033[32m Reload the service unit from disk\033[0m\n"
        systemctl daemon-reload ||:
        printf "\033[32m Unmask the service\033[0m\n"
        systemctl unmask ${service_name} ||:
        printf "\033[32m Set the preset flag for the service unit\033[0m\n"
        systemctl preset ${service_name} ||:
        printf "\033[32m Set the enabled flag for the service unit\033[0m\n"
        systemctl enable ${service_name} ||:
        systemctl restart ${service_name} ||:
    fi
}

upgrade() {
    printf "\033[32m Post Install of an upgrade\033[0m\n"
}

# Step 2, check if this is a clean install or an upgrade
action="$1"
if  [ "$1" = "configure" ] && [ -z "$2" ]; then
  action="install"
elif [ "$1" = "configure" ] && [ -n "$2" ]; then
    action="upgrade"
fi

case "$action" in
  "1" | "install")
    cleanInstall
    ;;
  "2" | "upgrade")
    printf "\033[32m Post Install of an upgrade\033[0m\n"
    upgrade
    ;;
  *)
    printf "\033[32m Alpine\033[0m"
    cleanInstall
    ;;
esac

cleanup