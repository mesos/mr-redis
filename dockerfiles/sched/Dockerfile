FROM ubuntu:14.04

EXPOSE 5454
EXPOSE 8080
EXPOSE 2379 

RUN mkdir -p /mrredis/bin
RUN mkdir -p /etcd/bin

COPY ./sched /mrredis/bin/
COPY ./config.json /mrredis/bin/
COPY ./MrRedisExecutor /mrredis/bin/
COPY ./redis-server /mrredis/bin/
COPY ./run_mrredis.sh /mrredis/bin/

COPY ./etcd /etcd/bin/
COPY ./runetcd.sh /etcd/bin/


CMD ["/mrredis/bin/run_mrredis.sh"]
