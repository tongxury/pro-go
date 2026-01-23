#!/bin/bash

# ==============================================================================
# 脚本名称: cleanup_traces.sh
# 功能描述: 清理 Linux 系统中的登录记录、操作历史和系统日志
# 使用说明: 请使用 bash 执行：bash cleanup_traces.sh 或 ./cleanup_traces.sh
# ==============================================================================

# 确保脚本即使在报错时也继续执行
set +e

echo "Starting cleanup process..."

# 1. 清理当前 Shell 内存中的历史记录
# 注意：history 命令是 shell 内核内置命令，在非交互式脚本中可能无效
if command -v history >/dev/null 2>&1; then
    history -c
    history -w
fi

# 2. 清理用户个人历史文件
# 使用路径遍历，确保覆盖所有可能的用户（如果以 root 运行）
if [ "$(id -u)" -eq 0 ]; then
    USERS_DIRS="/root $(ls -d /home/* 2>/dev/null)"
else
    USERS_DIRS="$HOME"
fi

for user_dir in $USERS_DIRS; do
    USER_HISTORY_FILES="
    $user_dir/.bash_history
    $user_dir/.zsh_history
    $user_dir/.python_history
    $user_dir/.viminfo
    $user_dir/.lesshst
    $user_dir/.mysql_history
    $user_dir/.rediscli_history
    $user_dir/.bash_logout
    "

    for file in $USER_HISTORY_FILES; do
        if [ -f "$file" ]; then
            echo "Clearing $file..."
            true > "$file"
            # 尝试修改权限或使用 chattr 确保不容易被 shell 进程在退出时写回（可选）
        fi
    done
done

# 3. 清理系统级日志文件 (需要 root 权限)
if [ "$(id -u)" -eq 0 ]; then
    # 系统审计与登陆历史
    SYS_LOG_FILES="
    /var/log/wtmp
    /var/log/btmp
    /var/log/lastlog
    /var/log/secure
    /var/log/auth.log
    /var/log/messages
    /var/log/syslog
    /var/log/maillog
    /var/run/utmp
    /var/log/audit/audit.log
    /var/log/cron
    /var/log/faillog
    /var/log/tallylog
    "

    for log in $SYS_LOG_FILES; do
        if [ -f "$log" ]; then
            echo "Truncating $log..."
            true > "$log"
        fi
    done

    # 清理 systemd 日志 (journalctl)
    if command -v journalctl >/dev/null 2>&1; then
        echo "Vacuuming journalctl..."
        journalctl --vacuum-time=1s >/dev/null 2>&1
        # 彻底清空（如果支持 --rotate）
        journalctl --rotate >/dev/null 2>&1
        journalctl --vacuum-time=1s >/dev/null 2>&1
    fi

    # 清理内核环形缓冲区 (dmesg)
    dmesg -c >/dev/null 2>&1

    # 清理邮件记录
    if [ -d "/var/spool/mail" ]; then
        for mailfile in /var/spool/mail/*; do
            if [ -f "$mailfile" ]; then
                true > "$mailfile"
            fi
        done
    fi
else
    echo "Warning: Not running as root. System logs were not cleared."
fi

# 4. 彻底禁用当前环境的历史记录 (仅对脚本进程有效)
export HISTSIZE=0
export HISTFILESIZE=0
unset HISTFILE

echo "Successfully cleared files. IMPORTANT: To clear your CURRENT terminal memory, please run: history -c && history -w"
