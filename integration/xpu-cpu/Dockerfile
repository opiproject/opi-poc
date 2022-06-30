FROM fedora:33

# TODO: consider replacing the base image with linuxserver/openssh-server:arm64v8-latest

RUN dnf install -y openssh-server

RUN mkdir /var/run/sshd
RUN useradd -rm -d /home/test -s /bin/bash -g root -G ssh_keys -u 1000 test 
RUN  echo 'test:test' | chpasswd
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd
RUN ssh-keygen -A

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
