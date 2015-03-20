FROM dockerfile/java:oracle-java8

RUN apt-get update

RUN apt-get install apt-transport-https
RUN wget -q http://support.quobyte.com/repo/1/<QUOBYTE_REPO_ID>/xUbuntu_14.04/Release.key -O - | sudo apt-key add -
RUN add-apt-repository "deb https://support.quobyte.com/repo/1/<QUOBYTE_REPO_ID>/xUbuntu_14.04 ./"
RUN apt-get update

RUN apt-get install -y quobyte-server quobyte-client dnsutils

ADD ./main.sh /opt/main.sh
ENTRYPOINT ["/opt/main.sh"]
