package databasetest

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"testing"
)

func startContainer(t *testing.T) *container {
	t.Helper()

	cmd := exec.Command("docker", "run", "-d", "-p", "5432:5432", "-e", "POSTGRES_PASSWORD=postgres", "postgres:15.2-alpine")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		t.Fatalf("could not starting container: %s", err)
	}

	id := out.String()[:12]
	t.Log("DB container id:", id)
	cmd = exec.Command("docker", "inspect", id)
	out.Reset()
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not inspect container: %s: %v", id, err)
	}

	var doc []struct {
		NetworkSettings struct {
			Ports struct {
				P5432 []struct {
					HostIP   string `json:"HostIp"`
					HostPort string `json:"HostPort"`
				} `json:"5432/tcp"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}

	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatalf("could not decode inspect response: %s", err)
	}

	network := doc[0].NetworkSettings.Ports.P5432[0]

	c := container{
		ID:   id,
		Host: network.HostIP + ":" + network.HostPort,
	}

	t.Log("DB container host:", c.Host)

	return &c
}

type container struct {
	ID   string
	Host string
}

func stopContainer(t *testing.T, c *container) {
	t.Helper()

	cmd := exec.Command("docker", "stop", c.ID)
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not stop container: %s: %v", c.ID, err)
	}
	t.Log("stopped:", c.ID)

	cmd = exec.Command("docker", "rm", c.ID)
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not remove container: %s: %v", c.ID, err)
	}
	t.Log("removed:", c.ID)
}
