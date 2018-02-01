From registry.ng.bluemix.net/sakshitest/terraform-ibm-provider-service:latest

ENV API_REPO /go/src/github.com/terraform-ibm-provider-api
COPY . $API_REPO
RUN cd $API_REPO && \
    go build -o apiserver
EXPOSE 9080
WORKDIR $API_REPO
CMD ["nohup ./apiserver &"]