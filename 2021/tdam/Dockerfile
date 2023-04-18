FROM ubuntu:22.10

ENV TZ=US/Denver
RUN ln -s /usr/share/info/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apt-get update && apt install openssh-server openssh-client vim curl bash iproute2 iputils-ping dnsutils net-tools tcpdump libcap2-bin wireshark -y

RUN echo "root:netscout"|chpasswd
RUN echo "server" > /root/whoami

RUN echo "PermitRootLogin=yes" > /etc/ssh/sshd_config
RUN ssh-keygen -A
RUN service ssh start

ENV version=v1
WORKDIR /build
COPY . .
RUN cp /build/tdam-client /root/.

CMD [ "/bin/bash", "-c", "service ssh start && /root/tdam-client" ]
