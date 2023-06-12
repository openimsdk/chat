FROM ubuntu

# 设置固定的项目路径
ENV WORKDIR /chat
ENV CMDDIR $WORKDIR/script
ENV CONFIG_NAME $WORKDIR/config/config.yaml

# 将可执行文件复制到目标目录
ADD ./bin/open_im_admin $WORKDIR/bin/open_im_admin
ADD ./bin/open_im_admin_api $WORKDIR/bin/open_im_admin_api
ADD ./bin/open_im_chat $WORKDIR/bin/open_im_chat
ADD ./bin/open_im_chat_api $WORKDIR/bin/open_im_chat_api
ADD ./script $WORKDIR/script
ADD ./config/config.yaml $WORKDIR/config/config.yaml

# 创建用于挂载的几个目录，添加可执行权限
RUN mkdir $WORKDIR/logs && \
    chmod +x $WORKDIR/bin/open_im_admin $WORKDIR/bin/open_im_chat $WORKDIR/bin/open_im_admin_api $WORKDIR/bin/open_im_chat_api
RUN apt-get -qq update \
    && apt-get -qq install -y --no-install-recommends ca-certificates curl

VOLUME ["/chat/logs","/chat/config","/chat/script"]

WORKDIR $CMDDIR
CMD ["./docker_start_all.sh"]
