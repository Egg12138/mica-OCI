package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	v1opts "github.com/containerd/containerd/pkg/runtimeoptions/v1"
)

var validShims = [2]string{
	"io.containerd.runc.v2",
	"io.containerd.runtime.v1.linux",
}

func isValidShim(shim string) bool {
	for _, valid := range validShims {
		if shim == valid {
			return true
		}
	}
	return false
}

func getSpecOpts(_ context.Context, _ containerd.Container, spec *oci.Spec) []oci.SpecOpts {
	var opts []oci.SpecOpts
	if spec.Process != nil {
		opts = append(opts, oci.WithProcessArgs(spec.Process.Args...))
	}
	return opts
}

func migrateContainerShim(
	ctx context.Context,
	client *containerd.Client,
	container containerd.Container,
	newShim string,
) (containerd.Container, error) {
	info, err := container.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container info: %w", err)
	}

	img, err := container.Image(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container image: %w", err)
	}

	spec, err := container.Spec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container spec: %w", err)
	}

	if err := container.Delete(ctx); err != nil {
		return nil, fmt.Errorf("failed to delete old container: %w", err)
	}

	specOpts := getSpecOpts(ctx, container, spec)
	shimOpts := info.Runtime.Options

	newContainer, err := client.NewContainer(
		ctx, info.ID,
		containerd.WithRuntime(newShim, shimOpts),
		containerd.WithImage(img),
		containerd.WithSnapshotter(info.Snapshotter),
		containerd.WithSnapshot(info.SnapshotKey),
		containerd.WithNewSpec(append(specOpts, oci.WithImageConfig(img))...),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new container with new shim: %w", err)
	}

	log.Printf("Successfully migrated container to new shim: %s\n", newShim)
	return newContainer, nil
}

func main() {
	ctx := namespaces.WithNamespace(context.TODO(), "default")
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Fatalf("Failed to connect to containerd: %v", err)
	}
	defer client.Close()

	imgs, err := client.ListImages(ctx)
	if err != nil {
		log.Fatalf("Failed to list images: %v", err)
	}

	if len(imgs) == 0 {
		log.Println("No images found. Attempting to import from stdin...")
		importedImgs, err := client.Import(ctx, os.Stdin)
		if err != nil {
			log.Fatalf("Failed to import images: %v", err)
		}
		for _, i := range importedImgs {
			img, err := client.GetImage(ctx, i.Name)
			if err != nil {
				log.Fatalf("Failed to get imported image %s: %v", i.Name, err)
			}
			if err := img.Unpack(ctx, "overlayfs"); err != nil {
				log.Fatalf("Failed to unpack image %s: %v", i.Name, err)
			}
			log.Printf("Imported image: %s (%s)", i.Name, i.Target.Platform.OS)
		}
	} else {
		log.Println("Available images:")
		log.Println("NAME           SIZE")
		for _, i := range imgs {
			lastName := strings.TrimSuffix(i.Name(), "/")
			lastName = lastName[strings.LastIndex(lastName, "/")+1:]
		sz, _ := i.Size(ctx)
		mb := float64(sz) / 1024 / 1024
		log.Printf("%s %.2f MB", lastName, mb)
	}
	}

	img, err := client.GetImage(ctx, "docker.io/library/busybox:latest")
	if err != nil {
		log.Fatalf("Failed to get busybox image: %v", err)
	}

	var opts v1opts.Options
	const myContainerName = "myContainer"

	cntr, err := client.NewContainer(
		ctx,
		myContainerName,
		containerd.WithSnapshotter("overlayfs"),
		containerd.WithNewSnapshot("myContainer-snapshot", img),
		containerd.WithImage(img),
		containerd.WithNewSpec(oci.WithImageConfig(img)),
		containerd.WithRuntime("io.containerd.runc.v2", &opts),
	)
	if err != nil {
		if !errdefs.IsAlreadyExists(err) {
			log.Fatalf("Failed to create container: %v", err)
		}

		log.Printf("Container '%s' already exists, loading it...", myContainerName)
		container, err := client.LoadContainer(ctx, myContainerName)
		if err != nil {
			log.Fatalf("Failed to load existing container: %v", err)
		}

		info, err := container.Info(ctx)
		if err != nil {
			log.Fatalf("Failed to get container info: %v", err)
		}

		snapshotter := client.SnapshotService(info.Snapshotter)
		_, err = snapshotter.Stat(ctx, "myContainer-snapshot")
		if err != nil && !errdefs.IsNotFound(err) {
			log.Fatalf("Failed to stat snapshot: %v", err)
		}

		if errdefs.IsNotFound(err) {
			log.Println("Creating new container with a new snapshot...")
			cntr, err = client.NewContainer(
				ctx,
				myContainerName,
				containerd.WithSnapshotter(info.Snapshotter),
				containerd.WithNewSnapshot("myContainer-snapshot", img),
				containerd.WithImage(img),
				containerd.WithNewSpec(oci.WithImageConfig(img)),
				containerd.WithRuntime("io.containerd.runc.v2", &opts),
			)
			if err != nil {
				log.Fatalf("Failed to recreate container with new snapshot: %v", err)
			}
		} else {
			log.Println("Reusing existing container...")
			cntr, err = client.LoadContainer(ctx, myContainerName)
			if err != nil {
				log.Fatalf("Failed to reload container: %v", err)
			}
		}
	}

	info, _ := cntr.Info(ctx)
	if !isValidShim(info.Runtime.Name) {
		log.Println("Container is not using a supported shim, migrating...")
		newContainer, err := migrateContainerShim(ctx, client, cntr, "io.containerd.runc.v2")
		if err != nil {
			log.Fatalf("Failed to migrate container: %v", err)
		}
		cntr = newContainer
	}

	spec, _ := cntr.Spec(ctx)
    spec.Process.Args = []string{"/bin/bash", "-c", "echo HEllo world"}
    spec.Process.User.UID = 0
	log.Printf("Container Info - Version: %s, Domain: %s, RootFS: %s, Runtime: %s",
		spec.Version, spec.Domainname, spec.Root.Path, info.Runtime.Name,
	)

	// Attach or create new task
	log.Println("Attaching to existing task or creating a new one...")

	// Create task with stdio
	task, err := cntr.Task(ctx, nil)
	if err != nil {
		if errdefs.IsNotFound(err) {
			// No existing task, create a new one
			task, err = cntr.NewTask(ctx, cio.LogFile(fmt.Sprintf("/tmp/%s-logs.txt", cntr.ID())))
			if err != nil {
				log.Fatalf("Failed to create new task: %v", err)
			}
		} else {
			log.Fatalf("Failed to attach to task: %v", err)
		}
	}
    proc, err := task.Exec(
        ctx, 
        fmt.Sprintf("%d", spec.Process.User.UID), 
        spec.Process, 
        cio.NewCreator(cio.WithStdio),
    )

    if err != nil {
        log.Fatalf("Failed to create new process: %v", err)
    }

    if err := proc.Start(ctx); err != nil {
        log.Fatalf("Failed to start process: %v", err)
    }

    statusChannel, err := proc.Wait(ctx)
    if err != nil {
        log.Fatalf("Failed to wait for process: %v", err)
    }
    status := <-statusChannel
    log.Printf("$%d[%s]", status.ExitCode(), status.ExitTime())
    task.Pause(ctx)

	log.Printf("Task PID: %d\n", task.Pid())
	cntr.Delete(ctx)
}