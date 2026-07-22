// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"context"
	"fmt"

	"code.gitea.io/gitea/models/db"
)

// ReputationLabel represent a reputation label
type ReputationLabel struct {
	ID          int64  `xorm:"pk autoincr"`
	Name        string `xorm:"UNIQUE"`
	Description string
	Color       string `xorm:"VARCHAR(7)"`
}

// UserReputationLabel represents a relation between user and //
// ReputationLabel
type UserReputationLabel struct {
	ID                int64 `xorm:"pk autoincr"`
	ReputationLabelID int64
	UserID            int64 `xorm:"INDEX"`
}

// UserReputationLabel represents a relation between user and //
// ReputationLabel
type RepoReputationLabel struct {
	ID                int64 `xorm:"pk autoincr"`
	ReputationLabelID int64
	RepoID            int64 `xorm:"INDEX"`
}

func init() {
	db.RegisterModel(new(ReputationLabel))
	db.RegisterModel(new(UserReputationLabel))
	db.RegisterModel(new(RepoReputationLabel))
}

// GetUserReputationLabels returns the user's reputation labels.
func GetUserReputationLabels(ctx context.Context, u *User) ([]*ReputationLabel, int64, error) {
	sess := db.GetEngine(ctx).
		Select("`reputation_label`.*").
		Join("INNER", "user_reputation_label", "`user_reputation_label`.reputation_label_id=reputation_label.id").
		Where("user_reputation_label.user_id=?", u.ID)

	labels := make([]*ReputationLabel, 0, 8)
	count, err := sess.FindAndCount(&labels)
	return labels, count, err
}

// CreateReputationLabel creates a new reputation label.
func CreateReputationLabel(ctx context.Context, label *ReputationLabel) error {
	_, err := db.GetEngine(ctx).Insert(label)
	return err
}

// GetAllReputationLabels returns all reputation label.
func GetAllReputationLabels(ctx context.Context) ([]*ReputationLabel, error) {
	labels := make([]*ReputationLabel, 0)
	return labels, db.GetEngine(ctx).OrderBy("id").Find(&labels)
}

// GetReputationLabel returns a reputation label.
func GetReputationLabel(ctx context.Context, name string) (*ReputationLabel, error) {
	label := new(ReputationLabel)
	has, err := db.GetEngine(ctx).Where("name=?", name).Get(label)
	if !has {
		return nil, err
	}
	return label, err
}

// UpdateReputationLabel updates a label based on its name.
func UpdateReputationLabel(ctx context.Context, label *ReputationLabel) error {
	_, err := db.GetEngine(ctx).Where("name=?", label.Name).Update(label)
	return err
}

// DeleteReputationLabel deletes a label.
func DeleteReputationLabel(ctx context.Context, label *ReputationLabel) error {
	_, err := db.GetEngine(ctx).Where("name=?", label.Name).Delete(label)
	return err
}

// AddUserReputationLabel adds a reputation label to a user.
func AddUserReputationLabel(ctx context.Context, u *User, label *ReputationLabel) error {
	return AddUserReputationLabels(ctx, u, []*ReputationLabel{label})
}

// AddUserReputationLabels adds labels to a user.
func AddUserReputationLabels(ctx context.Context, u *User, labels []*ReputationLabel) error {
	return db.WithTx(ctx, func(ctx context.Context) error {
		for _, label := range labels {
			// hydrate label and check if it exists
			has, err := db.GetEngine(ctx).Where("name=?", label.Name).Get(label)
			if err != nil {
				return err
			} else if !has {
				return fmt.Errorf("label with name %s doesn't exist", label.Name)
			}
			if err := db.Insert(ctx, &UserReputationLabel{
				ReputationLabelID: label.ID,
				UserID:  u.ID,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

// RemoveUserReputationLabel removes a label from a user.
func RemoveUserReputationLabel(ctx context.Context, u *User, label *ReputationLabel) error {
	return RemoveUserReputationLabels(ctx, u, []*ReputationLabel{label})
}

// RemoveUserReputationLabels removes labels from a user.
func RemoveUserReputationLabels(ctx context.Context, u *User, labels []*ReputationLabel) error {
	return db.WithTx(ctx, func(ctx context.Context) error {
		for _, label := range labels {
			if _, err := db.GetEngine(ctx).
				Join("INNER", "reputation_label", "reputation_label.id = `user_reputation_label`.reputation_label_id").
				Where("`user_reputation_label`.user_id=? AND `reputation_label`.name=?", u.ID, label.Name).
				Delete(&UserReputationLabel{}); err != nil {
				return err
			}
		}
		return nil
	})
}

// RemoveAllUserReputationLabels removes all labels from a user.
func RemoveAllUserReputationLabels(ctx context.Context, u *User) error {
	_, err := db.GetEngine(ctx).Where("user_id=?", u.ID).Delete(&UserReputationLabel{})
	return err
}

// AddRepoReputationLabel adds a reputation label to a repo.
func AddRepoReputationLabel(ctx context.Context, rID int64, label *ReputationLabel) error {
	return AddRepoReputationLabels(ctx, rID, []*ReputationLabel{label})
}

// AddRepoReputationLabels adds labels to a repo.
func AddRepoReputationLabels(ctx context.Context, rID int64, labels []*ReputationLabel) error {
	return db.WithTx(ctx, func(ctx context.Context) error {
		for _, label := range labels {
			// hydrate label and check if it exists
			has, err := db.GetEngine(ctx).Where("name=?", label.Name).Get(label)
			if err != nil {
				return err
			} else if !has {
				return fmt.Errorf("label with name %s doesn't exist", label.Name)
			}
			if err := db.Insert(ctx, &RepoReputationLabel{
				ReputationLabelID: label.ID,
				RepoID:  rID,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

// RemoveRepoReputationLabel removes a label from a repo.
func RemoveRepoReputationLabel(ctx context.Context, rID int64, label *ReputationLabel) error {
	return RemoveRepoReputationLabels(ctx, rID, []*ReputationLabel{label})
}

// RemoveRepoReputationLabels removes labels from a repo.
func RemoveRepoReputationLabels(ctx context.Context, rID int64, labels []*ReputationLabel) error {
	return db.WithTx(ctx, func(ctx context.Context) error {
		for _, label := range labels {
			if _, err := db.GetEngine(ctx).
				Join("INNER", "reputation_label", "reputation_label.id = `repo_reputation_label`.reputation_label_id").
				Where("`repo_reputation_label`.repo_id=? AND `reputation_label`.name=?", rID, label.Name).
				Delete(&RepoReputationLabel{}); err != nil {
				return err
			}
		}
		return nil
	})
}

// RemoveAllRepoReputationLabels removes all labels from a repo.
func RemoveAllRepoReputationLabels(ctx context.Context, rID int64) error {
	_, err := db.GetEngine(ctx).Where("repo_id=?", rID).Delete(&RepoReputationLabel{})
	return err
}
