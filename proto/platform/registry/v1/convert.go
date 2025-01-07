/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	apkotypes "chainguard.dev/apko/pkg/build/types"
)

func ToApkoProto(ic apkotypes.ImageConfiguration) *ApkoConfig {
	return &ApkoConfig{
		Contents: &ApkoConfig_Contents{
			Repositories:      ic.Contents.RuntimeRepositories,
			BuildRepositories: ic.Contents.BuildRepositories,
			Keyring:           ic.Contents.Keyring,
			Packages:          ic.Contents.Packages,
		},
		Environment: ic.Environment,
		Accounts: &ApkoConfig_Accounts{
			Users:  pbUsers(ic.Accounts.Users),
			Groups: pbGroups(ic.Accounts.Groups),
			RunAs:  ic.Accounts.RunAs,
		},
		Annotations: ic.Annotations,
		Paths:       pbPaths(ic.Paths),
		Entrypoint:  pbEntrypoint(ic.Entrypoint),
		Cmd:         ic.Cmd,
		WorkDir:     ic.WorkDir,
		Archs:       pbArchs(ic.Archs),

		// These are unused.
		Volumes:    ic.Volumes,
		StopSignal: ic.StopSignal,
		Include:    ic.Include,
		VcsUrl:     ic.VCSUrl,
	}
}

func pbArchs(archs []apkotypes.Architecture) []string {
	if archs == nil {
		return nil
	}
	pbArchs := make([]string, len(archs))
	for i, a := range archs {
		pbArchs[i] = a.String()
	}
	return pbArchs
}

func pbUsers(users []apkotypes.User) []*ApkoConfig_Accounts_User {
	if users == nil {
		return nil
	}
	pbUsers := make([]*ApkoConfig_Accounts_User, len(users))
	for i, u := range users {
		pbUsers[i] = &ApkoConfig_Accounts_User{
			Uid:      u.UID,
			UserName: u.UserName,
			Gid:      *u.GID,
		}
	}
	return pbUsers
}

func pbGroups(groups []apkotypes.Group) []*ApkoConfig_Accounts_Group {
	if groups == nil {
		return nil
	}
	pbGroups := make([]*ApkoConfig_Accounts_Group, len(groups))
	for i, g := range groups {
		pbGroups[i] = &ApkoConfig_Accounts_Group{
			Gid:       g.GID,
			GroupName: g.GroupName,
		}
	}
	return pbGroups
}

func pbPaths(paths []apkotypes.PathMutation) []*ApkoConfig_PathMutation {
	if paths == nil {
		return nil
	}
	pbPaths := make([]*ApkoConfig_PathMutation, len(paths))
	for i, p := range paths {
		pbPaths[i] = &ApkoConfig_PathMutation{
			Path:        p.Path,
			Type:        p.Type,
			Uid:         p.UID,
			Gid:         p.GID,
			Permissions: p.Permissions,
			Source:      p.Source,
			Recursive:   p.Recursive,
		}
	}
	return pbPaths
}

func pbEntrypoint(entrypoint apkotypes.ImageEntrypoint) *ApkoConfig_Entrypoint {
	return &ApkoConfig_Entrypoint{
		Type:          entrypoint.Type,
		Command:       entrypoint.Command,
		ShellFragment: entrypoint.ShellFragment,
		Services:      entrypoint.Services,
	}
}

func ToApkoNative(cfg *ApkoConfig) apkotypes.ImageConfiguration {
	if cfg == nil {
		return apkotypes.ImageConfiguration{}
	}

	return apkotypes.ImageConfiguration{
		Contents:    apkoContents(cfg.Contents),
		Environment: cfg.Environment,
		Accounts:    apkoAccounts(cfg.Accounts),
		Annotations: cfg.Annotations,
		Paths:       apkoPaths(cfg.Paths),
		Entrypoint:  apkoEntrypoint(cfg.Entrypoint),
		Cmd:         cfg.Cmd,
		WorkDir:     cfg.WorkDir,
		Archs:       apkoArchs(cfg.Archs),

		// These are unused.
		Volumes:    cfg.Volumes,
		StopSignal: cfg.StopSignal,
		Include:    cfg.Include,
		VCSUrl:     cfg.VcsUrl,
	}
}

func apkoContents(contents *ApkoConfig_Contents) apkotypes.ImageContents {
	if contents == nil {
		return apkotypes.ImageContents{}
	}

	return apkotypes.ImageContents{
		RuntimeRepositories: contents.Repositories,
		BuildRepositories:   contents.BuildRepositories,
		Keyring:             contents.Keyring,
		Packages:            contents.Packages,
	}
}

func apkoArchs(archs []string) []apkotypes.Architecture {
	if archs == nil {
		return nil
	}
	apkoArchs := make([]apkotypes.Architecture, len(archs))
	for i, a := range archs {
		apkoArchs[i] = apkotypes.ParseArchitecture(a)
	}
	return apkoArchs
}

func apkoAccounts(accounts *ApkoConfig_Accounts) apkotypes.ImageAccounts {
	if accounts == nil {
		return apkotypes.ImageAccounts{}
	}
	return apkotypes.ImageAccounts{
		Users:  apkoUsers(accounts.Users),
		Groups: apkoGroups(accounts.Groups),
		RunAs:  accounts.RunAs,
	}
}

func apkoUsers(users []*ApkoConfig_Accounts_User) []apkotypes.User {
	if users == nil {
		return nil
	}
	apkoUsers := make([]apkotypes.User, len(users))
	for i, u := range users {
		apkoUsers[i] = apkotypes.User{
			UID:      u.Uid,
			UserName: u.UserName,
			GID:      &u.Gid,
		}
	}
	return apkoUsers
}

func apkoGroups(groups []*ApkoConfig_Accounts_Group) []apkotypes.Group {
	if groups == nil {
		return nil
	}
	apkoGroups := make([]apkotypes.Group, len(groups))
	for i, g := range groups {
		apkoGroups[i] = apkotypes.Group{
			GID:       g.Gid,
			GroupName: g.GroupName,
		}
	}
	return apkoGroups
}

func apkoPaths(paths []*ApkoConfig_PathMutation) []apkotypes.PathMutation {
	if paths == nil {
		return nil
	}
	apkoPaths := make([]apkotypes.PathMutation, len(paths))
	for i, p := range paths {
		apkoPaths[i] = apkotypes.PathMutation{
			Path:        p.Path,
			Type:        p.Type,
			UID:         p.Uid,
			GID:         p.Gid,
			Permissions: p.Permissions,
			Source:      p.Source,
			Recursive:   p.Recursive,
		}
	}
	return apkoPaths
}

func apkoEntrypoint(entrypoint *ApkoConfig_Entrypoint) apkotypes.ImageEntrypoint {
	if entrypoint == nil {
		return apkotypes.ImageEntrypoint{}
	}
	return apkotypes.ImageEntrypoint{
		Type:          entrypoint.Type,
		Command:       entrypoint.Command,
		ShellFragment: entrypoint.ShellFragment,
		Services:      entrypoint.Services,
	}
}
