FROM ubuntu

# Set fixed project paths
ENV WORKDIR /chat
ENV CMDDIR $WORKDIR/script
ENV CONFIG_NAME $WORKDIR/config/config.yaml

# Copy executable files to the target directory
COPY ./bin/open_im_admin $WORKDIR/bin/open_im_admin
COPY ./bin/open_im_admin_api $WORKDIR/bin/open_im_admin_api
COPY ./bin/open_im_chat $WORKDIR/bin/open_im_chat
COPY ./bin/open_im_chat_api $WORKDIR/bin/open_im_chat_api
COPY ./script $WORKDIR/script
COPY ./config/config.yaml $WORKDIR/config/config.yaml

# Create directories for mounting and add executable permissions
RUN mkdir $WORKDIR/logs && \
    chmod +x $WORKDIR/bin/open_im_admin $WORKDIR/bin/open_im_chat $WORKDIR/bin/open_im_admin_api $WORKDIR/bin/open_im_chat_api

RUN apt-get -qq update \
    && apt-get -qq install -y --no-install-recommends ca-certificates curl

VOLUME ["/chat/logs","/chat/config","/chat/script"]

WORKDIR $CMDDIR
CMD ["./docker_start_all.sh"]