// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/ldap"
	ldapClient "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/ldap"
	mcuimodels "github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/web/server/models"
)

// VerifyLdapBind verifies LDAP authentication.
func (app *App) VerifyLdapBind(params ldap.VerifyLdapBindParams) middleware.Responder {
	if app.ldapClient == nil {
		return ldap.NewVerifyLdapBindInternalServerError().WithPayload(Err(errors.New("LDAP client is not initialized properly")))
	}

	result, err := app.ldapClient.LdapBind()
	if err != nil {
		return ldap.NewVerifyLdapBindBadRequest().WithPayload(Err(errors.Wrap(err, "unable to perform LDAP bind")))
	}

	return ldap.NewVerifyLdapBindOK().WithPayload(ldapTestResultToMCUILdapTestResult(result))
}

// VerifyLdapCloseConnection disconnect from a LDAP server.
func (app *App) VerifyLdapCloseConnection(params ldap.VerifyLdapCloseConnectionParams) middleware.Responder {
	if app.ldapClient == nil {
		return ldap.NewVerifyLdapCloseConnectionInternalServerError().WithPayload(Err(errors.New("LDAP client is not initialized properly")))
	}

	app.ldapClient.LdapCloseConnection()

	return ldap.NewVerifyLdapCloseConnectionCreated()
}

// VerifyLdapConnect checks LDAP server can be reached.
func (app *App) VerifyLdapConnect(params ldap.VerifyLdapConnectParams) middleware.Responder {
	app.ldapClient = ldapClient.New()
	success, err := app.ldapClient.LdapConnect(mcUILdapParamsToLdapParams(params.Credentials))

	if err != nil {
		return ldap.NewVerifyLdapConnectBadRequest().WithPayload(Err(errors.Wrap(err, "unable to connect to LDAP server as configed")))
	}

	return ldap.NewVerifyLdapConnectOK().WithPayload(ldapTestResultToMCUILdapTestResult(success))
}

// VerifyGroupSearch verifies the LDAP group can be found.
func (app *App) VerifyGroupSearch(params ldap.VerifyLdapGroupSearchParams) middleware.Responder {
	if app.ldapClient == nil {
		return ldap.NewVerifyLdapGroupSearchInternalServerError().WithPayload(Err(errors.New("LDAP client is not initialized properly")))
	}

	success, err := app.ldapClient.LdapGroupSearch()
	if err != nil {
		return ldap.NewVerifyLdapGroupSearchBadRequest().WithPayload(Err(errors.Wrap(err, "unable to perform LDAP Group Search")))
	}

	return ldap.NewVerifyLdapGroupSearchOK().WithPayload(ldapTestResultToMCUILdapTestResult(success))
}

// VerifyUserSearch verifies the LDAP user can be found.
func (app *App) VerifyUserSearch(params ldap.VerifyLdapUserSearchParams) middleware.Responder {
	if app.ldapClient == nil {
		return ldap.NewVerifyLdapUserSearchInternalServerError().WithPayload(Err(errors.New("LDAP client is not initialized properly")))
	}

	success, err := app.ldapClient.LdapUserSearch()
	if err != nil {
		return ldap.NewVerifyLdapUserSearchBadRequest().WithPayload(Err(errors.Wrap(err, "unable to perform LDAP User Search")))
	}

	return ldap.NewVerifyLdapUserSearchOK().WithPayload(ldapTestResultToMCUILdapTestResult(success))
}

// Need until the ldapClient code decouples from the presentation code in the
// management-cluster API logic.
func ldapTestResultToMCUILdapTestResult(tr *mcuimodels.LdapTestResult) *models.LdapTestResult {
	result := &models.LdapTestResult{
		Code: tr.Code,
		Desc: tr.Desc,
	}

	return result
}

// Need until the ldapClient code decouples from the presentation code in the
// management-cluster API logic.
func mcUILdapParamsToLdapParams(p *models.LdapParams) *mcuimodels.LdapParams {
	result := &mcuimodels.LdapParams{
		LdapBindDn:               p.LdapBindDn,
		LdapBindPassword:         p.LdapBindPassword,
		LdapGroupSearchBaseDn:    p.LdapGroupSearchBaseDn,
		LdapGroupSearchFilter:    p.LdapGroupSearchFilter,
		LdapGroupSearchGroupAttr: p.LdapGroupSearchGroupAttr,
		LdapGroupSearchNameAttr:  p.LdapGroupSearchNameAttr,
		LdapGroupSearchUserAttr:  p.LdapGroupSearchUserAttr,
		LdapRootCa:               p.LdapRootCa,
		LdapTestGroup:            p.LdapTestGroup,
		LdapTestUser:             p.LdapTestUser,
		LdapURL:                  p.LdapURL,
		LdapUserSearchBaseDn:     p.LdapUserSearchBaseDn,
		LdapUserSearchEmailAttr:  p.LdapUserSearchEmailAttr,
		LdapUserSearchFilter:     p.LdapUserSearchFilter,
		LdapUserSearchIDAttr:     p.LdapUserSearchIDAttr,
		LdapUserSearchNameAttr:   p.LdapUserSearchNameAttr,
		LdapUserSearchUsername:   p.LdapUserSearchUsername,
	}

	return result
}
