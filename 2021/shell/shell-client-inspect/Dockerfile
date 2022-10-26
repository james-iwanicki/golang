FROM ubuntu:22.04

ENV TZ=US/Denver
RUN ln -s /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt update && apt install openssh-server openssh-client vim bash iproute2 iputils-ping curl tcpdump libcap2-bin wireshark -y
RUN echo "root:netscout"|chpasswd
RUN echo "shell-client-inspect" > /root/whoami
COPY ./cni-inspect /root/.

RUN echo "PermitRootLogin=yes" >> /etc/ssh/sshd_config
RUN ssh-keygen -A
RUN service ssh start

CMD [ "/bin/bash", "-c", "service ssh start && tail -f /dev/null" ]

