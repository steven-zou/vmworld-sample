FROM library/photon:1.0

MAINTAINER szou@vmware.com

RUN mkdir /vmworld2017

COPY ./vmworld-sample /vmworld2017/

RUN chmod u+x /vmworld2017/vmworld-sample

WORKDIR /vmworld2017/

ENTRYPOINT ["/vmworld2017/vmworld-sample"]