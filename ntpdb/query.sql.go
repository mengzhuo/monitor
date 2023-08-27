// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: query.sql

package ntpdb

import (
	"context"
	"database/sql"
	"time"
)

const getMinLogScoreID = `-- name: GetMinLogScoreID :one
select id from log_scores order by id limit 1
`

// https://github.com/kyleconroy/sqlc/issues/1965
func (q *Queries) GetMinLogScoreID(ctx context.Context) (uint64, error) {
	row := q.db.QueryRowContext(ctx, getMinLogScoreID)
	var id uint64
	err := row.Scan(&id)
	return id, err
}

const getMonitorPriority = `-- name: GetMonitorPriority :many
select m.id, m.tls_name,
    avg(ls.rtt) / 1000 as avg_rtt,
    round((avg(ls.rtt)/1000) * (1+(2 * (1-avg(ls.step))))) as monitor_priority,
    avg(ls.step) as avg_step,
    if(avg(ls.step) < 0, false, true) as healthy,
    m.status as monitor_status, ss.status as status,
    count(*) as count
  from log_scores ls
  inner join monitors m
  left join server_scores ss on (ss.server_id = ls.server_id and ss.monitor_id = ls.monitor_id)
  where
    m.id = ls.monitor_id
  and ls.server_id = ?
  and m.type = 'monitor'
  and ls.ts > date_sub(now(), interval 12 hour)
  group by m.id, m.tls_name, m.status, ss.status
  order by healthy desc, monitor_priority, avg_step desc, avg_rtt
`

type GetMonitorPriorityRow struct {
	ID              uint32                 `json:"id"`
	TlsName         sql.NullString         `json:"tls_name"`
	AvgRtt          interface{}            `json:"avg_rtt"`
	MonitorPriority float64                `json:"monitor_priority"`
	AvgStep         interface{}            `json:"avg_step"`
	Healthy         interface{}            `json:"healthy"`
	MonitorStatus   MonitorsStatus         `json:"monitor_status"`
	Status          NullServerScoresStatus `json:"status"`
	Count           int64                  `json:"count"`
}

