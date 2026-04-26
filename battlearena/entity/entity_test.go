package entity

import (
	"testing"

	"github.com/ecumeurs/upsilonbattle/battlearena/property"
	"github.com/ecumeurs/upsilonbattle/battlearena/property/defaultproperty"
	"github.com/google/uuid"
)

func TestEntity(t *testing.T) {
	e := New()
	if e.ID == uuid.Nil {
		t.Error("New() should not return nil")
	}
}

func TestEntityGetPropertyWithoutBuffs(t *testing.T) {
	e := New()

	e.Properties[property.HP.String()] = defaultproperty.MakeIntCounterProperty(property.HP, 10, 10, property.Public, property.Character)

	if e.GetProperty(property.HP).Get().(int) != 10 {
		t.Error("GetProperty() should return 10")
	}
}

func TestEntityGetPropertyWithBuffs(t *testing.T) {
	e := New()

	e.Properties[property.HP.String()] = defaultproperty.MakeIntCounterProperty(property.HP, 10, 10, property.Public, property.Character)
	tmpBuff := property.MakeTemporaryProperties(10)
	tmpBuff.Properties[property.HP.String()] = defaultproperty.MakeIntCounterProperty(property.HP, 0, 10, property.Public, property.Character)

	e.RegisterBuff(tmpBuff)
	if e.GetProperty(property.HP).(*defaultproperty.DefaultIntCounterProperty).MaxValue != 20 {
		t.Error("GetProperty() should return 20")
	}
}

func TestEntityGetPropertyWithBuffsAndNegativeValue(t *testing.T) {
	e := New()

	e.Properties[property.HP.String()] = defaultproperty.MakeIntCounterProperty(property.HP, 10, 10, property.Public, property.Character)
	tmpBuff := property.MakeTemporaryProperties(10)
	tmpBuff.Properties[property.HP.String()] = defaultproperty.MakeIntCounterProperty(property.HP, 0, -5, property.Public, property.Character)

	e.RegisterBuff(tmpBuff)
	if e.GetProperty(property.HP).(*defaultproperty.DefaultIntCounterProperty).MaxValue != 5 {
		t.Error("GetProperty() should return 5")
	}
}

func TestBuffGetRemovedAfterTime(t *testing.T) {
	e := New()

	e.Properties[property.HP.String()] = defaultproperty.MakeIntCounterProperty(property.HP, 10, 10, property.Public, property.Character)
	tmpBuff := property.MakeTemporaryProperties(5)
	tmpBuff.Properties[property.HP.String()] = defaultproperty.MakeIntCounterProperty(property.HP, 0, 10, property.Public, property.Character)

	e.RegisterBuff(tmpBuff)

	for i := 0; i < 5; i++ {
		e.BuffTickDown()
	}

	if e.GetProperty(property.HP).(*defaultproperty.DefaultIntCounterProperty).MaxValue != 10 {
		t.Error("GetProperty() should return 10")
	}
}

func TestRemoveBuffsByOrigin(t *testing.T) {
	e := New()
	origin1 := uuid.New()
	origin2 := uuid.New()

	buff1 := property.TemporaryProperties{Forever: true, OriginEntityID: origin1, Properties: make(map[string]property.Property)}
	buff2 := property.TemporaryProperties{Forever: true, OriginEntityID: origin2, Properties: make(map[string]property.Property)}
	buff3 := property.TemporaryProperties{Forever: true, OriginEntityID: origin1, Properties: make(map[string]property.Property)}

	e.RegisterBuff(buff1)
	e.RegisterBuff(buff2)
	e.RegisterBuff(buff3)

	if len(e.Buffs) != 3 {
		t.Errorf("Expected 3 buffs, got %d", len(e.Buffs))
	}

	e.RemoveBuffsByOrigin(origin1)

	if len(e.Buffs) != 1 {
		t.Errorf("Expected 1 buff after removal, got %d", len(e.Buffs))
	}

	if e.Buffs[0].OriginEntityID != origin2 {
		t.Error("Remaining buff should be origin2")
	}
}
