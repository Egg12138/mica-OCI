#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <errno.h>
#include <fcntl.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <signal.h>
#include <sys/wait.h>

#include "mica_shim.h"

#define MICA_SOCKET_PATH "/var/run/micad.sock"
#define MICA_RUNTIME_PATH "/usr/local/bin/mica-runtime"
#define MICA_WORK_DIR "/var/run/mica"

static int mica_socket_connect(mica_shim_ctx_t *ctx)
{
    struct sockaddr_un addr;
    int ret = 0;

    if (ctx == NULL) {
        return -1;
    }

    ctx->socket_fd = socket(AF_UNIX, SOCK_STREAM, 0);
    if (ctx->socket_fd < 0) {
        return -1;
    }

    memset(&addr, 0, sizeof(addr));
    addr.sun_family = AF_UNIX;
    strncpy(addr.sun_path, ctx->config.socket_path ? ctx->config.socket_path : MICA_SOCKET_PATH, 
            sizeof(addr.sun_path) - 1);

    ret = connect(ctx->socket_fd, (struct sockaddr *)&addr, sizeof(addr));
    if (ret < 0) {
        close(ctx->socket_fd);
        return -1;
    }

    ctx->connected = true;
    return 0;
}

static void mica_socket_disconnect(mica_shim_ctx_t *ctx)
{
    if (ctx && ctx->connected) {
        close(ctx->socket_fd);
        ctx->connected = false;
    }
}

int mica_shim_init(mica_shim_ctx_t *ctx, const mica_shim_config_t *config)
{
    if (ctx == NULL || config == NULL) {
        return -1;
    }

    memset(ctx, 0, sizeof(mica_shim_ctx_t));
    ctx->config = *config;

    // Create work directory if not exists
    if (access(MICA_WORK_DIR, F_OK) != 0) {
        if (mkdir(MICA_WORK_DIR, 0755) != 0) {
            return -1;
        }
    }

    return mica_socket_connect(ctx);
}

void mica_shim_cleanup(mica_shim_ctx_t *ctx)
{
    if (ctx == NULL) {
        return;
    }

    mica_socket_disconnect(ctx);
    free(ctx->config.socket_path);
    free(ctx->config.runtime_path);
    free(ctx->config.namespace);
    free(ctx->config.container_id);
    free(ctx->config.bundle);
    free(ctx->config.work_dir);
}

int mica_shim_create(mica_shim_ctx_t *ctx, const char *bundle, const char *runtime)
{
    if (ctx == NULL || bundle == NULL) {
        return -1;
    }

    // Store bundle path
    ctx->config.bundle = strdup(bundle);
    if (ctx->config.bundle == NULL) {
        return -1;
    }

    // Store runtime path if provided
    if (runtime != NULL) {
        ctx->config.runtime_path = strdup(runtime);
        if (ctx->config.runtime_path == NULL) {
            return -1;
        }
    }

    return 0;
}

int mica_shim_start(mica_shim_ctx_t *ctx)
{
    pid_t pid;
    char *argv[] = {
        ctx->config.runtime_path ? ctx->config.runtime_path : MICA_RUNTIME_PATH,
        "start",
        ctx->config.container_id,
        ctx->config.bundle,
        NULL
    };

    if (ctx == NULL) {
        return -1;
    }

    pid = fork();
    if (pid < 0) {
        return -1;
    }

    if (pid == 0) {
        // Child process
        execv(argv[0], argv);
        exit(1);
    }

    // Parent process
    ctx->container_pid = pid;
    ctx->running = true;
    return 0;
}

int mica_shim_stop(mica_shim_ctx_t *ctx)
{
    if (ctx == NULL || !ctx->running) {
        return -1;
    }

    if (kill(ctx->container_pid, SIGTERM) != 0) {
        return -1;
    }

    // Wait for process to terminate
    waitpid(ctx->container_pid, NULL, 0);
    ctx->running = false;
    return 0;
}

int mica_shim_delete(mica_shim_ctx_t *ctx)
{
    if (ctx == NULL) {
        return -1;
    }

    if (ctx->running) {
        mica_shim_stop(ctx);
    }

    return 0;
}

int mica_shim_state(mica_shim_ctx_t *ctx, char *state, size_t max_len)
{
    if (ctx == NULL || state == NULL || max_len == 0) {
        return -1;
    }

    if (!ctx->running) {
        strncpy(state, "stopped", max_len - 1);
    } else {
        strncpy(state, "running", max_len - 1);
    }
    state[max_len - 1] = '\0';

    return 0;
}

int mica_shim_exec(mica_shim_ctx_t *ctx, const char *exec_id, const char *command)
{
    if (ctx == NULL || !ctx->running || exec_id == NULL || command == NULL) {
        return -1;
    }

    // TODO: Implement command execution in container
    return 0;
}

int mica_shim_kill(mica_shim_ctx_t *ctx, pid_t pid, int signal)
{
    if (ctx == NULL || !ctx->running) {
        return -1;
    }

    if (pid == 0) {
        pid = ctx->container_pid;
    }

    return kill(pid, signal);
}

int mica_shim_list_pids(mica_shim_ctx_t *ctx, pid_t *pids, size_t max_pids)
{
    if (ctx == NULL || pids == NULL || max_pids == 0) {
        return -1;
    }

    if (!ctx->running) {
        return 0;
    }

    // Only return container PID for now
    pids[0] = ctx->container_pid;
    return 1;
}

int mica_shim_stats(mica_shim_ctx_t *ctx, char *stats, size_t max_len)
{
    if (ctx == NULL || stats == NULL || max_len == 0) {
        return -1;
    }

    // TODO: Implement container stats collection
    strncpy(stats, "{}", max_len - 1);
    stats[max_len - 1] = '\0';
    return 0;
}

int mica_shim_update(mica_shim_ctx_t *ctx, const char *resources)
{
    if (ctx == NULL || !ctx->running || resources == NULL) {
        return -1;
    }

    // TODO: Implement resource update
    return 0;
}

int mica_shim_pause(mica_shim_ctx_t *ctx)
{
    if (ctx == NULL || !ctx->running) {
        return -1;
    }

    return kill(ctx->container_pid, SIGSTOP);
}

int mica_shim_resume(mica_shim_ctx_t *ctx)
{
    if (ctx == NULL || !ctx->running) {
        return -1;
    }

    return kill(ctx->container_pid, SIGCONT);
} 
