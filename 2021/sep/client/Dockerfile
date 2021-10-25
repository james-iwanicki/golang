FROM ubuntu:latest

ENV TZ=US/Denver
RUN ln -s /usr/share/info/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apt update && apt install openssh-server openssh-client vim curl bash iproute2 iputils-ping dnsutils net-tools tcpdump libcap2-bin wireshark -y

RUN echo "root:netscout"|chpasswd
RUN echo "server" > /root/whoami

RUN echo "PermitRootLogin=yes" > /etc/ssh/sshd_config
RUN ssh-keygen -A
RUN service ssh start

ENV version=v1
WORKDIR /build
COPY . .
RUN cp /build/demo-client /root/.

CMD [ "/bin/bash", "-c", "service ssh start && /root/demo-client" ]
#CMD [ "/root/demo-client" ]
