# AutoGet Organizer

This project provides AI-powered tools designed to automate the organization of downloaded files. The key functionalities include:

- **File Categorization and Movement**: Automatically moves downloaded files to appropriate directories based on their content and predefined categories.
- **JAV Actor Name Merging**: Consolidates and merges different names for the same JAV actor, ensuring consistent categorization and reducing redundancy.

### JAV Actor Name Merging

JAV movies are organized into folders named after the primary female actor. For movies featuring multiple actors, the file is stored in the folder of one designated actor. This is managed by an `actor.json` file, which maps various actor names to their canonical folder names:

```json
{
    "the folder name": ["name1", "name2"],
    "桥本有菜": ["桥本有菜", "橋本ありな", "橋本有菜", "新有菜", "新ありな"],
    ...
}
```

When a new movie is processed, the system first identifies the actor(s). It then checks if the actor's name already exists in `actor.json`. If not, a tool is run to generate potential name mappings. The system re-checks `actor.json` with these new mappings. If a match is found, the existing entry in `actor.json` is updated. If no match is found, a new folder is created for the actor, and `actor.json` is updated with the new entry.
