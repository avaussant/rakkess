/*
Copyright 2019 Cornelius Weig

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rakkess

import (
	"context"

	"github.com/corneliusweig/rakkess/pkg/rakkess/client"
	"github.com/corneliusweig/rakkess/pkg/rakkess/options"
	"github.com/corneliusweig/rakkess/pkg/rakkess/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Rakkess(ctx context.Context, opts *options.RakkessOptions) error {
	if err := util.ValidateVerbs(opts.Verbs); err != nil {
		return err
	}
	if err := util.ValidateOutputFormat(opts.Output); err != nil {
		return err
	}

	grs, err := client.FetchAvailableGroupResources(opts)
	if err != nil {
		return errors.Wrap(err, "fetch available group resources")
	}
	logrus.Info(grs)

	authClient, err := opts.GetAuthClient()
	if err != nil {
		return errors.Wrap(err, "get auth client")
	}

	namespace := opts.ConfigFlags.Namespace
	results, err := client.CheckResourceAccess(ctx, authClient, grs, opts.Verbs, namespace)
	if err != nil {
		return errors.Wrap(err, "check resource access")
	}

	util.PrintResults(opts.Streams.Out, opts.Verbs, outputFormat(opts), results)

	if namespace == nil || *namespace == "" {
		logrus.Warn("No namespace given, this implies cluster scope (try -n if this is not intended)")
	}

	return nil
}

func outputFormat(o *options.RakkessOptions) util.OutputFormat {
	if o.Output == "ascii-table" {
		return util.ASCIITable
	}
	return util.IconTable
}
