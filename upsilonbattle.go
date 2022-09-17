package upsilonbattle

import (
	"errors"
	"time"

	servicemgd "github.com/ecumeurs/service_mgd"
	"github.com/ecumeurs/upsilonbattle/battlearena"
	"github.com/google/uuid"
)

type UpsilonBattle struct {
	runningJobs map[uuid.UUID]*servicemgd.Job
	errors      []servicemgd.JobError

	battleArenas map[uuid.UUID]*battlearena.BattleArena
}

// Specific methods for UpsilonBattle

// AddBattleArena adds a battle arena to the list of battle arenas.
func (ub *UpsilonBattle) AddBattleArena(ba *battlearena.BattleArena) {
	ub.battleArenas[ba.Uuid] = ba
}

// RemoveBattleArena removes a battle arena from the list of battle arenas.
func (ub *UpsilonBattle) RemoveBattleArena(ba *battlearena.BattleArena) {
	delete(ub.battleArenas, ba.Uuid)
}

// BattleArena fetches a battle arena from the list of battle arenas.
func (ub *UpsilonBattle) BattleArena(ba uuid.UUID) (*battlearena.BattleArena, error) {
	if ba, ok := ub.battleArenas[ba]; ok {
		return ba, nil
	}
	return nil, errors.New("battle arena not found")
}

// BattleArenas returns the list of battle arenas.
func (ub *UpsilonBattle) BattleArenas() []*battlearena.BattleArena {
	bas := make([]*battlearena.BattleArena, 0, len(ub.battleArenas))
	for _, ba := range ub.battleArenas {
		bas = append(bas, ba)
	}
	return bas
}

// NewUpsilonBattle creates a new UpsilonBattle instance.
func NewUpsilonBattle() *UpsilonBattle {
	return &UpsilonBattle{
		runningJobs:  make(map[uuid.UUID]*servicemgd.Job),
		errors:       make([]servicemgd.JobError, 0),
		battleArenas: make(map[uuid.UUID]*battlearena.BattleArena),
	}
}

// Filling for prototype servicemgd

func (ub *UpsilonBattle) Name() string {
	return "UpsilonBattle"
}

func (ub *UpsilonBattle) Version() string {
	return "v0.0.1"
}

func (ub *UpsilonBattle) Running() bool {
	return true
}

func (ub *UpsilonBattle) Status() string {
	return "Running"
}

func (ub *UpsilonBattle) Features() []string {
	return []string{"BattleArena"}
}

func (ub *UpsilonBattle) Jobs() []uuid.UUID {
	jobs := make([]uuid.UUID, 0, len(ub.runningJobs))
	for _, job := range ub.runningJobs {
		jobs = append(jobs, job.Uuid)
	}
	return jobs
}

func (ub *UpsilonBattle) JobsList() []*servicemgd.Job {
	jobs := make([]*servicemgd.Job, 0, len(ub.runningJobs))
	for _, job := range ub.runningJobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// Job fetch the job from the list of running jobs.
func (ub *UpsilonBattle) Job(job uuid.UUID) (*servicemgd.Job, error) {
	if job, ok := ub.runningJobs[job]; ok {
		return job, nil
	}
	return nil, errors.New("job not found")
}

func (ub *UpsilonBattle) Abort(job uuid.UUID, reason string) (bool, error) {
	return false, errors.New("not implemented")
}

func (ub *UpsilonBattle) KeepErrorsFor(hours int) {
}

func (ub *UpsilonBattle) LastErrors() []servicemgd.JobError {
	return ub.errors
}

func (ub *UpsilonBattle) Lockdown(reason string) (bool, error) {
	return false, errors.New("not implemented")
}

func (ub *UpsilonBattle) Unlock(reason string) (bool, error) {
	return false, errors.New("not implemented")
}

func (ub *UpsilonBattle) GetLockdownStatus() servicemgd.LockdownStatus {
	return servicemgd.Unlocked
}

func (ub *UpsilonBattle) StoreState(reason string) (bool, error) {
	return false, errors.New("not implemented")
}

func (ub *UpsilonBattle) RestoreState(reason string) (bool, error) {
	return false, errors.New("not implemented")
}

func (ub *UpsilonBattle) LastStateDate() time.Time {
	return time.Now().UTC()
}
