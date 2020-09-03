FROM golang:1.14

#ENV GO111MODULE=on
#ENV GOCHACHE=off

#RUN go get -u github.com/go-sql-driver/mysql


#RUN groupadd --gid $GROUP_ID app && useradd -m -l --uid $USER_ID --gid $GROUP_ID $APP_USER
#RUN mkdir -p $APP_HOME && chown -R $APP_USER:$APP_USER $APP_HOME
#USER $APP_USER
#WORKDIR $APP_HOME

RUN groupadd -g 786 appuser && useradd -r -m -u 786 -g appuser appuser
RUN mkdir -p /go/src/webapp && chown -R appuser:appuser /go/src/webapp
USER appuser
WORKDIR /go/src/webapp
COPY . .
EXPOSE 8081

RUN go get -d -v ./...
RUN go install -v ./...


CMD ["webapp"]
