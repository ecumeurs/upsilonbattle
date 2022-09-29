# .\battlearena\entity\skills

[Up](../README.md)

# Skills

Skills are the way to make an entity act in a special way. Some will be attack skills, that will perform an attack assisted by various properties.
Skills also build up upon the entities properties, and will use them to compute the result of the skill.

For example, if entity has a CriticChance property, then all related skills and basic attack will make use of it. Absence of this entity's property will disallow the skills to use it.
Some skill may specifically require a property to be used, and will fail if the property is missing.
Some skill will provide a property to the entity (buffs/debuffs). In consequence, if an entity is subjected to a "CriticChance" buff, then all skills and basic attack will make use of it.

Skills, Properties and Rules are the three main components of the battlearena. They are all linked together and are used to compute the result of an attack.

NOTE: Skills are only the structure handling things, not the working of the skills. That's the role of the rules. (see [Rules](../rules/README.md))

Eventually the skill will be able to level up and handle experience gains, it will be handled here.

## Structure 

Skills are defined by four main components:

* A name: allows the entity's controller to use said skill.
* A targeting mechanism: tell what is affected by the skill, in which way. These are properties (see [properties](../properties/README.md))
  * Targeting mechanism also encompass the medium: it will differentiate between a skill targeting an entity(how to compute hit/miss), multiples entities(might not use the same hit/miss computation), a position(no hit or miss), delayed actions(reaction), counters(another kind of reaction), enviromental effects(skills affecting only the terrain), traps(delayed actions stored on cells).
  * Some skills will only be able to target in Line of Sight, others will be able to target anywhere on the map, or they might be able to use LOS of friendly entities.
* An effet: the effect of the skill on the target(s). These will also be properties. The way the interaction will be computed is not yet defined, but several options are available: 
  * Depending on the presence of some properties, the computation will change. Example: a skill that has a damaging property will invoke the damage dealing computation (see [rules](../rules/README.md))]), while a skill that has a healing property will invoke the healing computation. Note: a skill might invoke multiples rules at the same time. 
  * A specific property might be used to define the effect of the skill. All skill would get the "EffectType" property, and the value of this property would be used to invoke the correct computation. This would prevent more exotic skills to be defined. Less Chaos, more order.
  * Expect multiple kind of computations:
    * Attack/damaging skills.
    * Damage Reducing skills (especially in reaction skills, the other kinds will be considered as buffs)
    * Healing skills.
    * Buffing/cursing skills. 
    * Movement & placement skills.
* A Cost and a Cooldown: 
  * At the moment, no cost will be defined, some venues may be explored later on.
    * Cooldown duration
    * HP/SP/MP
    * Entity delay 
  * Cooldown enforce the skill to be unable to be used again for a certain amount of time. 




##  Futur

* Skills will be able to level up and evolve allowing one to upgrade skills specific properties.
* Skills will be able to be merged. This will allow to create more complex skills, and to create skills that are not available to the player.
* Eventually skills will be tradable (in a way) and will be able to be used by other entities.
* Skill requirements need to be explored:
  * Minimum ressource to be lost before being able to use the skill
  * Properties required to be able to use the skill (ex: a skill that require to have the Poisonned property to be present, up to a certain level; or a skill that won't be able to be used if the property is present(curses))
  * Number of foes in range, type of foes, etc.
  * Number of turns since the begining of the battle.
  * Item presence requirement(firearm, sword, etc.)
  * Item consumption (amunitions, components, etc.)
  * Grid element presence (water, lava, tree, etc.)
  