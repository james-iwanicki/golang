FROM ubuntu:kinetic-20220830

ENV TZ=US/Denver
RUN ln -s /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt update && apt install openssh-server openssh-client vim bash iproute2 iputils-ping curl tcpdump -y
RUN echo "root:netscout"|chpasswd
RUN echo "shell-client1" > /root/whoami

RUN echo "PermitRootLogin=yes" >> /etc/ssh/sshd_config
RUN ssh-keygen -A
RUN service ssh start

CMD [ "/bin/bash", "-c", "service ssh start && tail -f /dev/null" ]

