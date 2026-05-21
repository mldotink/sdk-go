//go:build e2e

package e2e

import (
	"context"
	"testing"
	"time"

	ink "github.com/mldotink/sdk-go"
)

func TestServiceVolume(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Minute)
	defer cancel()

	svcName := name("volsvc")
	volumeName := name("vol")
	result, err := client.CreateService(ctx, ink.CreateServiceInput{
		Name:          svcName,
		Source:        "image",
		Image:         "nginx:latest",
		Memory:        "256Mi",
		VCPUs:         "0.25",
		Ports:         publicHTTPPort(80),
		Project:       project,
		WorkspaceSlug: workspace,
		Volumes: []ink.VolumeSpec{{
			Name:      volumeName,
			MountPath: "/data",
			SizeGi:    1,
		}},
	})
	requireNoError(t, err, "create service with volume")
	cleanupServiceAndVolumes(t, result.ServiceID, svcName, project, volumeName)
	t.Logf("created service %s with volume %s", result.Name, volumeName)

	waitForActive(t, result.ServiceID, 3*time.Minute)
	vol := waitForVolumeStatus(t, volumeName, project, "provisioned", time.Minute)
	requireEqual(t, vol.MountPath, "/data", "volume mount path")
	requireEqual(t, vol.SizeGi, 1, "volume size")

	_, err = client.DeleteService(ctx, ink.DeleteServiceInput{
		ServiceID:     result.ServiceID,
		Project:       project,
		WorkspaceSlug: workspace,
	})
	requireNoError(t, err, "delete service")
	waitForServiceGone(t, result.ServiceID, time.Minute)
	vol = waitForVolumeStatus(t, volumeName, project, "detached", time.Minute)
	requireEqual(t, vol.Name, volumeName, "detached volume name")

	err = client.DeleteVolume(ctx, volumeName, project, workspace)
	requireNoError(t, err, "delete detached volume")
	waitForVolumeGone(t, volumeName, project, time.Minute)
}
