// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package plugin

import (
	"github.com/spf13/cobra"
)

func getPluginInvokedAs(descriptor *PluginDescriptor) string {
	name := descriptor.Name

	for _, entry := range descriptor.CommandMap {
		// this is a plugin-level map
		if entry.SourceCommandPath == "" {
			name = entry.DestinationCommandPath
		}
	}

	return name
}

func newRootCmd(descriptor *PluginDescriptor) *cobra.Command {
	cmdName := getPluginInvokedAs(descriptor)

	cmd := &cobra.Command{
		Use:     descriptor.Name,
		Short:   descriptor.Description,
		Aliases: descriptor.Aliases,
		// Disable footers in docs generated
		DisableAutoGenTag: true,
		// Hide the default completion command of the plugin.
		// Shell completion is enabled using the Tanzu CLI's `completion` command so a plugin
		// does not need its own `completion` command.  Having such a command is just
		// confusing for users. However, we don't disable it completely for two reasons:
		//   1. backwards-compatibility, as the command used to be available for some plugins
		//   2. to allow shell completion when using the plugin as a native program (mostly for testing)
		// Note that a plugin can completely disable this command itself using:
		//  plugin.Cmd.CompletionOptions.DisableDefaultCmd = true
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		Annotations: map[string]string{
			"target":                           string(descriptor.Target),
			cobra.CommandDisplayNameAnnotation: cmdName,
		},
	}
	// Instead of using templates, use go functions.
	// This allows for dead-code-elimination.
	// The below call will set the format for both usage and help printouts.
	cmd.SetUsageFunc(UsageFunc)
	// Keep this call for backwards-compatibility, in case a plugin uses the templating
	cobra.AddTemplateFuncs(TemplateFuncs)

	cmd.AddCommand(
		newDescribeCmd(descriptor.Description),
		newVersionCmd(descriptor.Version),
		newInfoCmd(descriptor),
	)

	return cmd
}
