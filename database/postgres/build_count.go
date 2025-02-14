// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// GetBuildCount gets a count of all builds from the database.
func (c *client) GetBuildCount() (int64, error) {
	c.Logger.Trace("getting count of builds from the database")

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := c.Postgres.
		Table(constants.TableBuild).
		Raw(dml.SelectBuildsCount).
		Pluck("count", &b).Error

	return b, err
}

// GetBuildCountByStatus gets a count of all builds by status from the database.
func (c *client) GetBuildCountByStatus(status string) (int64, error) {
	c.Logger.Tracef("getting count of builds by status %s from the database", status)

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := c.Postgres.
		Table(constants.TableBuild).
		Raw(dml.SelectBuildsCountByStatus, status).
		Pluck("count", &b).Error

	return b, err
}

// GetOrgBuildCount gets the count of all builds by repo ID from the database.
func (c *client) GetOrgBuildCount(org string, filters map[string]interface{}) (int64, error) {
	c.Logger.WithFields(logrus.Fields{
		"org": org,
	}).Tracef("getting count of builds for org %s from the database", org)

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := c.Postgres.
		Table(constants.TableBuild).
		Joins("JOIN repos ON builds.repo_id = repos.id and repos.org = ?", org).
		Where(filters).
		Count(&b).Error

	return b, err
}

// GetRepoBuildCount gets the count of all builds by repo ID from the database.
func (c *client) GetRepoBuildCount(r *library.Repo, filters map[string]interface{}) (int64, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"name": r.GetName(),
	}).Tracef("getting count of builds for repo %s from the database", r.GetFullName())

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := c.Postgres.
		Table(constants.TableBuild).
		Where("repo_id = ?", r.GetID()).
		Where(filters).
		Count(&b).Error

	return b, err
}
