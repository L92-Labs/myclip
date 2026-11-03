package main

import (
	"os"

	"github.com/egoist/mygo"

	"github.com/L92-Labs/myclip/internal/transfer"
	"github.com/L92-Labs/myclip/locales"
)

// Backup export and import: the prefs pane's two buttons land here. The
// dialogs block the main thread, like MyGo's other dialogs.

var backupFilter = []mygo.FileFilter{{Name: "MyClip backup", Extensions: []string{"json"}}}

// ExportBackup asks where to write a backup and writes settings + history;
// failures show in a dialog.
func (a *app) ExportBackup() error {
	if err := a.exportBackup(); err != nil {
		mygo.Dialog.Error(locales.BackupExportTitle(), err.Error())
	}
	return nil
}

func (a *app) exportBackup() error {
	path, err := mygo.Dialog.Save(mygo.SaveDialogOptions{
		Title:       locales.BackupExportTitle(),
		DefaultPath: "myclip-backup.json",
		Filters:     backupFilter,
	})
	if err != nil || path == "" { // "" is a canceled dialog
		return err
	}
	// The backup carries the whole history: owner-only, like the database.
	// An existing file is tightened before anything is written into it.
	if info, err := os.Stat(path); err == nil && info.Mode().Perm() != 0o600 {
		if err := os.Chmod(path, 0o600); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	err = transfer.Export(f, &a.settings, a.engine.Items())
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		os.Remove(path) // no half backups
	}
	return err
}

// ImportBackup asks for a backup file, confirms, then merges its clips
// and applies its settings; failures show in a dialog.
func (a *app) ImportBackup() error {
	if err := a.importBackup(); err != nil {
		mygo.Dialog.Error(locales.BackupImportTitle(), err.Error())
	}
	return nil
}

func (a *app) importBackup() error {
	paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
		Title:   locales.BackupImportTitle(),
		Filters: backupFilter,
	})
	if err != nil || len(paths) == 0 { // nil paths is a canceled dialog
		return err
	}
	f, err := os.Open(paths[0])
	if err != nil {
		return err
	}
	archive, err := transfer.Read(f)
	f.Close()
	if err != nil {
		return err
	}
	andSettings := ""
	if archive.Settings != nil {
		andSettings = locales.BackupImportSettings()
	}
	r, err := mygo.Dialog.Message(mygo.MessageOptions{
		Type:          mygo.MessageQuestion,
		Title:         locales.BackupImportTitle(),
		Message:       locales.BackupImportAsk(len(archive.Items), andSettings),
		Buttons:       []string{locales.BackupImport(), locales.BackupCancel()},
		DefaultButton: 1, // decline by default
		CancelButton:  1,
	})
	if err != nil || r.Button != 0 {
		return nil
	}
	// The backup's settings go in first, when it carries them: importing
	// then applies the restored rules to the imported clips, not the
	// outgoing ones (and applyRules never unpins, so a late install could
	// not undo the outgoing rules' pins).
	if archive.Settings != nil {
		a.settings = *archive.Settings
		a.Changed()
	}
	added, skipped, err := a.engine.ImportItems(archive.HistoryItems())
	if err != nil {
		return err
	}
	if archive.Settings != nil {
		// The open preferences window's controls mirror the old values:
		// reload them, or saving one pane would restore the stale rest.
		if a.prefs != nil {
			a.prefs.Load()
			a.fitPrefs()
		}
	}
	go mygo.Dialog.Message(mygo.MessageOptions{
		Type:    mygo.MessageInfo,
		Message: locales.BackupImportDone(added, skipped),
	})
	return nil
}
