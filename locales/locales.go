// Copyright 1999-2020. Plesk International GmbH. All rights reserved.

package locales

import (
	"fmt"
)

// L is ready locale object
var L *locale

type locale struct {
	lang string
	defs map[string]map[string]string
}

func newLocale(lang string, defs map[string]map[string]string) *locale {
	return &locale{
		lang: lang,
		defs: defs,
	}
}

func (l *locale) Get(def string, vars ...interface{}) string {
	if tr, ok := l.defs[l.lang][def]; ok {
		return fmt.Sprintf(tr, vars...)
	}

	if tr, ok := l.defs["en_US"][def]; ok {
		return fmt.Sprintf(tr, vars...)
	}

	return def
}

func init() {
	// TODO: Proper locale detections & locale support

	var defs map[string]map[string]string
	defs = map[string]map[string]string{
		"en_US": {
			"app.description":                   "Manage local applications",
			"app.deploy.cmd":                    "deploy [PATH]",
			"app.deploy.description":            "Deploy application in set PATH to DOMAIN on SERVER",
			"app.deploy.success":                "Application successfully deployed",
			"app.register.cmd":                  "register [SERVER] [DOMAIN] [PATH]",
			"app.register.description":          "Register a new appliation in the specified path",
			"app.register.success":              "New app registered",
			"app.register.features.flag":        "Features necessary for the application",
			"app.register.target.path.flag":     "Path to deploy to on the server relative to document root",
			"app.register.overwrite.flag":       "Overwrite configuration if it exists",
			"app.register.flag.invalid":         "--path must point to a directory",
			"app.register.flag.feature.unknown": "Unknown feature %s, skipping",

			"database.description":        "Manage databases on the server",
			"database.create.cmd":         "create [SERVER] [DOMAIN] [NAME]",
			"database.create.description": "Create a database of specified type on the server",
			"database.create.flag.type":   "Type of the database server: mysql (default), mssql, or postgresql",
			"database.create.success":     "Database successfully created",
			"database.delete.cmd":         "delete [SERVER] [NAME]",
			"database.delete.description": "Delete database from the server",
			"database.delete.success":     "Database successfully deleted",
			"database.list.cmd":           "list [SERVER] [DOMAIN]",
			"database.list.description":   "List databases on server, optionally filtered by domain",
			"database.deploy.cmd":         "upload [SERVER] [DOMAIN] [NAME] [FILE]",
			"database.deploy.description": "Upload an SQL dump file to the server and deploy it",
			"database.deploy.success":     "Database successfully deployed",

			"domain.description":        "Manage domains on the server",
			"domain.create.cmd":         "create [SERVER] [DOMAIN] [IPv4] [IPv6]",
			"domain.create.description": "Create a new domain on the server with specific IPv4 & IPv6 addresses",
			"domain.create.success":     "Domain successfully created",
			"domain.list.cmd":           "list [SERVER]",
			"domain.list.description":   "List domains on the specific server",
			"domain.delete.cmd":         "delete [SERVER] [DOMAIN...]",
			"domain.delete.description": "Delete domain(s) from the specified server",
			"domain.delete.success":     "Domain successfully deleted",
			"domain.reload.cmd":         "reload [SERVER]",
			"domain.reload.description": "Reload cached domains for the specified server",
			"domain.reload.success":     "Domains successfully reloaded",

			"files.upload.cmd":            "sync [SERVER] [DOMAIN] [FILE ...]",
			"files.upload.description":    "Upload files to the target domain",
			"files.upload.flag.overwrite": "Overwrite existing files",
			"files.upload.flag.dry-run":   "Do not upload files, only show actions",
			"files.upload.success":        "Files successfully uploaded",

			"server.description":              "Manage known servers",
			"server.delete.cmd":               "delete [IP ADDRESS|HOSTNAME ...]",
			"server.delete.description":       "Remove registered server(s) and flush API key(s)",
			"server.delete.success":           "Server %s removed",
			"server.list.cmd":                 "list",
			"server.list.description":         "List registered servers",
			"server.register.cmd":             "register [IP ADDRESS|HOSTNAME]",
			"server.register.description":     "Register a new server on this device",
			"server.register.ignore.ssl.flag": "Ignore SSL certificate mismatch",
			"server.register.success":         "Server successfully registered",
			"server.reload.cmd":               "reload [IP ADDRESS|HOSTNAME]",
			"server.reload.description":       "Reload cached server data",
			"server.reload.success":           "Server data successfully reloaded",
			"server.reauth.cmd":               "reauth [IP ADDRESS|HOSTNAME]",
			"server.reauth.description":       "Update server API key",
			"server.reauth.success":           "Server successfully re-authenticated",
			"server.login.cmd":                "login [SERVER]",
			"server.login.description":        "Get login link for the specified server",
			"server.login.success":            "Generated login link: %s",

			"errors.server.remove.failure":     "Could not remove server %s: %s",
			"errors.feature.not.supported":     "Feature %s is not supported on Windows",
			"errors.feature.unknown":           "Unknown feature %s",
			"errors.path.is.not.directory":     "%s does not point to a directory",
			"errors.path.is.directory":         "%s points to a directory",
			"errors.cannot.parse.config":       "Configuration %s cannot be parsed, falling back to default config",
			"errors.path.already.exists":       "%s already exists",
			"errors.ip.address.required":       "One IPv4/IPv6 address must be specified",
			"errors.ip.address.not.cached":     "Could not find specified IP address(es) in cache, use command \"servers reload %s\" to reload cached data",
			"errors.ip.address.not.found":      "Could not find specified IP address(es) on the server %s",
			"errors.server.not.found":          "Server %s is not registered, use command \"servers register\" to add it",
			"errors.domain.not.cached":         "Domain %s does not exist on server %s, use command \"servers reload %s\" to reload cached data",
			"errors.database.not.found":        "Database %s does not exist on server %s, use command \"databases create\" to add it",
			"errors.database.server.not.found": "Database server with type %s does not exist on server %s",
			"errors.domain.unknown":            "Domain %s does not exist on the server",
			"errors.ip.address.class.limit":    "Cannot assign more than one IPv4 and one IPv6 to the domain",
			"errors.abspath.failed":            "Could not read absolute path of %s: %s, skipping",
			"errors.stat.failed":               "Could not stat() path %s: %s, skipping",
			"errors.upload.failed":             "Failed to upload %s: %s",
			"errors.unknown.database.type":     "Unknown database type %s",
			"errors.execution.failed.generic":  "Command execution failed with %s",
			"errors.domain.already.exists":     "Domain %s already exists",
			"errors.mkdir.failed":              "Could not create directory %s over FTP: %s",
			"errors.stor.failed":               "Could not store file %s over FTP: %s",

			"debug.mkdir.success": "Successfully created directory %s on server",
			"debug.stor.success":  "Successfully stored file %s on server",
			"debug.cwd":           "Changing directory to %s",
			"debug.dir.skip":      "Skipping directory %s that exists on target",
			"debug.file.skip":     "Skipping file %s that exists on target (no overwrite specified)",

			"api.errors.cligate.error.responce": "Execution failed with error code \"%d\", stdout: [%s], stderr: [%s]",
			"api.errors.failed.request":         "Request failed with code \"%d\": Reason: \"%s\"; Errors: \"%s\"",
			"api.errors.auth.wrong.pass":        "Could not authenticate on %s using provided password",
			"api.errors.auth.failed.reauth":     "Could not authenticate using stored credentials, use \"servers reauth %s\" to fix",
			"api.errors.auth.failed":            "Failed to acquire an API key using provided password: %s",
			"api.errors.auth.cli.failed":        "Failed to acquire an API key using provided password: [%d: %s]",
			"api.errors.domain.info.failed":     "Failed to get domain info: [%d: %s]",
			"api.errors.domain.info.not.found":  "Failed to get domain info (no FTP Login field)",

			"upload.dry.run.upload": "Dry run; would upload file %s to %s",
		},
	}

	L = newLocale("en_US", defs)
}
