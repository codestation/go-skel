Controller
========

The `controller` package contains the implementation of the application interfaces and adapters, which provide the 
means for the user to interact with the system and for the system to interact with external services and systems.

This layer includes REST APIs, gRPC services, CLI commands, and other user-facing components that translate user input 
into use case invocations and present the results of those invocations to the user. By separating the delivery logic
from other layers of the application, the delivery package promotes flexibility, scalability, and maintainability 
of the codebase.
