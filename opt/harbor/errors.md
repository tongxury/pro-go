
### 安装时遇到过的问题

1. 部分pod 需要设置 
```shell
securityContext:
  ...
  allowPrivilegeEscalation: true
  privilege: true
```

2. 1:C 28 Nov 2024 03:50:55.610 # WARNING Memory overcommit must be enabled! Without it, a background save or replication may fail under low memory condition. Being disabled, it can also cause failures without low memory condition, see https://github.com/jemalloc/jemalloc/issues/1328. To fix this issue add 'vm.overcommit_memory = 1' to /etc/sysctl.conf and then reboot or run the command 'sysctl vm.overcommit_memory=1' for this to take effect.
```shell
vm.overcommit_memory = 1
```