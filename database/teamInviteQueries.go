package database

const TeamInviteCreateQueryString = `
INSERT INTO team_invites (
  team_id,
  email,
  invite_creator_id
) VALUES (
  $1, $2, $3
)
RETURNING *
;
`

const TeamInvitesForCurrentUserQueryString = `
SELECT * FROM team_invites
WHERE email = $1
;
`

const TeamInviteAcceptQueryString = `
UPDATE team_invites
SET status = 'accepted'
WHERE team_invite_id = $1
RETURNING *
;
`

const TeamInviteGetByIdQueryString = `
SELECT * FROM team_invites
WHERE team_invite_id = $1
;
`
const TeamInviteDeclineQueryString = `
UPDATE team_invites
SET status = 'declined'
WHERE team_invite_id = $1
;
`

const TeamInviteDeleteQueryString = `
DELETE FROM team_invites
WHERE team_invite_id = $1
;
`
