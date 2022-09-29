# .\battlearena\entity\entitygenerator

[Up](../README.md)

# Entity Generator

this module allow to create entity based on a template. This will be used at grid generation to populate the map. 
This won't be implemented in v0.0.2 but maybe soon ;)

For now, expect entity generator to create entities with random properties being set.

## Entity Template

Entity template are stored in the `templates/entities` folder. They are simple json files, with the following structure:

```json
{
    "name": "template name",
    "entity": {
        "name": "entity name",
        "properties": {
            "property1": {"min": 1, "max": 10},
            "property2": "value2"
        }
    }
}
```