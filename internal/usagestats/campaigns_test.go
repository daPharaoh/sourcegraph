package usagestats

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/database/dbconn"
	"github.com/sourcegraph/sourcegraph/internal/database/dbtesting"
	"github.com/sourcegraph/sourcegraph/internal/extsvc"
	"github.com/sourcegraph/sourcegraph/internal/types"
)

func TestCampaignsUsageStatistics(t *testing.T) {
	ctx := context.Background()
	dbtesting.SetupGlobalTestDB(t)

	// Create stub repo.
	repoStore := database.Repos(dbconn.Global)
	esStore := database.ExternalServices(dbconn.Global)

	now := time.Now()
	svc := types.ExternalService{
		Kind:        extsvc.KindGitHub,
		DisplayName: "Github - Test",
		Config:      `{"url": "https://github.com"}`,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := esStore.Upsert(ctx, &svc); err != nil {
		t.Fatalf("failed to insert external services: %v", err)
	}
	repo := &types.Repo{
		Name: "test/repo",
		ExternalRepo: api.ExternalRepoSpec{
			ID:          fmt.Sprintf("external-id-%d", svc.ID),
			ServiceType: extsvc.TypeGitHub,
			ServiceID:   "https://github.com/",
		},
		Sources: map[string]*types.SourceInfo{
			svc.URN(): {
				ID:       svc.URN(),
				CloneURL: "https://secrettoken@test/repo",
			},
		},
	}
	if err := repoStore.Create(ctx, repo); err != nil {
		t.Fatal(err)
	}

	// Create a user.
	user, err := database.GlobalUsers.Create(ctx, database.NewUser{Username: "test"})
	if err != nil {
		t.Fatal(err)
	}

	// Create campaign specs 1, 2.
	_, err = dbconn.Global.Exec(`
		INSERT INTO campaign_specs
			(id, rand_id, raw_spec, namespace_user_id)
		VALUES
			(1, '123', '{}', $1),
			(2, '456', '{}', $1),
			(3, '789', '{}', $1)
	`, user.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Create event logs
	_, err = dbconn.Global.Exec(`
		INSERT INTO event_logs
			(id, name, argument, url, user_id, anonymous_user_id, source, version, timestamp)
		VALUES
		-- User 23, created a campaign last month and closes it
			(1, 'CampaignSpecCreated', '{"changeset_specs_count": 3}', '', 23, '', 'backend', 'version', date_trunc('month', CURRENT_DATE) - INTERVAL '2 days'),
			(2, 'CampaignSpecCreated', '{"changeset_specs_count": 1}', '', 23, '', 'backend', 'version', date_trunc('month', CURRENT_DATE) - INTERVAL '2 days'),
			(3, 'CampaignSpecCreated', '{}', '', 23, '', 'backend', 'version', date_trunc('month', CURRENT_DATE) - INTERVAL '2 days'),
			(4, 'ViewCampaignApplyPage', '{}', 'https://sourcegraph.test:3443/users/mrnugget/campaigns/apply/RANDID', 23, '5d302f47-9e91-4b3d-9e96-469b5601a765', 'WEB', 'version', date_trunc('month', CURRENT_DATE) - INTERVAL '2 days'),
			(5, 'CampaignCreated', '{"campaign_id": 1}', '', 23, '', 'backend', 'version', date_trunc('month', CURRENT_DATE) - INTERVAL '2 days'),
			(6, 'ViewCampaignDetailsPageAfterCreate', '{}', 'https://sourcegraph.test:3443/users/mrnugget/campaigns/gitignore-files', 23, '5d302f47-9e91-4b3d-9e96-469b5601a765', 'WEB', 'version', date_trunc('month', CURRENT_DATE) - INTERVAL '2 days'),
			(7, 'ViewCampaignDetailsPageAfterUpdate', '{}', 'https://sourcegraph.test:3443/users/mrnugget/campaigns/gitignore-files', 23, '5d302f47-9e91-4b3d-9e96-469b5601a765', 'WEB', 'version', date_trunc('month', CURRENT_DATE) - INTERVAL '2 days'),
			(8, 'ViewCampaignDetailsPagePage', '{}', 'https://sourcegraph.test:3443/users/mrnugget/campaigns/gitignore-files', 23, '5d302f47-9e91-4b3d-9e96-469b5601a765', 'WEB', 'version', date_trunc('month', CURRENT_DATE) - INTERVAL '2 days'),
			(9, 'CampaignCreatedOrUpdated', '{"campaign_id": 1}', '', 23, '', 'backend', 'version', date_trunc('month', CURRENT_DATE) - INTERVAL '2 days'),
			(10, 'CampaignClosed', '{"campaign_id": 1}', '', 23, '', 'backend', 'version', date_trunc('month', CURRENT_DATE) - INTERVAL '2 days'),
			(11, 'CampaignDeleted', '{"campaign_id": 1}', '', 23, '', 'backend', 'version', date_trunc('month', CURRENT_DATE) - INTERVAL '2 days'),
		-- User 24, created a campaign today and closes it
			(14, 'CampaignSpecCreated', '{}', '', 24, '', 'backend', 'version', now()),
			(15, 'ViewCampaignApplyPage', '{}', 'https://sourcegraph.test:3443/users/mrnugget/campaigns/apply/RANDID-2', 24, '5d302f47-9e91-4b3d-9e96-469b5601a765', 'WEB', 'version', now()),
			(16, 'CampaignCreated', '{"campaign_id": 2}', '', 24, '', 'backend', 'version', now()),
			(17, 'ViewCampaignDetailsPageAfterCreate', '{}', 'https://sourcegraph.test:3443/users/mrnugget/campaigns/foobar-files', 24, '5d302f47-9e91-4b3d-9e96-469b5601a765', 'WEB', 'version', now()),
			(18, 'ViewCampaignDetailsPageAfterUpdate', '{}', 'https://sourcegraph.test:3443/users/mrnugget/campaigns/foobar-files', 24, '5d302f47-9e91-4b3d-9e96-469b5601a765', 'WEB', 'version', now()),
			(19, 'CampaignCreatedOrUpdated', '{"campaign_id": 2}', '', 24, '', 'backend', 'version', now()),
			(20, 'CampaignClosed', '{"campaign_id": 2}', '', 24, '', 'backend', 'version', now()),
			(21, 'CampaignDeleted', '{"campaign_id": 2}', '', 24, '', 'backend', 'version', now()),
		-- User 25, only views the campaigns, today
			(27, 'ViewCampaignDetailsPagePage', '{}', 'https://sourcegraph.test:3443/users/mrnugget/campaigns/gitignore-files', 25, '5d302f47-9e91-4b3d-9e96-469b5601a765', 'WEB', 'version', now()),
			(28, 'ViewCampaignsListPage', '{}', 'https://sourcegraph.test:3443/users/mrnugget/campaigns', 25, '5d302f47-9e91-4b3d-9e96-469b5601a765', 'WEB', 'version', now()),
			(29, 'ViewCampaignDetailsPagePage', '{}', 'https://sourcegraph.test:3443/users/mrnugget/campaigns/foobar-files', 25, '5d302f47-9e91-4b3d-9e96-469b5601a765', 'WEB', 'version', now())
	`)
	if err != nil {
		t.Fatal(err)
	}

	campaignCreationDate1 := now.Add(-24 * 7 * 8 * time.Hour)  // 8 weeks ago
	campaignCreationDate2 := now.Add(-24 * 3 * time.Hour)      // 3 days ago
	campaignCreationDate3 := now.Add(-24 * 7 * 60 * time.Hour) // 60 weeks ago

	// Create campaigns 1, 2
	_, err = dbconn.Global.Exec(`
		INSERT INTO campaigns
			(id, name, campaign_spec_id, created_at, last_applied_at, namespace_user_id, closed_at)
		VALUES
			(1, 'test',   1, $2::timestamp, NOW(), $1, NULL),
			(2, 'test-2', 2, $3::timestamp, NOW(), $1, NOW()),
			(3, 'test-3', 3, $4::timestamp, NOW(), $1, NULL)
	`, user.ID, campaignCreationDate1, campaignCreationDate2, campaignCreationDate3)
	if err != nil {
		t.Fatal(err)
	}

	// Create 6 changesets.
	// 2 tracked: one OPEN, one MERGED.
	// 4 created by a campaign: 2 open (one with diffstat, one without), 2 merged (one with diffstat, one without)
	// missing diffstat shouldn't happen anymore (due to migration), but it's still a nullable field.
	_, err = dbconn.Global.Exec(`
		INSERT INTO changesets
			(id, repo_id, external_service_type, owned_by_campaign_id, campaign_ids, external_state, publication_state, diff_stat_added, diff_stat_changed, diff_stat_deleted)
		VALUES
		    -- tracked
			(1, $1, 'github', NULL, '{"1": {"detached": false}}', 'OPEN',   'PUBLISHED', 9, 7, 5),
			(2, $1, 'github', NULL, '{"2": {"detached": false}}', 'MERGED', 'PUBLISHED', 7, 9, 5),
			-- created by campaign
			(4,  $1, 'github', 1, '{"1": {"detached": false}}', 'OPEN',   'PUBLISHED', 5, 7, 9),
			(5,  $1, 'github', 1, '{"1": {"detached": false}}', 'OPEN',   'PUBLISHED', NULL, NULL, NULL),
			(6,  $1, 'github', 1, '{"1": {"detached": false}}', 'DRAFT',  'PUBLISHED', NULL, NULL, NULL),
			(7,  $1, 'github', 2, '{"2": {"detached": false}}',  NULL,    'UNPUBLISHED', 9, 7, 5),
			(8,  $1, 'github', 2, '{"2": {"detached": false}}', 'MERGED', 'PUBLISHED', 9, 7, 5),
			(9,  $1, 'github', 2, '{"2": {"detached": false}}', 'MERGED', 'PUBLISHED', NULL, NULL, NULL),
			(10, $1, 'github', 2, '{"2": {"detached": false}}',  NULL,    'UNPUBLISHED', 9, 7, 5),
			(11, $1, 'github', 2, '{"2": {"detached": false}}', 'CLOSED', 'PUBLISHED', NULL, NULL, NULL),
			(12, $1, 'github', 3, '{"3": {"detached": false}}', 'OPEN',   'PUBLISHED', 5, 7, 9),
			(13, $1, 'github', 3, '{"3": {"detached": false}}', 'OPEN',   'PUBLISHED', NULL, NULL, NULL)
	`, repo.ID)
	if err != nil {
		t.Fatal(err)
	}
	have, err := GetCampaignsUsageStatistics(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want := &types.CampaignsUsageStatistics{
		ViewCampaignApplyPageCount:               2,
		ViewCampaignDetailsPageAfterCreateCount:  2,
		ViewCampaignDetailsPageAfterUpdateCount:  2,
		CampaignsCount:                           3,
		CampaignsClosedCount:                     1,
		ActionChangesetsUnpublishedCount:         2,
		ActionChangesetsCount:                    8,
		ActionChangesetsDiffStatAddedSum:         19,
		ActionChangesetsDiffStatChangedSum:       21,
		ActionChangesetsDiffStatDeletedSum:       23,
		ActionChangesetsMergedCount:              2,
		ActionChangesetsMergedDiffStatAddedSum:   9,
		ActionChangesetsMergedDiffStatChangedSum: 7,
		ActionChangesetsMergedDiffStatDeletedSum: 5,
		ManualChangesetsCount:                    2,
		ManualChangesetsMergedCount:              1,
		CampaignSpecsCreatedCount:                4,
		ChangesetSpecsCreatedCount:               4,
		CurrentMonthContributorsCount:            1,
		CurrentMonthUsersCount:                   2,
		CampaignsCohorts: []*types.CampaignsCohort{
			{
				Week:                     campaignCreationDate1.Truncate(24 * 7 * time.Hour).Format("2006-01-02"),
				CampaignsOpen:            1,
				ChangesetsImported:       1,
				ChangesetsPublished:      3,
				ChangesetsPublishedOpen:  2,
				ChangesetsPublishedDraft: 1,
			},
			{
				Week:                      campaignCreationDate2.Truncate(24 * 7 * time.Hour).Format("2006-01-02"),
				CampaignsClosed:           1,
				ChangesetsImported:        1,
				ChangesetsUnpublished:     2,
				ChangesetsPublished:       3,
				ChangesetsPublishedMerged: 2,
				ChangesetsPublishedClosed: 1,
			},
			// campaign 3 should be ignored because it's too old
		},
	}
	if diff := cmp.Diff(want, have); diff != "" {
		t.Fatal(diff)
	}
}
