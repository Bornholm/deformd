package module_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"testing"
	"time"

	"github.com/Bornholm/deformd/internal/handler"
	"github.com/Bornholm/deformd/internal/handler/module"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/wneessen/go-mail"
)

type fakeSMTPContainer struct {
	testcontainers.Container
	URL      string
	SMTPPort int
}

func setupFakeSMTP(ctx context.Context) (*fakeSMTPContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "bornholm/fake-smtp:2022.11.222058",
		ExposedPorts: []string{"8080/tcp", "2525/tcp"},
		WaitingFor:   wait.ForHTTP("/").WithPort("8080/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	httpMappedPort, err := container.MappedPort(ctx, "8080")
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s:%s", ip, httpMappedPort.Port())

	smtpMappedPort, err := container.MappedPort(ctx, "2525")
	if err != nil {
		return nil, err
	}

	return &fakeSMTPContainer{Container: container, URL: url, SMTPPort: smtpMappedPort.Int()}, nil
}

type fakeSMTPResult struct {
	Data struct {
		Emails []struct {
			Subject string
			From    []struct {
				Name    string
				Address string
			}
		}
	}
}

func TestEmailModule(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Parallel()

	ctx := context.Background()

	fakeSMTPCtn, err := setupFakeSMTP(ctx)
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	// Clean up the container after the test is complete
	defer func() {
		if err := fakeSMTPCtn.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %v", errors.WithStack(err))
		}
	}()

	script, err := loadTestDataScript("email.js")
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	handler := handler.New(
		script,
		handler.WithModules(
			module.EmailModuleFactory(
				"localhost",
				mail.WithPort(fakeSMTPCtn.SMTPPort),
				mail.WithTLSPolicy(mail.TLSOpportunistic),
			),
		),
		handler.WithMaxDuration(5*time.Second),
	)

	form := url.Values{}

	if err := handler.Process(ctx, form); err != nil {
		t.Error(errors.WithStack(err))
	}

	res, err := http.Get(fakeSMTPCtn.URL + "/api/v1/emails?from=foo@bar.com")
	if err != nil {
		t.Fatal(errors.WithStack(err))
	}

	defer res.Body.Close()

	result := &fakeSMTPResult{}
	decoder := json.NewDecoder(res.Body)

	if err := decoder.Decode(&result); err != nil {
		t.Fatal(errors.WithStack(err))
	}

	lastEmail := result.Data.Emails[0]

	if e, g := "Test", lastEmail.Subject; e != g {
		t.Errorf("lastEmail.Subject: expected '%s', got '%s'", e, g)
	}

	if e, g := "foo@bar.com", lastEmail.From[0].Address; e != g {
		t.Errorf("lastEmail.From.Address: expected '%s', got '%s'", e, g)
	}
}

func loadTestDataScript(name string) (string, error) {
	filepath := filepath.Join("./testdata", name)

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", errors.Wrapf(err, "could not read file '%s'", filepath)
	}

	return string(data), nil
}