func (q *Queries) GetMonitorPriority(ctx context.Context, serverID uint32) ([]GetMonitorPriorityRow, error) {
	rows, err := q.db.QueryContext(ctx, getMonitorPriority, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMonitorPriorityRow
	for rows.Next() {
		var i GetMonitorPriorityRow
		if err := rows.Scan(
			&i.ID,
			&i.TlsName,
			&i.AvgRtt,
			&i.MonitorPriority,
			&i.AvgStep,
			&i.Healthy,
			&i.MonitorStatus,
			&i.Status,
			&i.Count,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMonitorTLSName = `-- name: GetMonitorTLSName :one
SELECT id, type, user_id, account_id, name, location, ip, ip_version, tls_name, api_key, status, config, client_version, last_seen, last_submit, created_on FROM monitors
WHERE tls_name = ? LIMIT 1
`

func (q *Queries) GetMonitorTLSName(ctx context.Context, tlsName sql.NullString) (Monitor, error) {
	row := q.db.QueryRowContext(ctx, getMonitorTLSName, tlsName)
	var i Monitor
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.UserID,
		&i.AccountID,
		&i.Name,
		&i.Location,
		&i.Ip,
		&i.IpVersion,
		&i.TlsName,
		&i.ApiKey,
		&i.Status,
		&i.Config,
		&i.ClientVersion,
		&i.LastSeen,
		&i.LastSubmit,
		&i.CreatedOn,
	)
	return i, err
}

const getScorerLogScores = `-- name: GetScorerLogScores :many
select ls.id, ls.monitor_id, ls.server_id, ls.ts, ls.score, ls.step, ls.offset, ls.rtt, ls.attributes from
  log_scores ls use index (primary),
  monitors m
WHERE
  ls.id > ? AND
  ls.id < ?+100000 AND
  m.type = 'monitor' AND
  monitor_id = m.id
ORDER by ls.id
LIMIT ?
`

type GetScorerLogScoresParams struct {
	LogScoreID uint64 `json:"log_score_id"`
	Limit      int32  `json:"limit"`
}

func (q *Queries) GetScorerLogScores(ctx context.Context, arg GetScorerLogScoresParams) ([]LogScore, error) {
	rows, err := q.db.QueryContext(ctx, getScorerLogScores, arg.LogScoreID, arg.LogScoreID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []LogScore
	for rows.Next() {
		var i LogScore
		if err := rows.Scan(
			&i.ID,
			&i.MonitorID,
			&i.ServerID,
			&i.Ts,
			&i.Score,
			&i.Step,
			&i.Offset,
			&i.Rtt,
			&i.Attributes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getScorerRecentScores = `-- name: GetScorerRecentScores :many
 select ls.id, ls.monitor_id, ls.server_id, ls.ts, ls.score, ls.step, ls.offset, ls.rtt, ls.attributes
   from log_scores ls
   inner join
   (select ls2.monitor_id, max(ls2.ts) as sts
      from log_scores ls2,
         monitors m,
         server_scores ss
      where ls2.server_id = ?
         and ls2.monitor_id=m.id and m.type = 'monitor'
         and (ls2.monitor_id=ss.monitor_id and ls2.server_id=ss.server_id)
         and ss.status in (?,?)
         and ls2.ts <= ?
         and ls2.ts >= date_sub(?, interval ? second)
      group by ls2.monitor_id
   ) as g
   where
     ls.server_id = ? AND
     g.sts = ls.ts AND
     g.monitor_id = ls.monitor_id
  order by ls.ts
`

type GetScorerRecentScoresParams struct {
	ServerID       uint32             `json:"server_id"`
	MonitorStatus  ServerScoresStatus `json:"monitor_status"`
	MonitorStatus2 ServerScoresStatus `json:"monitor_status_2"`
	Ts             time.Time          `json:"ts"`
	TimeLookback   interface{}        `json:"time_lookback"`
}

func (q *Queries) GetScorerRecentScores(ctx context.Context, arg GetScorerRecentScoresParams) ([]LogScore, error) {
	rows, err := q.db.QueryContext(ctx, getScorerRecentScores,
		arg.ServerID,
		arg.MonitorStatus,
		arg.MonitorStatus2,
		arg.Ts,
		arg.Ts,
		arg.TimeLookback,
		arg.ServerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []LogScore
	for rows.Next() {
		var i LogScore
		if err := rows.Scan(
			&i.ID,
			&i.MonitorID,
			&i.ServerID,
			&i.Ts,
			&i.Score,
			&i.Step,
			&i.Offset,
			&i.Rtt,
			&i.Attributes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getScorerStatus = `-- name: GetScorerStatus :many
select s.id, s.scorer_id, s.log_score_id, s.modified_on,m.name from scorer_status s, monitors m
WHERE m.type = 'score' and (m.id=s.scorer_id)
`

type GetScorerStatusRow struct {
	ID         uint32    `json:"id"`
	ScorerID   uint32    `json:"scorer_id"`
	LogScoreID uint64    `json:"log_score_id"`
	ModifiedOn time.Time `json:"modified_on"`
	Name       string    `json:"name"`
}

func (q *Queries) GetScorerStatus(ctx context.Context) ([]GetScorerStatusRow, error) {
	rows, err := q.db.QueryContext(ctx, getScorerStatus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetScorerStatusRow
	for rows.Next() {
		var i GetScorerStatusRow
		if err := rows.Scan(
			&i.ID,
			&i.ScorerID,
			&i.LogScoreID,
			&i.ModifiedOn,
			&i.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getScorers = `-- name: GetScorers :many
SELECT m.id as ID, s.id as status_id,
  m.status, s.log_score_id, m.name
FROM monitors m, scorer_status s
WHERE
  m.type = 'score'
  and m.status = 'active'
  and (m.id=s.scorer_id)
`

type GetScorersRow struct {
	ID         uint32         `json:"id"`
	StatusID   uint32         `json:"status_id"`
	Status     MonitorsStatus `json:"status"`
	LogScoreID uint64         `json:"log_score_id"`
	Name       string         `json:"name"`
}

func (q *Queries) GetScorers(ctx context.Context) ([]GetScorersRow, error) {
	rows, err := q.db.QueryContext(ctx, getScorers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetScorersRow
	for rows.Next() {
		var i GetScorersRow
		if err := rows.Scan(
			&i.ID,
			&i.StatusID,
			&i.Status,
			&i.LogScoreID,
			&i.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getServer = `-- name: GetServer :one
SELECT id, ip, ip_version, user_id, account_id, hostname, stratum, in_pool, in_server_list, netspeed, created_on, updated_on, score_ts, score_raw, deletion_on FROM servers WHERE id=?
`

func (q *Queries) GetServer(ctx context.Context, id uint32) (Server, error) {
	row := q.db.QueryRowContext(ctx, getServer, id)
	var i Server
	err := row.Scan(
		&i.ID,
		&i.Ip,
		&i.IpVersion,
		&i.UserID,
		&i.AccountID,
		&i.Hostname,
		&i.Stratum,
		&i.InPool,
		&i.InServerList,
		&i.Netspeed,
		&i.CreatedOn,
		&i.UpdatedOn,
		&i.ScoreTs,
		&i.ScoreRaw,
		&i.DeletionOn,
	)
	return i, err
}

const getServerIP = `-- name: GetServerIP :one
SELECT id, ip, ip_version, user_id, account_id, hostname, stratum, in_pool, in_server_list, netspeed, created_on, updated_on, score_ts, score_raw, deletion_on FROM servers WHERE ip=?
`

func (q *Queries) GetServerIP(ctx context.Context, ip string) (Server, error) {
	row := q.db.QueryRowContext(ctx, getServerIP, ip)
	var i Server
	err := row.Scan(
		&i.ID,
		&i.Ip,
		&i.IpVersion,
		&i.UserID,
		&i.AccountID,
		&i.Hostname,
		&i.Stratum,
		&i.InPool,
		&i.InServerList,
		&i.Netspeed,
		&i.CreatedOn,
		&i.UpdatedOn,
		&i.ScoreTs,
		&i.ScoreRaw,
		&i.DeletionOn,
	)
	return i, err
}

const getServerScore = `-- name: GetServerScore :one
SELECT id, monitor_id, server_id, score_ts, score_raw, stratum, status, created_on, modified_on FROM server_scores
  WHERE
    server_id=? AND
    monitor_id=?
`

type GetServerScoreParams struct {
	ServerID  uint32 `json:"server_id"`
	MonitorID uint32 `json:"monitor_id"`
}

func (q *Queries) GetServerScore(ctx context.Context, arg GetServerScoreParams) (ServerScore, error) {
	row := q.db.QueryRowContext(ctx, getServerScore, arg.ServerID, arg.MonitorID)
	var i ServerScore
	err := row.Scan(
		&i.ID,
		&i.MonitorID,
		&i.ServerID,
		&i.ScoreTs,
		&i.ScoreRaw,
		&i.Stratum,
		&i.Status,
		&i.CreatedOn,
		&i.ModifiedOn,
	)
	return i, err
}

const getServers = `-- name: GetServers :many
SELECT s.id, s.ip, s.ip_version, s.user_id, s.account_id, s.hostname, s.stratum, s.in_pool, s.in_server_list, s.netspeed, s.created_on, s.updated_on, s.score_ts, s.score_raw, s.deletion_on
    FROM servers s
    LEFT JOIN server_scores ss
        ON (s.id=ss.server_id)
WHERE (monitor_id = ?
    AND s.ip_version = ?
    AND (ss.score_ts IS NULL
          OR (ss.score_raw > -90 AND ss.status = "active"
               AND ss.score_ts < DATE_SUB( NOW(), INTERVAL ? second))
          OR (ss.score_raw > -90 AND ss.status = "testing"
              AND ss.score_ts < DATE_SUB( NOW(), INTERVAL ? second))
          OR (ss.score_ts < DATE_SUB( NOW(), INTERVAL 120 minute)))
    AND (s.score_ts IS NULL OR
        (s.score_ts < DATE_SUB( NOW(), INTERVAL ? second) ))
    AND (deletion_on IS NULL or deletion_on > NOW()))
ORDER BY ss.score_ts
LIMIT  ?
OFFSET ?
`

type GetServersParams struct {
	MonitorID              uint32           `json:"monitor_id"`
	IpVersion              ServersIpVersion `json:"ip_version"`
	IntervalSeconds        interface{}      `json:"interval_seconds"`
	IntervalSecondsTesting interface{}      `json:"interval_seconds_testing"`
	IntervalSecondsAll     interface{}      `json:"interval_seconds_all"`
	Limit                  int32            `json:"limit"`
	Offset                 int32            `json:"offset"`
}

func (q *Queries) GetServers(ctx context.Context, arg GetServersParams) ([]Server, error) {
	rows, err := q.db.QueryContext(ctx, getServers,
		arg.MonitorID,
		arg.IpVersion,
		arg.IntervalSeconds,
		arg.IntervalSecondsTesting,
		arg.IntervalSecondsAll,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Server
	for rows.Next() {
		var i Server
		if err := rows.Scan(
			&i.ID,
			&i.Ip,
			&i.IpVersion,
			&i.UserID,
			&i.AccountID,
			&i.Hostname,
			&i.Stratum,
			&i.InPool,
			&i.InServerList,
			&i.Netspeed,
			&i.CreatedOn,
			&i.UpdatedOn,
			&i.ScoreTs,
			&i.ScoreRaw,
			&i.DeletionOn,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getServersMonitorReview = `-- name: GetServersMonitorReview :many
select server_id from servers_monitor_review
where (next_review <= NOW() OR next_review is NULL)
order by next_review
limit 10
`

func (q *Queries) GetServersMonitorReview(ctx context.Context) ([]uint32, error) {
	rows, err := q.db.QueryContext(ctx, getServersMonitorReview)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uint32
	for rows.Next() {
		var server_id uint32
		if err := rows.Scan(&server_id); err != nil {
			return nil, err
		}
		items = append(items, server_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSystemSetting = `-- name: GetSystemSetting :one
select value from system_settings where ` + "`" + `key` + "`" + ` = ?
`

func (q *Queries) GetSystemSetting(ctx context.Context, key string) (string, error) {
	row := q.db.QueryRowContext(ctx, getSystemSetting, key)
	var value string
	err := row.Scan(&value)
	return value, err
}

const insertLogScore = `-- name: InsertLogScore :execresult
INSERT INTO log_scores
  (server_id, monitor_id, ts, score, step, offset, rtt, attributes)
  values (?, ?, ?, ?, ?, ?, ?, ?)
`

type InsertLogScoreParams struct {
	ServerID   uint32          `json:"server_id"`
	MonitorID  sql.NullInt32   `json:"monitor_id"`
	Ts         time.Time       `json:"ts"`
	Score      float64         `json:"score"`
	Step       float64         `json:"step"`
	Offset     sql.NullFloat64 `json:"offset"`
	Rtt        sql.NullInt32   `json:"rtt"`
	Attributes sql.NullString  `json:"attributes"`
}

func (q *Queries) InsertLogScore(ctx context.Context, arg InsertLogScoreParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertLogScore,
		arg.ServerID,
		arg.MonitorID,
		arg.Ts,
		arg.Score,
		arg.Step,
		arg.Offset,
		arg.Rtt,
		arg.Attributes,
	)
}

const insertScorer = `-- name: InsertScorer :execresult
insert into monitors
   (type, user_id, account_id,
    name, location, ip, ip_version,
    tls_name, api_key, status, config, client_version, created_on)
    VALUES ('score', NULL, NULL,
            ?, '', NULL, NULL,
            ?, NULL, 'active',
            '', '', NOW())
`

type InsertScorerParams struct {
	Name    string         `json:"name"`
	TlsName sql.NullString `json:"tls_name"`
}

func (q *Queries) InsertScorer(ctx context.Context, arg InsertScorerParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertScorer, arg.Name, arg.TlsName)
}

const insertScorerStatus = `-- name: InsertScorerStatus :exec
insert into scorer_status
   (scorer_id, log_score_id, modified_on)
   values (?,?,NOW())
`

type InsertScorerStatusParams struct {
	ScorerID   uint32 `json:"scorer_id"`
	LogScoreID uint64 `json:"log_score_id"`
}

func (q *Queries) InsertScorerStatus(ctx context.Context, arg InsertScorerStatusParams) error {
	_, err := q.db.ExecContext(ctx, insertScorerStatus, arg.ScorerID, arg.LogScoreID)
	return err
}

const insertServerScore = `-- name: InsertServerScore :exec
insert into server_scores
  (monitor_id, server_id, score_raw, created_on)
  values (?, ?, ?, ?)
`

type InsertServerScoreParams struct {
	MonitorID uint32    `json:"monitor_id"`
	ServerID  uint32    `json:"server_id"`
	ScoreRaw  float64   `json:"score_raw"`
	CreatedOn time.Time `json:"created_on"`
}

func (q *Queries) InsertServerScore(ctx context.Context, arg InsertServerScoreParams) error {
	_, err := q.db.ExecContext(ctx, insertServerScore,
		arg.MonitorID,
		arg.ServerID,
		arg.ScoreRaw,
		arg.CreatedOn,
	)
	return err
}

const listMonitors = `-- name: ListMonitors :many
SELECT id, type, user_id, account_id, name, location, ip, ip_version, tls_name, api_key, status, config, client_version, last_seen, last_submit, created_on FROM monitors
ORDER BY name
`

func (q *Queries) ListMonitors(ctx context.Context) ([]Monitor, error) {
	rows, err := q.db.QueryContext(ctx, listMonitors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Monitor
	for rows.Next() {
		var i Monitor
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.UserID,
			&i.AccountID,
			&i.Name,
			&i.Location,
			&i.Ip,
			&i.IpVersion,
			&i.TlsName,
			&i.ApiKey,
			&i.Status,
			&i.Config,
			&i.ClientVersion,
			&i.LastSeen,
			&i.LastSubmit,
			&i.CreatedOn,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateMonitorSeen = `-- name: UpdateMonitorSeen :exec
UPDATE monitors
  SET last_seen = ?
  WHERE id = ?
`

type UpdateMonitorSeenParams struct {
	LastSeen sql.NullTime `json:"last_seen"`
	ID       uint32       `json:"id"`
}

func (q *Queries) UpdateMonitorSeen(ctx context.Context, arg UpdateMonitorSeenParams) error {
	_, err := q.db.ExecContext(ctx, updateMonitorSeen, arg.LastSeen, arg.ID)
	return err
}

const updateMonitorSubmit = `-- name: UpdateMonitorSubmit :exec
UPDATE monitors
  SET last_submit = ?, last_seen = ?
  WHERE id = ?
`

type UpdateMonitorSubmitParams struct {
	LastSubmit sql.NullTime `json:"last_submit"`
	LastSeen   sql.NullTime `json:"last_seen"`
	ID         uint32       `json:"id"`
}

func (q *Queries) UpdateMonitorSubmit(ctx context.Context, arg UpdateMonitorSubmitParams) error {
	_, err := q.db.ExecContext(ctx, updateMonitorSubmit, arg.LastSubmit, arg.LastSeen, arg.ID)
	return err
}

const updateMonitorVersion = `-- name: UpdateMonitorVersion :exec
UPDATE monitors
  SET client_version = ?
  WHERE id = ?
`

type UpdateMonitorVersionParams struct {
	ClientVersion string `json:"client_version"`
	ID            uint32 `json:"id"`
}

func (q *Queries) UpdateMonitorVersion(ctx context.Context, arg UpdateMonitorVersionParams) error {
	_, err := q.db.ExecContext(ctx, updateMonitorVersion, arg.ClientVersion, arg.ID)
	return err
}

const updateScorerStatus = `-- name: UpdateScorerStatus :exec
update scorer_status
  set log_score_id = ?
  where scorer_id = ?
`

type UpdateScorerStatusParams struct {
	LogScoreID uint64 `json:"log_score_id"`
	ScorerID   uint32 `json:"scorer_id"`
}

func (q *Queries) UpdateScorerStatus(ctx context.Context, arg UpdateScorerStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateScorerStatus, arg.LogScoreID, arg.ScorerID)
	return err
}

const updateServer = `-- name: UpdateServer :exec
UPDATE servers
  SET score_ts  = ?,
      score_raw = ?
  WHERE
    id = ?
    AND (score_ts < ? OR score_ts is NULL)
`

type UpdateServerParams struct {
	ScoreTs  sql.NullTime `json:"score_ts"`
	ScoreRaw float64      `json:"score_raw"`
	ID       uint32       `json:"id"`
}

func (q *Queries) UpdateServer(ctx context.Context, arg UpdateServerParams) error {
	_, err := q.db.ExecContext(ctx, updateServer,
		arg.ScoreTs,
		arg.ScoreRaw,
		arg.ID,
		arg.ScoreTs,
	)
	return err
}

const updateServerScore = `-- name: UpdateServerScore :exec
UPDATE server_scores
  SET score_ts  = ?,
      score_raw = ?
  WHERE id = ?
`

type UpdateServerScoreParams struct {
	ScoreTs  sql.NullTime `json:"score_ts"`
	ScoreRaw float64      `json:"score_raw"`
	ID       uint64       `json:"id"`
}

func (q *Queries) UpdateServerScore(ctx context.Context, arg UpdateServerScoreParams) error {
	_, err := q.db.ExecContext(ctx, updateServerScore, arg.ScoreTs, arg.ScoreRaw, arg.ID)
	return err
}

const updateServerScoreStatus = `-- name: UpdateServerScoreStatus :exec
update server_scores
  set status = ?
  where monitor_id = ? and server_id = ?
`

type UpdateServerScoreStatusParams struct {
	Status    ServerScoresStatus `json:"status"`
	MonitorID uint32             `json:"monitor_id"`
	ServerID  uint32             `json:"server_id"`
}

func (q *Queries) UpdateServerScoreStatus(ctx context.Context, arg UpdateServerScoreStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateServerScoreStatus, arg.Status, arg.MonitorID, arg.ServerID)
	return err
}

const updateServerScoreStratum = `-- name: UpdateServerScoreStratum :exec
UPDATE server_scores
  SET stratum  = ?
  WHERE id = ?
`

type UpdateServerScoreStratumParams struct {
	Stratum sql.NullInt32 `json:"stratum"`
	ID      uint64        `json:"id"`
}

func (q *Queries) UpdateServerScoreStratum(ctx context.Context, arg UpdateServerScoreStratumParams) error {
	_, err := q.db.ExecContext(ctx, updateServerScoreStratum, arg.Stratum, arg.ID)
	return err
}

const updateServerStratum = `-- name: UpdateServerStratum :exec
UPDATE servers
  SET stratum = ?
  WHERE id = ?
`

type UpdateServerStratumParams struct {
	Stratum sql.NullInt32 `json:"stratum"`
	ID      uint32        `json:"id"`
}

func (q *Queries) UpdateServerStratum(ctx context.Context, arg UpdateServerStratumParams) error {
	_, err := q.db.ExecContext(ctx, updateServerStratum, arg.Stratum, arg.ID)
	return err
}

const updateServersMonitorReview = `-- name: UpdateServersMonitorReview :exec
update servers_monitor_review
  set last_review=NOW(), next_review=?
  where server_id=?
`

type UpdateServersMonitorReviewParams struct {
	NextReview sql.NullTime `json:"next_review"`
	ServerID   uint32       `json:"server_id"`
}

func (q *Queries) UpdateServersMonitorReview(ctx context.Context, arg UpdateServersMonitorReviewParams) error {
	_, err := q.db.ExecContext(ctx, updateServersMonitorReview, arg.NextReview, arg.ServerID)
	return err
}

const updateServersMonitorReviewChanged = `-- name: UpdateServersMonitorReviewChanged :exec
update servers_monitor_review
  set last_review=NOW(), last_change=NOW(), next_review=?
  where server_id=?
`

type UpdateServersMonitorReviewChangedParams struct {
	NextReview sql.NullTime `json:"next_review"`
	ServerID   uint32       `json:"server_id"`
}

func (q *Queries) UpdateServersMonitorReviewChanged(ctx context.Context, arg UpdateServersMonitorReviewChangedParams) error {
	_, err := q.db.ExecContext(ctx, updateServersMonitorReviewChanged, arg.NextReview, arg.ServerID)
	return err
}
