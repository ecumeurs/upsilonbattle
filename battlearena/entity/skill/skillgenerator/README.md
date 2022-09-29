# .\battlearena\entity\skill\skillgenerator

[Up](../README.md)

# Skill Generator

this module allow to create skill based on a template. 
This won't be implemented in v0.0.3 but maybe soon ;)

For now, expect skill generator to create skill with random (skill)properties being set.

## Skill Template

Entity template are stored in the `templates/skills` folder. They are simple json files, with the following structure:

```json
{
    "name": "template name",
    "skill": {
        "name": "skill name",
        "properties": {
            "property1": {"min": 1, "max": 10},
            "property2": "value2"
        }
    }
}
```