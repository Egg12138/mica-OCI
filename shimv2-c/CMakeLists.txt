cmake_minimum_required(VERSION 3.10)
project(mica-shim C)

set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED ON)

# Add source files
set(SOURCES
    mica_shim.c
    containerd-shim-mica-v1.c
)

# Add header files
set(HEADERS
    mica_shim.h
)

# Create executable
add_executable(containerd-shim-mica-v1 ${SOURCES} ${HEADERS})

# Set install path
install(TARGETS containerd-shim-mica-v1
    RUNTIME DESTINATION bin
)

# Add compile options
target_compile_options(containerd-shim-mica-v1 PRIVATE
    -Wall
    -Wextra
    -Werror
)

# Add include directories
target_include_directories(containerd-shim-mica-v1 PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}
) 