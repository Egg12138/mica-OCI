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
 * Create: 2025-05-20
 * Description: provide mica shim implementation
 ******************************************************************************/

#ifndef MICA_SHIM_H
#define MICA_SHIM_H

#include <stdbool.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

// Mica shim configuration
typedef struct {
    char *socket_path;
    char *runtime_path;
    char *namespace;
    char *container_id;
    char *bundle;
    char *work_dir;
} mica_shim_config_t;

// Mica shim context
typedef struct {
    mica_shim_config_t config;
    int socket_fd;
    bool connected;
    pid_t container_pid;
    bool running;
} mica_shim_ctx_t;

// Initialize mica shim context
int mica_shim_init(mica_shim_ctx_t *ctx, const mica_shim_config_t *config);

// Cleanup mica shim context
void mica_shim_cleanup(mica_shim_ctx_t *ctx);

// Create container
int mica_shim_create(mica_shim_ctx_t *ctx, const char *bundle, const char *runtime);

// Start container
int mica_shim_start(mica_shim_ctx_t *ctx);

// Stop container
int mica_shim_stop(mica_shim_ctx_t *ctx);

// Delete container
int mica_shim_delete(mica_shim_ctx_t *ctx);

// Get container state
int mica_shim_state(mica_shim_ctx_t *ctx, char *state, size_t max_len);

// Execute command in container
int mica_shim_exec(mica_shim_ctx_t *ctx, const char *exec_id, const char *command);

// Kill process in container
int mica_shim_kill(mica_shim_ctx_t *ctx, pid_t pid, int signal);

// List processes in container
int mica_shim_list_pids(mica_shim_ctx_t *ctx, pid_t *pids, size_t max_pids);

// Get container stats
int mica_shim_stats(mica_shim_ctx_t *ctx, char *stats, size_t max_len);

// Update container resources
int mica_shim_update(mica_shim_ctx_t *ctx, const char *resources);

// Pause container
int mica_shim_pause(mica_shim_ctx_t *ctx);

// Resume container
int mica_shim_resume(mica_shim_ctx_t *ctx);

#ifdef __cplusplus
}
#endif

#endif // MICA_SHIM_H 