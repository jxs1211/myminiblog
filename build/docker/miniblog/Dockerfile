# Copyright 2020 Jayflow <jxs1211@gmail.com>. All rights reserved.    
# Use of this source code is governed by a MIT style    
# license that can be found in the LICENSE file.

FROM BASE_IMAGE
LABEL maintainer="<jxs1211@gmail.com>"

WORKDIR /opt/miniblog

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone && \
      mkdir -p /var/log/miniblog

COPY miniblog /opt/miniblog/bin/

ENTRYPOINT ["/opt/miniblog/bin/miniblog"]
CMD ["-c", "/etc/miniblog/miniblog.yaml"]
