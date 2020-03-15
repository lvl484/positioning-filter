FROM consul:1.7.1

COPY agent-config.json /consul/configs/

ENTRYPOINT ["consul", "agent", "-config-dir=/consul/configs", "-data-dir=/tmp/consul"]