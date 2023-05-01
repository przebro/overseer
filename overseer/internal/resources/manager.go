package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/datastore"
	"github.com/przebro/overseer/overseer/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const collectionName = "resources"

type ResourceManagerImpl struct {
	log   zerolog.Logger
	flock sync.Mutex
	rw    *resourceReadWriter
}

// TicketManager - base resources required by task to run
type TicketManager interface {
	Add(ctx context.Context, name string, odate date.Odate) (bool, error)
	Delete(ctx context.Context, name string, odate date.Odate) (bool, error)
}

// FlagManager - Resource that helps run tasks
type FlagManager interface {
	Set(ctx context.Context, name string, policy uint8) (bool, error)
	DestroyFlag(ctx context.Context, name string) (bool, error)
}

// ResourceManager - manages resources that are required by tasks
type ResourceManager interface {
	TicketManager
	FlagManager
	ListResources(ctx context.Context, filter ResourceFilter) []ResourceModel
}

// NewManager - crates new resources manager
func NewManager(rconfig config.ResourcesConfigurartion, provider *datastore.Provider) (*ResourceManagerImpl, error) {

	var err error

	rw, err := newResourceReadWriter(collectionName, provider)
	if err != nil {
		return nil, err
	}

	rm := &ResourceManagerImpl{
		log:   log.With().Str("component", "resource-manager").Logger(),
		rw:    rw,
		flock: sync.Mutex{},
	}

	return rm, nil
}
func (rm *ResourceManagerImpl) Add(ctx context.Context, name string, odate date.Odate) (bool, error) {

	var odateval int64 = 0
	if odate != date.OdateNone {
		odateval = odate.ToUnix()
	}
	err := rm.rw.Insert(ctx, &ResourceModel{
		ID:    getkey(ResourceTypeTicket, name, string(odate)),
		Type:  ResourceTypeTicket,
		Value: odateval,
	})

	if err != nil {
		return false, errors.New("ticket with given name and odate already exists")
	}
	rm.log.Info().Msg(fmt.Sprintf("ticket:%s,odate:%s", name, odate))

	return true, nil
}

func (rm *ResourceManagerImpl) Delete(ctx context.Context, name string, odate date.Odate) (bool, error) {

	key := getkey(ResourceTypeTicket, name, string(odate))

	if err := rm.rw.Delete(ctx, key); err != nil {
		return false, err
	}

	return true, nil

}

func (rm *ResourceManagerImpl) CheckTickets(in []types.CollectedTicketModel) []types.CollectedTicketModel {

	result := make([]types.CollectedTicketModel, 0, len(in))
	for _, item := range in {
		model := ResourceModel{}
		key := getkey(ResourceTypeTicket, item.Name, string(item.Odate))
		exists := false
		if err := rm.rw.Get(context.TODO(), key, &model); err == nil {
			exists = true
		}

		result = append(result, types.CollectedTicketModel{Name: item.Name, Odate: item.Odate, Exists: exists})
	}

	return result
}

func (rm *ResourceManagerImpl) Set(ctx context.Context, name string, policy uint8) (bool, error) {

	defer rm.flock.Unlock()
	rm.flock.Lock()

	key := getkey(ResourceTypeTicket, name, "")
	v := ResourceModel{}

	if err := rm.rw.Get(ctx, key, &v); err != nil {
		rm.rw.Insert(ctx, &ResourceModel{
			ID:    getkey(ResourceTypeTicket, name, ""),
			Type:  ResourceTypeFlag,
			Value: pack(policy, 0),
		})

		return true, nil
	}

	p, cnt := unpack(v.Value)

	if p == FlagPolicyExclusive {

		rm.log.Debug().Str("name", name).Str("type", "EXL").Msg("ACQ ERR FLAG")
		return false, errors.New("flag in use with exclusive policy")
	}

	if p == FlagPolicyShared && policy == FlagPolicyExclusive && cnt != 0 {
		rm.log.Debug().Str("name", name).Str("type", "SHR").Msg("ACQ ERR FLAG")
		return false, errors.New("unable to set shared, flag in use with exclusive policy")
	}

	cnt++
	p = policy
	v.Value = pack(p, cnt)
	rm.log.Debug().Str("name", v.ID).Str("type", "SHR").Msg("ACQ FLAG SUCCESS")
	rm.rw.Update(ctx, &v)

	return true, nil
}

// Unset - remove a flag
func (rm *ResourceManagerImpl) unset(name string) (bool, error) {

	defer rm.flock.Unlock()
	rm.flock.Lock()

	v := ResourceModel{}

	key := getkey(ResourceTypeFlag, name, "")
	if err := rm.rw.Get(context.TODO(), key, &v); err != nil {
		return false, errors.New("flag with given name does not exists")
	}

	_, cnt := unpack(v.Value)
	cnt--

	rm.log.Debug().Str("flag", string(v.ID)).Int64("count", v.Value).Msg("unset")
	if cnt == 0 {
		rm.rw.Delete(context.TODO(), key)
		rm.log.Debug().Str("flag", string(v.ID)).Msg("flag removed")
	} else {
		v.Value--
		rm.rw.Update(context.TODO(), &v)
	}

	return true, nil

}

