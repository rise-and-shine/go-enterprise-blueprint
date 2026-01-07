# Modular Enterprise Golang Application

---

## Documentation

### Document first approach

We maintain document first approach. We try to keep documentation in a human readable format, mostly in md files. Our documentation should serve as a source of truth. Everything (use cases, flows, tests) should be implemented based on the documentation. This way, we can ensure that the documentation and the implementation are in sync.

USE system-analyst skill when creating documentation.

### Docs structure

ABOUT.md - description of the current project, vision, goals, etc...

architecture/ - single architecture.md for describing code and infrastructure architectures. Should include references to architecture diagrams (UML, drawio, png, anything) in current folder.

flows/ - markdown files for describing business or technical flows of the system. Usually these files represent sequence or communication between different use cases. May include refenence to png, UML diagrams.

usecases/ - markdown files for structured definition of use cases. Each use case has it's type, and each type has it's strict documentation template. Each use case has it's own file. Also these files serve a as specification for our API.

uml/ - UML diagrams that are referenced from markdown files. Can be divided into different folders for different types of diagrams.

integrations/ - markdown files that describe how our system is integrated with other systems. Divided into different folders for different types of integrations or different systems.

misc/ - different files that don't fit into any other category. E.g.: pdf documentations, references to some resources or any other thing.

gen/dbdocs/ - generated files for database documentation, based on `tbls`

---

## Testing and development

USE golang-developer skill for testing and development.
