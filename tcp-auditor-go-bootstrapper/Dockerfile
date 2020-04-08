FROM mcr.microsoft.com/dotnet/core/runtime:2.2-alpine

RUN apk add --no-cache curl

ENV CONSULCONTAINER=consul:8500
ENV PUBLISHENVSETTINGS=0
ENV MANIFESTCONTAINER=tcp-manifest-service
ENV PUBLISHAPPSETTINGS=1
ENV DEPLOYMENTSTAGE=CI


WORKDIR /
COPY /init/bootstrap.sh /init/
COPY /config/stage/CI/ /config/ci/appsettings/
COPY /config/stage/QA/ /config/qa/appsettings/
COPY /config/stage/PROD/ /config/prod/appsettings/
COPY entrypoint.sh .
ENTRYPOINT [ "sh", "-c","/entrypoint.sh" ]