func (rm *ResourceManagerImpl) DestroyFlag(ctx context.Context, name string) (bool, error) {

	defer rm.flock.Unlock()
	rm.flock.Lock()

	key := getkey(ResourceTypeFlag, name, "")
	v := ResourceModel{}
	if err := rm.rw.Get(ctx, key, &v); err != nil {
		return false, errors.New("flag with given name does not exists")
	}

	if err := rm.rw.Delete(ctx, key); err != nil {
		return false, err
	}

	return true, nil
}

func (rm *ResourceManagerImpl) ListResources(ctx context.Context, filter ResourceFilter) []ResourceModel {

	result := []ResourceModel{}

	for _, rtype := range filter.Type {

		if rtype == ResourceTypeTicket {
			if filter.TicketFilterOptions != nil {
				result = append(result, rm.findTickets(ctx, filter.Name, filter.TicketFilterOptions)...)
			}
		}
		if rtype == ResourceTypeFlag {
			if filter.FlagFilterOptions != nil {
				result = append(result, rm.findFlags(ctx, filter.Name, filter.FlagFilterOptions)...)
			}
		}
	}

	return result
}

func (rm *ResourceManagerImpl) ProcessTicketAction(tickets []types.TicketActionModel) bool {

	var result bool = true

	for _, ticket := range tickets {
		if ticket.Action == "ADD" {
			if ok, e := rm.Add(context.TODO(), ticket.Name, ticket.Odate); !ok {
				rm.log.Error().Str("action", string(ticket.Action)).Err(e).Msg("ticket action error")
				result = false
			}
		} else {
			if ok, e := rm.Delete(context.TODO(), ticket.Name, ticket.Odate); !ok {
				rm.log.Error().Str("action", string(ticket.Action)).Err(e).Msg("ticket action error")
				result = false
			}
		}
	}
	return result
}

func (rm *ResourceManagerImpl) ProcessAcquireFlag(input []types.FlagModel) (bool, []string) {
	var requiredFlagNames []string = []string{}
	ok := true

	aflags := []string{}

	for _, f := range input {
		policy := f.Policy
		// if flag is not acquired, rollback changes that were made
		if ok, _ = rm.Set(context.TODO(), f.Name, policy); !ok {
			requiredFlagNames = append(requiredFlagNames, f.Name)
			break
		}
		aflags = append(aflags, f.Name)
	}

	if !ok {

		for _, f := range aflags {
			rm.unset(f)
		}
	}

	return ok, requiredFlagNames
}
func (rm *ResourceManagerImpl) ProcessReleaseFlag(input []string) (bool, []string) {
	ok := true

	var flagNames []string = []string{}
	var success bool = true

	for _, f := range input {
		if ok, _ := rm.unset(f); !ok {
			flagNames = append(flagNames, f)
			success = success && false
		} else {
			success = success && true
		}
	}

	return ok, flagNames
}

// Start - starts the task pool
func (rm *ResourceManagerImpl) Start() error {

	return nil
}

// Shutdown - shutdowns task pool
func (rm *ResourceManagerImpl) Shutdown() error {

	return nil
}

func (rm *ResourceManagerImpl) findTickets(ctx context.Context, name string, options *TicketFilterOptions) []ResourceModel {

	result := []ResourceModel{}

	crsr, err := rm.rw.AllTickets(ctx)
	if err != nil {
		return result
	}
	for crsr.Next(ctx) {
		var v ResourceModel
		if err := crsr.Decode(&v); err != nil {
			continue
		}
		v.ID = getname(v.ID)
		result = append(result, v)
	}

	return result
}

func (rm *ResourceManagerImpl) findFlags(ctx context.Context, name string, options *FlagFilterOptions) []ResourceModel {

	result := []ResourceModel{}

	crsr, err := rm.rw.AllFlags(ctx)
	if err != nil {
		return result
	}
	for crsr.Next(ctx) {
		var v ResourceModel
		if err := crsr.Decode(&v); err != nil {
			continue
		}
		v.ID = getname(v.ID)
		result = append(result, v)
	}

	return result
}

func buildExpr(value string) string {

	expr := ""

	if value == "" {
		return `[\w\-]*|^$`
	}

	expr += "^"

	for _, c := range value {

		if c == '*' {
			expr += `[\w\-]*`
			continue
		}
		if c == '?' {
			expr += `[\w\-]{1}`
			continue
		}

		expr += string(c)
	}

	expr += "$"

	return expr
}

func buildDateExpr(value string) string {

	expr := ""

	if value == "" {
		return `[\d]*|^$`
	}

	expr += "^"

	for _, c := range value {

		if c == '*' {
			expr += `[\d]*`
			continue
		}
		if c == '?' {
			expr += `[\d]{1}`
			continue
		}

		expr += string(c)
	}

	expr += "$"

	return expr
}

func getname(key string) string {
	return strings.Split(key, ":")[1]
}
func getkey(rtype ResourceType, name, value string) string {
	if rtype == ResourceTypeTicket {
		return "t:" + name + ":" + value
	}

	return "f:" + name
}

// pack - packs the cnt in lower 32 bits and p in upper 32 bits
func pack(p uint8, cnt uint32) int64 {
	return int64(p)<<32 | int64(cnt)

}

func unpack(p int64) (uint8, uint32) {
	return uint8(p >> 32), uint32(p)
}
