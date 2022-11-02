FROM ubuntu:kinetic-20220830
ENV TZ=US/Denver
RUN ln -s /usr/share/info/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt update && apt install openssh-server openssh-client bash vim iproute2 iputils-ping dnsutils net-tools -y
RUN echo "root:netscout" | chgpasswd
RUN echo "shell-mirror" > /root/whoami

RUN echo "PermitRootLogin=yes" > /etc/ssh/sshd_config
RUN ssh-keygen -A
RUN service ssh start

CMD [ "/bin/bash", "-c", "service ssh start && tail -f /dev/null" ]
