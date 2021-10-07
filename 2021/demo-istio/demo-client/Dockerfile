FROM ubuntu:latest

ENV TZ=US/Denver
RUN ln -s /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt update && apt install openssh-server openssh-client iproute2 iputils-ping curl tcpdump bash vim -y
RUN echo "root:netscout"|chpasswd
RUN echo "demo-client" > /root/whoami

RUN echo "PermitRootLogin=yes" >> /etc/ssh/sshd_config
RUN ssh-keygen -A
RUN service ssh start

CMD [ "/bin/bash", "-c", "service ssh start && tail -f /dev/null" ]

