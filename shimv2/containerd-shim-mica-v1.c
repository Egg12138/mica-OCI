/******************************************************************************
 * Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved.
 * iSulad licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
 * PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: xiaojunzhe
 * Create: 2024-05-20
 * Description: provide mica shim service implementation
 ******************************************************************************/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <signal.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <errno.h>

#include "mica_shim.h"

#define SHIM_SOCKET_PATH "/var/run/mica-shim.sock"
#define SHIM_NAMESPACE "default"
#define SHIM_CONTAINER_ID "test-container"

static mica_shim_ctx_t g_shim_ctx;
static bool g_running = true;

static void signal_handler(int signo)
{
    if (signo == SIGTERM || signo == SIGINT) {
        g_running = false;
    }
}

static int setup_signal_handlers(void)
{
    struct sigaction sa;
    memset(&sa, 0, sizeof(sa));
    sa.sa_handler = signal_handler;
    sigemptyset(&sa.sa_mask);

    if (sigaction(SIGTERM, &sa, NULL) < 0) {
        return -1;
    }
    if (sigaction(SIGINT, &sa, NULL) < 0) {
        return -1;
    }

    return 0;
}

static int handle_create_request(const char *bundle, const char *runtime)
{
    return mica_shim_create(&g_shim_ctx, bundle, runtime);
}

static int handle_start_request(void)
{
    return mica_shim_start(&g_shim_ctx);
}

static int handle_stop_request(void)
{
    return mica_shim_stop(&g_shim_ctx);
}

static int handle_delete_request(void)
{
    return mica_shim_delete(&g_shim_ctx);
}

static int handle_state_request(char *state, size_t max_len)
{
    return mica_shim_state(&g_shim_ctx, state, max_len);
}

static int handle_exec_request(const char *exec_id, const char *command)
{
    return mica_shim_exec(&g_shim_ctx, exec_id, command);
}

static int handle_kill_request(pid_t pid, int signal)
{
    return mica_shim_kill(&g_shim_ctx, pid, signal);
}

static int handle_list_pids_request(pid_t *pids, size_t max_pids)
{
    return mica_shim_list_pids(&g_shim_ctx, pids, max_pids);
}

static int handle_stats_request(char *stats, size_t max_len)
{
    return mica_shim_stats(&g_shim_ctx, stats, max_len);
}

static int handle_update_request(const char *resources)
{
    return mica_shim_update(&g_shim_ctx, resources);
}

static int handle_pause_request(void)
{
    return mica_shim_pause(&g_shim_ctx);
}

static int handle_resume_request(void)
{
    return mica_shim_resume(&g_shim_ctx);
}

int main(int argc, char *argv[])
{
    mica_shim_config_t config = {
        .socket_path = MICA_SOCKET_PATH,
        .runtime_path = NULL,
        .namespace = SHIM_NAMESPACE,
        .container_id = SHIM_CONTAINER_ID,
        .bundle = NULL,
        .work_dir = NULL
    };

    if (setup_signal_handlers() < 0) {
        fprintf(stderr, "Failed to setup signal handlers\n");
        return 1;
    }

    if (mica_shim_init(&g_shim_ctx, &config) < 0) {
        fprintf(stderr, "Failed to initialize mica shim\n");
        return 1;
    }

    // Main event loop
    while (g_running) {
        // TODO: Implement request handling logic
        sleep(1);
    }

    mica_shim_cleanup(&g_shim_ctx);
    return 0;
} 