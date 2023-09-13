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
