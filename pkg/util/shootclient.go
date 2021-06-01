package util

import (
	"context"

	gutil "github.com/gardener/gardener/extensions/pkg/util"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
	"github.com/gardener/gardener/pkg/utils/secrets"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewClientForShoot gets a client for a specific shoot. It wraps
// https://pkg.go.dev/github.com/gardener/gardener/extensions/pkg/util#NewClientForShoot
// and, if the connection is not resolved, attempts to use the external one
func NewClientForShoot(ctx context.Context, c client.Client, namespace string, opts client.Options) (*rest.Config, client.Client, error) {
	config, shoot, err := gutil.NewClientForShoot(ctx, c, namespace, opts)
	if err == nil {
		return config, shoot, err
	}
	// err was not nil, so try external
	return NewClientForShootExternal(ctx, c, namespace, opts)

}

// NewClientForShoot gets a client for a specific shoot. It is almost identical to
// https://pkg.go.dev/github.com/gardener/gardener/extensions/pkg/util#NewClientForShoot
// except that it explicitly use the external/public secret rather than internal
func NewClientForShootExternal(ctx context.Context, c client.Client, namespace string, opts client.Options) (*rest.Config, client.Client, error) {
	var (
		gardenerSecret = &corev1.Secret{}
		err            error
	)

	err = c.Get(ctx, kutil.Key(namespace, v1beta1constants.SecretNameGardener), gardenerSecret)
	if err != nil {
		return nil, nil, err
	}

	shootRESTConfig, err := gutil.NewRESTConfigFromKubeconfig(gardenerSecret.Data[secrets.DataKeyKubeconfig])
	if err != nil {
		return nil, nil, err
	}
	shootClient, err := client.New(shootRESTConfig, opts)
	if err != nil {
		return nil, nil, err
	}
	return shootRESTConfig, shootClient, nil
}
